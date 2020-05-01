package secret

import "errors"

var (
	// ErrNoSuchProvider is returned when a provider has been
	// requested that does not exist within the system.
	ErrNoSuchProvider = errors.New("no provider with that name exists")

	// ErrTerminal is returned when the fetch should not be
	// retried for a given secret.
	ErrTerminal = errors.New("an unrecoverable error occured during fetch")

	// ErrNotFound is returned when the fetch is able to connect
	// to the secret service, but no secret is present.
	ErrNotFound = errors.New("specified secret was not found")
)
