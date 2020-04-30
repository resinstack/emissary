package tmpl

import (
	"text/template"
)

// Tmpl contains a parsed template including its metadata.  When
// rendered, it is ready to write out to a file specified by Dest and
// with mode Mode.  Once the file is successfully written, the command
// specified by OnRender will be executed.
type Tmpl struct {
	Dest     string
	OnRender string
	Mode     int

	Content string `fm:"content"`

	Template *template.Template
}
