package server

import (
	"rush/attendance"
	"rush/auth"
	"rush/session"
	"rush/user"
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
