package user

import "rush/permission"

type User struct {
	Id         string          `json:"id"`
	Name       string          `json:"name"`
	Role       permission.Role `json:"role"`
	Generation float64         `json:"generation"`
	IsActive   bool            `json:"is_active"`
	Email      string          `json:"email"`
	// The unique name consisting of the user name and a number.
	// It's used as an external ID for the users so that
	// it's easier for them to identify themselves such as in Google Forms.
	ExternalName string `json:"external_name"`
}
