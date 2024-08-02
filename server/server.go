package server

import (
	"fmt"
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

type userRepo interface {
	GetAll() ([]user.User, error)
	// Skips `offset` users and returns up to `pageSize` users, an indicator if it has more users and total count.
	List(offset int, pageSize int) (*user.ListResult, error)
	Add(user *user.User) error
}

type sessionRepo interface {
	Get(id string) (*session.Session, error)
	GetAll() ([]session.Session, error)
	List(offset int, pageSize int) (*session.ListResult, error)
	Add(name string, description string, hostedBy int, createdBy int, startsAt time.Time, score int) (string, error)
	Update(id string, updateForm *session.UpdateForm) (*session.Session, error)
}

type sessionFormHandler interface {
	GenerateForm(title string, description string, users []user.User) (string, error)
	ReadUsers(formId string) ([]string, error)
}

type Server struct {
	userRepo           userRepo
	sessionRepo        sessionRepo
	sessionFormHandler sessionFormHandler
}

func New(userRepo userRepo, sessionRepo sessionRepo, sessionFormHandler sessionFormHandler) *Server {
	return &Server{
		userRepo:           userRepo,
		sessionRepo:        sessionRepo,
		sessionFormHandler: sessionFormHandler,
	}
}

func (s *Server) GetAllUsers() ([]*User, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
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
		return nil, err
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

func (s *Server) AddUser(name string, university string, phone string, generation float64, isActive bool) error {
	return s.userRepo.Add(&user.User{
		Name:       name,
		University: university,
		Phone:      phone,
		Generation: generation,
		IsActive:   isActive,
	})
}

func (s *Server) GetSession(id string) (*Session, error) {
	session, err := s.sessionRepo.Get(id)
	if err != nil {
		return nil, err
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
		return nil, err
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
	users, err := s.userRepo.GetAll()
	if err != nil {
		return "", fmt.Errorf("failed to get users: %w", err)
	}

	sort.Slice(users, func(i, j int) bool {
		if users[i].Generation != users[j].Generation {
			return users[i].Generation > users[j].Generation
		}
		return users[i].Name < users[j].Name
	})

	dbSession, err := s.sessionRepo.Get(sessionId)
	if err != nil {
		return "", fmt.Errorf("failed to get session: %w", err)
	}

	if dbSession.GoogleFormUri != "" {
		return "", fmt.Errorf("form already created: %s", dbSession.GoogleFormUri)
	}

	formTitle := fmt.Sprintf("[출석] %s", dbSession.Name)
	startsAt := dbSession.StartsAt
	expiresAt := startsAt.Add(-time.Second)
	formDescription := fmt.Sprintf(`%s을(를) 위한 출석용 구글폼입니다.
폼 마감 시각은 %s입니다. %s 이후 요청은 무시됩니다.`, dbSession.Name, expiresAt.Format("2006-01-02 15:04:05"), startsAt.Format("2006-01-02 15:04:05"))

	formUri, err := s.sessionFormHandler.GenerateForm(formTitle, formDescription, users)
	if err != nil {
		return "", fmt.Errorf("failed to generate form: %w", err)
	}

	_, err = s.sessionRepo.Update(sessionId, &session.UpdateForm{GoogleFormUri: &formUri, ReturnUpdatedSession: false})
	if err != nil {
		return "", fmt.Errorf("failed to update session: %w", err)
	}

	return formUri, nil
}

func (s *Server) AddSession(name string, description string, startsAt time.Time, score int) (string, error) {
	return s.sessionRepo.Add(name, description, 0, 0, startsAt, score)
}
