package session

import (
	"time"
)

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
	// If the session is closed, no one can fix the metadata besides the developer.
	// It's to prevent cheating. E.g., false
	IsClosed bool `json:"is_closed"`
}
