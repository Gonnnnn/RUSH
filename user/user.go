package user

import "rush/permission"

type User struct {
	// The unique identifier of the user.
	Id string `json:"id"`
	// The name of the user. E.g., "김건"
	Name string `json:"name"`
	// The role of the user. E.g., "member"
	Role permission.Role `json:"role"`
	// The generation of the user. E.g., 9
	Generation float64 `json:"generation"`
	// TODO(#223): Fix the repo to handle active users only by default.
	// The activity status of the user. E.g., true
	IsActive bool `json:"is_active"`
	// The email address of the user. E.g., "kim.geon@gmail.com"
	Email string `json:"email"`
	// The unique name consisting of the user name and a number.
	// It's used as an external ID for the users so that
	// it's easier for them to identify themselves such as in Google Forms.
	ExternalName string `json:"external_name"`
}
