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
	ProviderIds map[Provider]string
	// The email address of the user provided by the provider.
	// E.g., Firebase: john.doe@gmail.com
	Emails map[Provider]string
}
