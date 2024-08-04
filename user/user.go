package user

type User struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	University string  `json:"university"`
	Phone      string  `json:"phone"`
	Generation float64 `json:"generation"`
	IsActive   bool    `json:"is_active"`
	// The ID of the user to expose to the external system outside of the service
	// such as the Google Form.
	ExternalId string `json:"external_id"`
}
