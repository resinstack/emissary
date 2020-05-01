package tmpl

import (
	"os"
	"text/template"
)

// Tmpl contains a parsed template including its metadata.  When
// rendered, it is ready to write out to a file specified by Dest and
// with mode Mode.  Once the file is successfully written, the command
// specified by OnRender will be executed.
type Tmpl struct {
	Dest     string
	OnRender string
	Mode     os.FileMode

	Content string `fm:"content"`

	Template *template.Template
}

// A SecretProvider is a function that fetches a secret from a remote
// secure storage.  The first argument is the engine to retrieve the
// secret from, and the second argument is the name of the secret to
// return.
type SecretProvider func(string, string) (string, error)
