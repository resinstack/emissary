package tmpl

import (
	"io/ioutil"
	"text/template"

	"github.com/ericaro/frontmatter"
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

	t.Template, err = t.Template.Parse(t.Content)
	if err != nil {
		return nil, err
	}

	return t, nil
}
