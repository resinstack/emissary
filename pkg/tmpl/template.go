package tmpl

import (
	"io/ioutil"
	"os"
	"text/template"

	"github.com/ericaro/frontmatter"

	"github.com/the-maldridge/emissary/pkg/secret"
)

var (
	fmap = template.FuncMap{
		"poll": secret.Poll,
	}
)

// Parse attempts to read the file at f and returns a Tmpl pointer
// that contains both the template, and the metadata for where to
// write the template.
func Parse(f string) (*Tmpl, error) {
	t := new(Tmpl)
	t.Template = template.New("")

	fbytes, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	if err := frontmatter.Unmarshal(fbytes, t); err != nil {
		return nil, err
	}

	t.Template.Funcs(fmap)

	t.Template, err = t.Template.Parse(t.Content)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Render spits out the contents of the template and renders it to a
// file on disk.
func (t *Tmpl) Render() error {
	f, err := os.OpenFile(t.Dest, os.O_RDWR|os.O_CREATE, t.Mode)
	if err != nil {
		return err
	}

	if err := t.Template.Execute(f, nil); err != nil {
		return err
	}
	return nil
}
