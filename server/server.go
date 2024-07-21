package server

import (
	"rush/attendance"
	"rush/session"
	"rush/user"
	"strings"
	"time"
)

type User struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	University string `json:"university"`
	Phone      string `json:"phone"`
	Generation string `json:"generation"`
	IsActive   bool   `json:"is_active"`
}

type Session struct {
	Id            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	HostedBy      string    `json:"hosted_by"`
	CreatedBy     string    `json:"created_by"`
	JoinningUsers []string  `json:"joinning_users"`
	CreatedAt     time.Time `json:"created_at"`
	StartsAt      time.Time `json:"starts_at"`
	Score         int       `json:"score"`
	IsClosed      bool      `json:"is_closed"`
}

type AttendanceReport struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	SessionIds  []string  `json:"session_ids"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
}

type userRepo interface {
	GetAll() ([]user.User, error)
	Add(user *user.User) error
}

type sessionRepo interface {
	Get(id string) (*session.Session, error)
	GetAll() ([]session.Session, error)
	Add(name string, description string, hostedBy int, createdBy int, startsAt time.Time, score int) (string, error)
}

type attendanceRepo interface {
	GetAll() ([]attendance.AttendanceReport, error)
	Add(name string, description string, sessionIds []string, createdBy int) error
}

type Server struct {
	userRepo       userRepo
	sessionRepo    sessionRepo
	attendanceRepo attendanceRepo
}

func New(userRepo userRepo, sessionRepo sessionRepo, attendanceRepo attendanceRepo) *Server {
	return &Server{
		userRepo:       userRepo,
		sessionRepo:    sessionRepo,
		attendanceRepo: attendanceRepo,
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

func (s *Server) AddUser(name string, university string, phone string, generation string, isActive bool) error {
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

func (s *Server) GetAllSessions() ([]*Session, error) {
	sessions, err := s.sessionRepo.GetAll()
	if err != nil {
		return nil, err
	}

	converted := []*Session{}
	for _, session := range sessions {
		converted = append(converted, fromSession(&session))
	}
	return converted, nil
}

func (s *Server) AddSession(name string, description string, startsAt time.Time, score int) (string, error) {
	return s.sessionRepo.Add(name, description, 0, 0, startsAt, score)
}

func (s *Server) GetAllAttendanceReports() ([]AttendanceReport, error) {
	reports, err := s.attendanceRepo.GetAll()
	if err != nil {
		return nil, err
	}

	converted := []AttendanceReport{}
	for _, report := range reports {
		converted = append(converted, AttendanceReport{
			Id:          string(report.Id),
			Name:        report.Name,
			Description: report.Description,
			SessionIds:  strings.Split(string(report.SessionIds), ","),
			CreatedAt:   report.CreatedAt,
			CreatedBy:   string(report.CreatedBy),
		})
	}

	return converted, nil
}

func (s *Server) AddAttendanceReport(name string, sessionIds []string, createdBy int) error {
	return s.attendanceRepo.Add(name, "", sessionIds, createdBy)
}
