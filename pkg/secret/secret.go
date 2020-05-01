package secret

import (
	"log"
	"time"
)

var (
	factories map[string]ProviderFactory
	providers map[string]Provider
)

func init() {
	factories = make(map[string]ProviderFactory)
	providers = make(map[string]Provider)
}

// RegisterProvider registers a new provider in the secrets system.
func RegisterProvider(n string, p ProviderFactory) {
	if _, ok := factories[n]; ok {
		// Already registered
		return
	}
	factories[n] = p
}

// InitializeProviders handles the setup of all providers during
// startup, which allows the actual connections of some providers to
// be deferred until after all have been registered.
func InitializeProviders() {
	for pn, pf := range factories {
		p, err := pf()
		if err != nil {
			log.Println(err)
			continue
		}
		providers[pn] = p
	}
}

// Poll polls for a secret.  It will only return an error in the case
// of an unrecoverable transport error or an issue with credentialing.
func Poll(pn, id string) (string, error) {
	p, ok := providers[pn]
	if !ok {
		return "", ErrNoSuchProvider
	}

	for {
		s, err := p.FetchSecret(id)
		switch err {
		case nil:
			return s, nil
		case ErrTerminal:
			return "", err
		default:
			time.Sleep(time.Second * 10)
		}
	}
	return p.FetchSecret(id)
}
