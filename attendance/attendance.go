package attendance

import (
	"time"
)

type AttendanceReport struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	SessionIds  []string  `json:"session_ids"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   int       `json:"created_by"`
}
