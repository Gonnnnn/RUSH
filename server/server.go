package server

import (
	"rush/attendance"
	"rush/auth"
	"rush/permission"
	"rush/session"
	"rush/user"
	"time"

	"github.com/benbjohnson/clock"
)

//go:generate mockgen -source=server.go -destination=mock/server_mock.go -package=mock

type User struct {
	// The ID of the user. E.g., "abc123"
	Id string `json:"id"`
	// The name of the user. E.g., "김건"
	Name string `json:"name"`
	// The generation of the user. E.g., 9.5
	// It's either has 0 or 5 for the decimal part.
	Generation float64 `json:"generation"`
	// The activity status of the user. E.g., true
	IsActive bool `json:"is_active"`
	// The email address of the user. E.g., "kim.geon@gmail.com"
	Email string `json:"email"`
	// The external name of the user. E.g., "김건3"
	// It's used as an external ID for the users so that it's easier for them to identify themselves such as in Google Forms.
	ExternalName string `json:"external_name"`
}

type SessionForAdmin struct {
	// The ID of the session. E.g., "abc123"
	Id string `json:"id"`
	// The name of the session. E.g., "456회 정기 세션"
	Name string `json:"name"`
	// The description of the session. E.g., "연대 트랙, ..."
	Description string `json:"description"`
	// The ID of the user who created the session. E.g., "abc123"
	CreatedBy string `json:"created_by"`
	// The URI of the Google form for the session. E.g., "https://docs.google.com/forms/d/e/1FAIpQLSd..."
	GoogleFormUri string `json:"google_form_uri"`
	// The ID of the Google form for the session. E.g., "1FAIpQLSd..."
	GoogleFormId string `json:"google_form_id"`
	// The time in UTC when the session is created.
	CreatedAt time.Time `json:"created_at"`
	// The time in UTC when the session starts.
	StartsAt time.Time `json:"starts_at"`
	// The score of the session. E.g., 100
	Score int `json:"score"`
	// The status of the session's attendance.
	// It indicates if it is applied, ignored, etc.
	AttendanceStatus session.AttendanceStatus `json:"attendance_status"`
	// The flag to indicate how the attendance is applied. E.g., "manual" or "form".
	AttendanceAppliedBy SessionAttendanceAppliedBy `json:"attendance_applied_by"`
}

// Session for a user. It includes fields that are safe for a user to know.
type Session struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	StartsAt    time.Time `json:"starts_at"`
	Score       int       `json:"score"`
}

type SessionAttendanceAppliedBy string

const (
	// Literally unknown. It's something but that is not known yet. Maybe DB has been updated but
	// the code is not updated.
	SessionAttendanceAppliedByUnknown SessionAttendanceAppliedBy = "unknown"
	// Not specified yet. It's used for the initial state where the attendance is not applied yet.
	SessionAttendanceAppliedByUnspecified SessionAttendanceAppliedBy = "unspecified"
	// The attendance is applied manually.
	SessionAttendanceAppliedByManual SessionAttendanceAppliedBy = "manual"
	// The attendance is applied by the form submissions.
	SessionAttendanceAppliedByForm SessionAttendanceAppliedBy = "form"
)

type Attendance struct {
	// The ID of the attendance. E.g., "abc123"
	Id string `json:"id"`
	// The ID of the session that the attendance is related to. E.g., "abc123"
	SessionId string `json:"session_id"`
	// The name of the session that the attendance is related to. E.g., "456회 정기 세션"
	// It's not synced with the actual session data.
	SessionName string `json:"session_name"`
	// The score of the session. E.g., 1
	SessionScore int `json:"session_score"`
	// The time when the session started.
	SessionStartedAt time.Time `json:"session_started_at"`
	// The ID of the user who joined the session. E.g., "abc123"
	UserId string `json:"user_id"`
	// The name of the user who joined the session. E.g., "김건"
	// It's not synced with the actual user data. It's used to show the user name in the attendance list.
	// Thus the external name is used instead of the user name.
	UserExternalName string `json:"user_external_name"`
	// The generation of the user who joined the session. E.g., 9.5
	UserGeneration float64 `json:"user_generation"`
	// The time in UTC when the user joined the session.
	UserJoinedAt time.Time `json:"user_joined_at"`
	// The time in UTC when the attendance is created.
	CreatedAt time.Time `json:"created_at"`
}

// The API request session. It contains the user information and some more to
// specify the session for the API request.
type UserSession struct {
	UserId    string          `json:"user_id"`
	Role      permission.Role `json:"role"`
	ExpiresAt time.Time       `json:"expires_at"`
}

type oauthClient interface {
	// Handles the third party token that is used for signing in.
	GetEmail(token string) (string, error)
}

type authHandler interface {
	// Extracts session from the rush token.
	GetSession(token string) (auth.Session, error)
	// Returns the rush token that is used for API calls after signing in.
	SignIn(userId string, role permission.Role) (string, error)
}

type userRepo interface {
	Get(id string) (*user.User, error)
	// Returns all the users.
	GetAll() ([]user.User, error)
	// Returns all the active users.
	GetAllActive() ([]user.User, error)
	// Skips `offset` users and returns up to `pageSize` users, an indicator if it has more users and total count.
	List(offset int, pageSize int) (*user.ListResult, error)
	// Returns the user that has the email. Typically used to get the user by the email from the OAuth2.0 token.
	GetByEmail(email string) (*user.User, error)
	// Returns the users that have the external names. Typically used to get users by the external names from the form.
	GetAllByExternalNames(externalNames []string) ([]user.User, error)
}

type userAdder interface {
	Add(name string, generation float64, isActive bool, email string) error
}

type userUpdater interface {
	// Updates the user with the given ID. It should include the logics to
	// be executed when updating a user so that all other data can be updated.
	Update(id string, updateForm user.UpdateForm) error
}

type sessionRepo interface {
	Get(id string) (session.Session, error)
	GetAll() ([]session.Session, error)
	List(offset int, pageSize int) (*session.ListResult, error)
	Add(name string, description string, createdBy string, startsAt time.Time, score int) (string, error)
}

// The repo that includes logics to update or delete the open sessions.
type openSessionRepo interface {
	UpdateOpenSession(id string, updateForm session.OpenSessionUpdateForm) (session.Session, error)
	DeleteOpenSession(id string) error
	// Marks the attendance status of the open session to be applied.
	// Use it after inserting attendances for the open session.
	MarkAsAttendanceApplied(id string) error
	// Marks the attendance status of the open session to be ignored.
	// Use it when the session's attendance or anything about session is not correct,
	// or suspicious, so that the server decides to not apply the attendance.
	MarkAttendanceIsIgnored(id string, reason string) error
}

type attendanceFormHandler interface {
	// Generates a form with the title, description, and user external names/generations for attendance.
	GenerateForm(title string, description string, userOptions []attendance.UserOption) (attendance.Form, error)
	// Extracts the submissions submitted to the form by the users.
	GetSubmissions(formId string) ([]attendance.FormSubmission, error)
}

type attendanceRepo interface {
	// Returns all the attendance requests. It is used to provide admins with the attendance result of all users.
	GetAll() ([]attendance.Attendance, error)
	// Inserts the attendance requests in bulk. It's used to insert the attendance requests after closing the session.
	BulkInsert(requests []attendance.AddAttendanceReq) error
	// Returns the attendances that are related to the user. Typically used to get the attendances for each user.
	FindByUserId(userId string) ([]attendance.Attendance, error)
	// Returns the attendances that are related to the session. Typically used for admins to see if attendance is applied well.
	FindBySessionId(sessionId string) ([]attendance.Attendance, error)
}

type Server struct {
	// Used to get the user email of the provider from the third party token.
	oauthClient oauthClient
	// Used to sign in and get the rush token for API calls.
	authHandler authHandler
	userRepo    userRepo
	// Used to add a user.
	userAdder userAdder
	// Used to update a user. It has logics that should be handled when updating a user.
	userUpdater userUpdater
	sessionRepo sessionRepo
	// Used to handle open sessions.
	openSessionRepo openSessionRepo
	// Used to generate the form for attendance and get the submissions from the form.
	attendanceFormHandler attendanceFormHandler
	attendanceRepo        attendanceRepo
	// The location of the time for the form. It's used to convert the time in the form to the local time.
	formTimeLocation *time.Location
	// Used to get the current time.
	clock clock.Clock
}

func New(oauthClient oauthClient, authHandler authHandler, userRepo userRepo, userAdder userAdder, userUpdater userUpdater, sessionRepo sessionRepo, openSessionRepo openSessionRepo,
	attendanceFormHandler attendanceFormHandler, attendanceRepo attendanceRepo, formTimeLocation *time.Location, clock clock.Clock) *Server {
	return &Server{
		oauthClient:           oauthClient,
		authHandler:           authHandler,
		userRepo:              userRepo,
		userAdder:             userAdder,
		userUpdater:           userUpdater,
		sessionRepo:           sessionRepo,
		openSessionRepo:       openSessionRepo,
		attendanceFormHandler: attendanceFormHandler,
		attendanceRepo:        attendanceRepo,
		formTimeLocation:      formTimeLocation,
		clock:                 clock,
	}
}
