package user

type User struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	University string  `json:"university"`
	Phone      string  `json:"phone"`
	Generation float64 `json:"generation"`
	IsActive   bool    `json:"is_active"`
	// The unique name consisting of the user name and a number.
	// It's used as an external ID for the users so that
	// it's easier for them to identify themselves such as in Google Forms.
	ExternalName string `json:"external_name"`
}
