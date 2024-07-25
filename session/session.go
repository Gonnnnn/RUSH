package session

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Session struct {
	Id            string      `json:"id"`
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	HostedBy      int         `json:"hosted_by"`
	CreatedBy     int         `json:"created_by"`
	JoinningUsers integerList `json:"joinning_users"`
	CreatedAt     time.Time   `json:"created_at"`
	StartsAt      time.Time   `json:"starts_at"`
	Score         int         `json:"score"`
	// If the session is closed, no one can fix the metadata. It's to prevent cheating.
	IsClosed bool `json:"is_closed"`
}

type integerList string

func (i *integerList) UnmarshalJSON(data []byte) error {
	stringified := string(data)
	for _, element := range strings.Split(stringified, ",") {
		if _, err := strconv.Atoi(element); err != nil {
			return fmt.Errorf("non integer included: %s", element)
		}
	}
	*i = integerList(stringified)
	return nil
}

func (i integerList) MarshalJSON() ([]byte, error) {
	split := strings.Split(string(i), ",")
	for _, element := range split {
		if _, err := strconv.Atoi(element); err != nil {
			return nil, fmt.Errorf("non integer included: %s", element)
		}
	}
	return []byte(string(i)), nil
}
