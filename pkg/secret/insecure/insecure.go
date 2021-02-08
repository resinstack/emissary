package insecure

import (
	"errors"
	"io/ioutil"
	"log"
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
	urlfile := os.Getenv("EMISSARY_INSECURE_URLFILE")
	if urlfile == "" {
		return nil, errors.New("required variable EMISSARY_INSECURE_URLFILE not set, insecure engine will be unavailable")
	}

	base, err := ioutil.ReadFile(urlfile)
	if err != nil {
		return nil, errors.New(urlfile + " is not readable")
	}
	baseS := strings.TrimSpace(string(base[:]))

	i.base, err = url.Parse(baseS)
	if err != nil {
		return nil, errors.New(baseS + " is not a URL, insecure engine will be unavailable")
	}
	log.Println("Insecure Secrets engine is initialized and available")
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
