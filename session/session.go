package session

import (
	"time"
)

// Session represents a session of a running event. It has data to identify the session,
// and the score. Plus, it has data about its attendance.
// session attendance can be applied by Google form submissions or manually.
// Once it's applied, the session data is immutable. Whether it's applied is indicated by `AttendanceStatus`.
type Session struct {
	// The ID of the session. It's a unique identifier. E.g., "abc123"
	Id string `json:"id"`
	// The name of the session. E.g., "456회 ..."
	Name string `json:"name"`
	// The description of the session. E.g., "연대 트랙..."
	Description string `json:"description"`
	// The ID of the user who created the session. E.g., "abc123"
	CreatedBy string `json:"created_by"`
	// The Google form ID for attendance of the session. E.g., "abc123"
	GoogleFormId string `json:"google_form_id"`
	// The Google form URI for attendance of the session. E.g., "https://forms.gle/abc123"
	GoogleFormUri string `json:"google_form_uri"`
	// The time in UTC when the session was created.
	CreatedAt time.Time `json:"created_at"`
	// The time in UTC when the session starts.
	StartsAt time.Time `json:"starts_at"`
	// The attendance score of the session that the user can get. E.g., 2
	Score int `json:"score"`
	// The status of the session's attendance.
	// It indicates if it is applied, ignored, etc.
	AttendanceStatus AttendanceStatus `json:"attendance_status"`
}

type AttendanceStatus string

const (
	// not applied yet.
	AttendanceStatusNotAppliedYet AttendanceStatus = "not_applied_yet"
	// The attendance has been applied. Once it's applied, the session data is immutable.
	AttendanceStatusApplied AttendanceStatus = "applied"
	// It has been tried to apply the attendance but ignored for some reasons. It should be checked manually.
	AttendanceStatusIgnored AttendanceStatus = "ignored"
)

// If the session attendance is applied, the session data is immutable. It checks if the data can be updated.
func (s *Session) CanUpdateMetadata() bool {
	return s.AttendanceStatus == AttendanceStatusNotAppliedYet || s.AttendanceStatus == AttendanceStatusIgnored
}

// Checks the session data and returns true if the attendance can be applied by Google form submissions.
func (s *Session) CanApplyAttendanceByFormSubmissions() bool {
	if s.AttendanceStatus == AttendanceStatusApplied {
		return false
	}

	if s.GoogleFormId == "" || s.GoogleFormUri == "" {
		return false
	}

	return true
}

// Checks the session data and returns true if the attendance can be applied manually.
func (s *Session) CanApplyAttendanceManually() bool {
	if s.AttendanceStatus == AttendanceStatusApplied {
		return false
	}

	if s.GoogleFormId != "" || s.GoogleFormUri != "" {
		return false
	}

	return true
}
