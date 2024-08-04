package server

import (
	"errors"
	"fmt"
	"rush/attendance"
	"rush/auth"
	"rush/session"
	"rush/user"
	"sort"
	"time"
)

type User struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	University string  `json:"university"`
	Phone      string  `json:"phone"`
	Generation float64 `json:"generation"`
	IsActive   bool    `json:"is_active"`
}

type Session struct {
	Id            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	HostedBy      string    `json:"hosted_by"`
	CreatedBy     string    `json:"created_by"`
	GoogleFormUri string    `json:"google_form_uri"`
	JoinningUsers []string  `json:"joinning_users"`
	CreatedAt     time.Time `json:"created_at"`
	StartsAt      time.Time `json:"starts_at"`
	Score         int       `json:"score"`
	IsClosed      bool      `json:"is_closed"`
}

type Attendance struct {
	Id          string    `json:"id"`
	SessionId   string    `json:"session_id"`
	SessionName string    `json:"session_name"`
	UserId      string    `json:"user_id"`
	UserName    string    `json:"user_name"`
	JoinedAt    time.Time `json:"joined_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type tokenInspector interface {
	// Handles the third party token that is used for signing in.
	GetUserIdentifier(token string) (auth.UserIdentifier, error)
	// Returns the provider of the token. authHandler uses it to extract the email address.
	Provider() auth.Provider
}

type authHandler interface {
	// Handles the rush token that is used for API calls after signing in.
	GetUserIdentifier(token string) (auth.UserIdentifier, error)
	// Returns the rush token that is used for API calls after signing in.
	SignIn(userIdentifier auth.UserIdentifier) (string, error)
}

type userRepo interface {
	Get(id string) (*user.User, error)
	GetAll() ([]user.User, error)
	// Skips `offset` users and returns up to `pageSize` users, an indicator if it has more users and total count.
	List(offset int, pageSize int) (*user.ListResult, error)
	Add(user *user.User) error
	GetByEmail(email string) (*user.User, error)
	GetAllByExternalIds(externalIds []string) ([]user.User, error)
}

type sessionRepo interface {
	Get(id string) (*session.Session, error)
	GetAll() ([]session.Session, error)
	List(offset int, pageSize int) (*session.ListResult, error)
	Add(name string, description string, hostedBy int, createdBy int, startsAt time.Time, score int) (string, error)
	Update(id string, updateForm *session.UpdateForm) (*session.Session, error)
}

type attendanceFormHandler interface {
	GenerateForm(title string, description string, users []user.User) (attendance.Form, error)
	GetSubmissions(formId string) ([]attendance.FormSubmission, error)
}

type attendanceRepo interface {
	BulkInsert(requests []attendance.AddAttendanceReq) error
	FindByUserId(userId string) ([]attendance.Attendance, error)
}

type Server struct {
	tokenInspector        tokenInspector
	authHandler           authHandler
	userRepo              userRepo
	sessionRepo           sessionRepo
	attendanceFormHandler attendanceFormHandler
	attendanceRepo        attendanceRepo
	formTimeLocation      *time.Location
}

func New(tokenInspector tokenInspector, authHandler authHandler, userRepo userRepo, sessionRepo sessionRepo,
	attendanceFormHandler attendanceFormHandler, attendanceRepo attendanceRepo, formTimeLocation *time.Location) *Server {
	return &Server{
		tokenInspector:        tokenInspector,
		authHandler:           authHandler,
		userRepo:              userRepo,
		sessionRepo:           sessionRepo,
		attendanceFormHandler: attendanceFormHandler,
		attendanceRepo:        attendanceRepo,
		formTimeLocation:      formTimeLocation,
	}
}

func (s *Server) SignIn(token string) (string, error) {
	userIdentifier, err := s.tokenInspector.GetUserIdentifier(token)
	if err != nil {
		return "", newBadRequestError(fmt.Errorf("failed to get user identifier: %w", err))
	}

	// TODO(#67): Distinguish errors.
	email, ok := userIdentifier.Email(s.tokenInspector.Provider())
	if !ok {
		return "", newInternalServerError(errors.New("failed to get email from user identifier although there should be"))
	}

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", newNotFoundError(fmt.Errorf("failed to get user by email: %w", err))
	}

	rushToken, err := s.authHandler.SignIn(
		auth.NewUserIdentifier(
			map[auth.Provider]string{auth.ProviderRush: user.Id},
			map[auth.Provider]string{auth.ProviderRush: email},
		),
	)
	if err != nil {
		return "", newInternalServerError(fmt.Errorf("failed to sign in: %w", err))
	}

	return rushToken, nil
}

func (s *Server) IsTokenValid(token string) bool {
	if _, err := s.authHandler.GetUserIdentifier(token); err != nil {
		return false
	}
	return true
}

func (s *Server) GetUserIdentifier(token string) (string, error) {
	userIdentifier, err := s.authHandler.GetUserIdentifier(token)
	if err != nil {
		return "", newBadRequestError(fmt.Errorf("failed to get user identifier: %w", err))
	}

	userId, ok := userIdentifier.ProviderId(auth.ProviderRush)
	if !ok {
		return "", newInternalServerError(errors.New("failed to get user ID from user identifier although there should be"))
	}

	return userId, nil
}

func (s *Server) GetAllUsers() ([]*User, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, newInternalServerError(fmt.Errorf("failed to get users: %w", err))
	}

	converted := []*User{}
	for _, user := range users {
		converted = append(converted, fromUser(&user))
	}
	return converted, nil
}

type ListUsersResult struct {
	Users      []User `json:"users"`
	IsEnd      bool   `json:"is_end"`
	TotalCount int    `json:"total_count"`
}

func (s *Server) ListUsers(offset int, pageSize int) (*ListUsersResult, error) {
	listResult, err := s.userRepo.List(offset, pageSize)
	if err != nil {
		return nil, newInternalServerError(fmt.Errorf("failed to list users: %w", err))
	}

	converted := []User{}
	for _, user := range listResult.Users {
		converted = append(converted, *fromUser(&user))
	}

	return &ListUsersResult{
		Users:      converted,
		IsEnd:      listResult.IsEnd,
		TotalCount: listResult.TotalCount,
	}, nil
}

func (s *Server) GetUser(id string) (*User, error) {
	user, err := s.userRepo.Get(id)
	if err != nil {
		return nil, newNotFoundError(fmt.Errorf("failed to get user: %w", err))
	}
	return fromUser(user), nil
}

func (s *Server) AddUser(name string, university string, phone string, generation float64, isActive bool) error {
	err := s.userRepo.Add(&user.User{
		Name:       name,
		University: university,
		Phone:      phone,
		Generation: generation,
		IsActive:   isActive,
	})
	if err != nil {
		return newInternalServerError(fmt.Errorf("failed to add user: %w", err))
	}
	return nil
}

func (s *Server) GetSession(id string) (*Session, error) {
	session, err := s.sessionRepo.Get(id)
	if err != nil {
		return nil, newNotFoundError(fmt.Errorf("failed to get session: %w", err))
	}
	return fromSession(session), nil
}

type ListSessionsResult struct {
	Sessions   []Session `json:"sessions"`
	IsEnd      bool      `json:"is_end"`
	TotalCount int       `json:"total_count"`
}

func (s *Server) ListSessions(offset int, pageSize int) (*ListSessionsResult, error) {
	listResult, err := s.sessionRepo.List(offset, pageSize)
	if err != nil {
		return nil, newInternalServerError(fmt.Errorf("failed to list sessions: %w", err))
	}

	converted := []Session{}
	for _, session := range listResult.Sessions {
		converted = append(converted, *fromSession(&session))
	}

	return &ListSessionsResult{
		Sessions:   converted,
		IsEnd:      listResult.IsEnd,
		TotalCount: listResult.TotalCount,
	}, nil
}

func (s *Server) CreateSessionForm(sessionId string) (string, error) {
	dbSession, err := s.sessionRepo.Get(sessionId)
	if err != nil {
		return "", newNotFoundError(fmt.Errorf("failed to get session: %w", err))
	}
	if dbSession.IsClosed {
		return "", newBadRequestError(errors.New("session is already closed"))
	}

	users, err := s.userRepo.GetAll()
	if err != nil {
		return "", newInternalServerError(fmt.Errorf("failed to get users: %w", err))
	}

	sort.Slice(users, func(i, j int) bool {
		if users[i].Generation != users[j].Generation {
			return users[i].Generation > users[j].Generation
		}
		return users[i].Name < users[j].Name
	})

	if dbSession.GoogleFormUri != "" {
		return "", newBadRequestError(fmt.Errorf("form already exists: URI is %s", dbSession.GoogleFormUri))
	}

	formTitle := fmt.Sprintf("[출석] %s", dbSession.Name)
	startsAt := dbSession.StartsAt.In(s.formTimeLocation)
	expiresAt := startsAt.Add(-time.Second)
	formDescription := fmt.Sprintf(`%s을(를) 위한 출석용 구글폼입니다.
폼 마감 시각은 %s입니다. %s 이후 요청은 무시됩니다.`, dbSession.Name, expiresAt.Format("2006-01-02 15:04:05"), startsAt.Format("2006-01-02 15:04:05"))

	attendanceForm, err := s.attendanceFormHandler.GenerateForm(formTitle, formDescription, users)
	if err != nil {
		return "", newInternalServerError(fmt.Errorf("failed to generate form: %w", err))
	}

	_, err = s.sessionRepo.Update(sessionId, &session.UpdateForm{
		GoogleFormId: &attendanceForm.Id, GoogleFormUri: &attendanceForm.Uri, ReturnUpdatedSession: false})
	if err != nil {
		return "", newInternalServerError(fmt.Errorf("failed to update session: %w", err))
	}

	return attendanceForm.Uri, nil
}

func (s *Server) AddSession(name string, description string, startsAt time.Time, score int) (string, error) {
	id, err := s.sessionRepo.Add(name, description, 0, 0, startsAt, score)
	if err != nil {
		return "", newInternalServerError(fmt.Errorf("failed to add session: %w", err))
	}
	return id, nil
}

func (s *Server) CloseSession(sessionId string) error {
	dbSession, err := s.sessionRepo.Get(sessionId)
	if err != nil {
		return newNotFoundError(fmt.Errorf("failed to get session: %w", err))
	}

	if dbSession.IsClosed {
		return newBadRequestError(errors.New("session is already closed"))
	}

	formSubmissions, err := s.attendanceFormHandler.GetSubmissions(dbSession.GoogleFormId)
	if err != nil {
		return newInternalServerError(fmt.Errorf("failed to get form submissions: %w", err))
	}

	submissionsOnTime := []attendance.FormSubmission{}
	for _, submission := range formSubmissions {
		if submission.SubmissionTime.Before(dbSession.StartsAt) {
			submissionsOnTime = append(submissionsOnTime, submission)
		}
	}

	externalIds := []string{}
	for _, submissionOnTime := range submissionsOnTime {
		externalIds = append(externalIds, submissionOnTime.UserExternalId)
	}

	users, err := s.userRepo.GetAllByExternalIds(externalIds)
	if err != nil {
		return newInternalServerError(fmt.Errorf("failed to get users by external IDs: %w", err))
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].ExternalId < users[j].ExternalId
	})
	sort.Slice(submissionsOnTime, func(i, j int) bool {
		return submissionsOnTime[i].UserExternalId < submissionsOnTime[j].UserExternalId
	})

	addAttendanceReqs := []attendance.AddAttendanceReq{}
	for index, submissionOnTime := range submissionsOnTime {
		user := users[index]
		addAttendanceReqs = append(addAttendanceReqs, attendance.AddAttendanceReq{
			SessionId:   sessionId,
			SessionName: dbSession.Name,
			UserId:      user.Id,
			UserName:    user.Name,
			JoinedAt:    submissionOnTime.SubmissionTime,
		})
	}

	if err := s.attendanceRepo.BulkInsert(addAttendanceReqs); err != nil {
		return newInternalServerError(fmt.Errorf("failed to bulk insert attendance: %w", err))
	}

	isClosed := true
	_, err = s.sessionRepo.Update(sessionId, &session.UpdateForm{IsClosed: &isClosed, ReturnUpdatedSession: false})

	if err != nil {
		return newInternalServerError(fmt.Errorf("failed to update session to be closed: %w", err))
	}

	return nil
}

func (s *Server) GetAttendanceByUserId(userId string) ([]Attendance, error) {
	attendances, err := s.attendanceRepo.FindByUserId(userId)
	if err != nil {
		return nil, newInternalServerError(fmt.Errorf("failed to find attendance by user ID: %w", err))
	}
	converted := []Attendance{}
	for _, attendance := range attendances {
		converted = append(converted, *fromAttendance(&attendance))
	}
	return converted, nil
}
