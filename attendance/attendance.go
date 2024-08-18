package attendance

import "time"

type Attendance struct {
	// The unique identifier for the attendance record. E.g. "1"
	Id string `json:"id"`
	// The unique identifier for the session that the user joined. E.g. "1"
	SessionId string `json:"session_id"`
	// The name of the session that the user joined. E.g. "Yonsei University track"
	SessionName string `json:"session_name"`
	// The score of the session. E.g. 2
	SessionScore int `json:"session_score"`
	// The time when the session started.
	SessionStartedAt time.Time `json:"session_started_at"`
	// The unique identifier for the user. E.g. "1"
	UserId string `json:"user_id"`
	// The name of the user. E.g. "Alice"
	UserExternalName string `json:"user_external_name"`
	// The generation of the user. E.g. 9.5
	UserGeneration float64 `json:"user_generation"`
	// The time when the user joined the session.
	UserJoinedAt time.Time `json:"user_joined_at"`
	// The time when the attendance record was created.
	CreatedAt time.Time `json:"created_at"`
}
