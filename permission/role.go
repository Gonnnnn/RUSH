// Global package for user permission on the attendance system.
package permission

type Role string

const (
	RoleNotSpecified Role = ""
	RoleUnknown      Role = "unknown"
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
