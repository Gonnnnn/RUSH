package attendance

import "time"

type Attendance struct {
	// The unique identifier for the attendance record. E.g. "1"
	Id string `json:"id"`
	// The unique identifier for the session that the user joined. E.g. "1"
	SessionId string `json:"session_id"`
	// The name of the session that the user joined. E.g. "Yonsei University track"
	SessionName string `json:"session_name"`
	// The unique identifier for the user. E.g. "1"
	UserId string `json:"user_id"`
	// The name of the user. E.g. "Alice"
	UserName string `json:"user_name"`
	// The time when the user joined the session.
	JoinedAt time.Time `json:"joined_at"`
	// The time when the attendance record was created.
	CreatedAt time.Time `json:"created_at"`
}
