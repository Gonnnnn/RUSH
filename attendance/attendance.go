package attendance

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type AttendanceReport struct {
	Id          int         `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	SessionIds  integerList `json:"session_ids"`
	CreatedAt   time.Time   `json:"created_at"`
	CreatedBy   int         `json:"created_by"`
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
