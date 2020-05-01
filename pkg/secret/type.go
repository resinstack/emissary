package secret

// A Provider loads a secret from a remote secret storage
// server.  Retry logic is handled at a higher level in the secret
// package.
type Provider interface {
	FetchSecret(string) (string, error)
}

// A ProviderFactory initializes and returns a Provider initialized
// and connected if applicable.
type ProviderFactory func() (Provider, error)
