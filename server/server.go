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
	// The ID of the user. E.g., "abc123"
	Id string `json:"id"`
	// The name of the user. E.g., "김건"
	// If the user name is duplicated, it might have a suffix. E.g., "김건B"
	Name string `json:"name"`
	// The university of the user. E.g., "Yonsei"
	University string `json:"university"`
	// The phone number of the user. E.g., "010-1234-5678"
	// The format is not strictly validated.
	Phone string `json:"phone"`
	// The generation of the user. E.g., 9.5
	// It's either has 0 or 5 for the decimal part.
	Generation float64 `json:"generation"`
	// The activity status of the user. E.g., true
	IsActive bool `json:"is_active"`
}

type Session struct {
	// The ID of the session. E.g., "abc123"
	Id string `json:"id"`
	// The name of the session. E.g., "456회 정기 세션"
	Name string `json:"name"`
	// The description of the session. E.g., "연대 트랙, ..."
	Description string `json:"description"`
	// The ID of the user who hosts the session. E.g., "abc123"
	HostedBy string `json:"hosted_by"`
	// The ID of the user who created the session. E.g., "abc123"
	CreatedBy string `json:"created_by"`
	// The URI of the Google form for the session. E.g., "https://docs.google.com/forms/d/e/1FAIpQLSd..."
	GoogleFormUri string `json:"google_form_uri"`
	// The external IDs of the users who joined the session. E.g., ["abc123", "def456"]
	JoinningUsers []string `json:"joinning_users"`
	// The time in UTC when the session is created.
	CreatedAt time.Time `json:"created_at"`
	// The time in UTC when the session starts.
	StartsAt time.Time `json:"starts_at"`
	// The score of the session. E.g., 100
	Score int `json:"score"`
	// The indicator if the session is closed. E.g., true
	IsClosed bool `json:"is_closed"`
}

type Attendance struct {
	// The ID of the attendance. E.g., "abc123"
	Id string `json:"id"`
	// The ID of the session that the attendance is related to. E.g., "abc123"
	SessionId string `json:"session_id"`
	// The name of the session that the attendance is related to. E.g., "456회 정기 세션"
	// It's not synced with the actual session data.
	SessionName string `json:"session_name"`
	// The ID of the user who joined the session. E.g., "abc123"
	UserId string `json:"user_id"`
	// The name of the user who joined the session. E.g., "김건"
	// It's not synced with the actual user data.
	UserName string `json:"user_name"`
	// The time in UTC when the user joined the session.
	JoinedAt time.Time `json:"joined_at"`
	// The time in UTC when the attendance is created.
	CreatedAt time.Time `json:"created_at"`
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
	// Returns the user that has the email. Typically used to get the user by the email from the OAuth2.0 token.
	GetByEmail(email string) (*user.User, error)
	// Returns the users that have the external IDs. Typically used to get users by the external IDs from the form.
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
	// Generates a form with the title, description, and users for attendance.
	GenerateForm(title string, description string, users []user.User) (attendance.Form, error)
	// Extracts the submissions submitted to the form by the users.
	GetSubmissions(formId string) ([]attendance.FormSubmission, error)
}

type attendanceRepo interface {
	// Inserts the attendance requests in bulk. It's used to insert the attendance requests after closing the session.
	BulkInsert(requests []attendance.AddAttendanceReq) error
	// Returns the attendances that are related to the user. Typically used to get the attendances for each user.
	FindByUserId(userId string) ([]attendance.Attendance, error)
}

type Server struct {
	// Used to get the user identifier, such as an email, from the third party token.
	tokenInspector tokenInspector
	// Used to sign in and get the rush token for API calls.
	authHandler authHandler
	userRepo    userRepo
	sessionRepo sessionRepo
	// Used to generate the form for attendance and get the submissions from the form.
	attendanceFormHandler attendanceFormHandler
	attendanceRepo        attendanceRepo
	// The location of the time for the form. It's used to convert the time in the form to the local time.
	formTimeLocation *time.Location
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
