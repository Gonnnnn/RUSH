package user

type User struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	Role       Role    `json:"role"`
	Generation float64 `json:"generation"`
	IsActive   bool    `json:"is_active"`
	Email      string  `json:"email"`
	// The unique name consisting of the user name and a number.
	// It's used as an external ID for the users so that
	// it's easier for them to identify themselves such as in Google Forms.
	ExternalName string `json:"external_name"`
}

type Role string

const (
	// Super admin has the highest authority.
	// It can perform any APIs. It is for the developers or the one who has
	// the ownership of the attendance system.
	RoleSuperAdmin Role = "super_admin"
	// Admin is the second highest authority.
	// It can perform certain APIs such as creating a session, updating a session, etc.
	RoleAdmin Role = "admin"
	// Member is the lowest authority. They can only view certain data
	// or access their own data.
	RoleMember Role = "member"
)
