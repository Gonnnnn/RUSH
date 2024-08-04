package session

import (
	"time"
)

type Session struct {
	Id            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	HostedBy      int       `json:"hosted_by"`
	CreatedBy     int       `json:"created_by"`
	GoogleFormId  string    `json:"google_form_id"`
	GoogleFormUri string    `json:"google_form_uri"`
	JoinningUsers []string  `json:"joinning_users"`
	CreatedAt     time.Time `json:"created_at"`
	StartsAt      time.Time `json:"starts_at"`
	Score         int       `json:"score"`
	// If the session is closed, no one can fix the metadata. It's to prevent cheating.
	IsClosed bool `json:"is_closed"`
}
