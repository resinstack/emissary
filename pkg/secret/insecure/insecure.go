package insecure

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/ddo/rq"
	"github.com/ddo/rq/client"

	"github.com/resinstack/emissary/pkg/secret"
)

type insecure struct {
	base *url.URL
}

func init() {
	secret.RegisterProvider("insecure", initialize)
}

func initialize() (secret.Provider, error) {
	i := new(insecure)
	base := os.Getenv("INSECURE_BASE")
	if base == "" {
		return nil, errors.New("required variable INSECURE_BASE not set, insecure engine will be unavailable")
	}
	var err error
	i.base, err = url.Parse(base)
	if err != nil {
		return nil, errors.New("INSECURE_BASE does not container a URL, insecure engine will be unavailable")
	}
	return i, nil
}

func (i *insecure) FetchSecret(s string) (string, error) {
	ru, _ := url.Parse(i.base.String())
	ru.Path = path.Join(ru.Path, s)

	data, resp, err := client.Send(rq.Get(ru.String()), true)
	if err != nil {
		// Terminal error, the host could not be contacted
		return "", secret.ErrTerminal
	}
	ds := string(data)
	switch resp.StatusCode {
	case http.StatusOK:
		return strings.TrimSpace(ds), nil
	case http.StatusNotFound:
		return "", secret.ErrNotFound
	default:
		return "", secret.ErrTerminal
	}
}
