package auth

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
}

func NewUserIdentifier(providerIds map[Provider]string, emails map[Provider]string) UserIdentifier {
	return UserIdentifier{providerIds: providerIds, emails: emails}
}

func (u *UserIdentifier) Email(provider Provider) (string, bool) {
	email, ok := u.emails[provider]
	return email, ok
}

func (u *UserIdentifier) ProviderId(provider Provider) (string, bool) {
	id, ok := u.providerIds[provider]
	return id, ok
}
