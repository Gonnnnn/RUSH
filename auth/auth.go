package auth

import "rush/permission"

type inspector struct{}

type Provider string

const (
	// Unknown provider. There should be something wrong.
	ProviderUnknown Provider = "unknown"
	// Rush itself.
	ProviderRush Provider = "rush"
	// Firebase Auth. https://firebase.google.com/docs/auth
	ProviderFirebase Provider = "firebase"
)

type UserIdentifier struct {
	// THe user ID provided by the provider.
	// Other packages could use it to identify the user matching the provider.
	providerIds map[Provider]string
	// The email address of the user provided by the provider.
	// E.g., Firebase: john.doe@gmail.com
	emails map[Provider]string
	// The role of the user. It is used to determine the access level of the user.
	// E.g., member, admin, etc.
	role map[Provider]permission.Role
}

// TODO(#138): Identifier could be implemented for each provider and have one common interface.
func NewUserIdentifier(providerIds map[Provider]string, emails map[Provider]string, role map[Provider]permission.Role) UserIdentifier {
	if providerIds == nil {
		providerIds = make(map[Provider]string)
	}
	if emails == nil {
		emails = make(map[Provider]string)
	}
	if role == nil {
		role = make(map[Provider]permission.Role)
	}
	return UserIdentifier{providerIds: providerIds, emails: emails, role: role}
}

func (u *UserIdentifier) Email(provider Provider) (string, bool) {
	email, ok := u.emails[provider]
	return email, ok
}

func (u *UserIdentifier) ProviderId(provider Provider) (string, bool) {
	id, ok := u.providerIds[provider]
	return id, ok
}

func (u *UserIdentifier) RushRole() permission.Role {
	role, ok := u.role[ProviderRush]
	if !ok {
		return permission.RoleNotSpecified
	}
	return role
}
