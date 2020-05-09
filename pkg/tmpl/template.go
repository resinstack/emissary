package tmpl

import (
	"io/ioutil"
	"log"
	"os/exec"
	"text/template"

	"github.com/ericaro/frontmatter"
	"github.com/google/renameio"
	"github.com/google/shlex"

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
	f, err := renameio.TempFile("", t.Dest)
	if err != nil {
		return err
	}
	defer f.Cleanup()

	if err := f.Chmod(t.Mode); err != nil {
		return err
	}

	if err := t.Template.Execute(f, nil); err != nil {
		return err
	}
	if err := f.CloseAtomicallyReplace(); err != nil {
		return err
	}

	if t.OnRender != "" {
		cmd, err := shlex.Split(t.OnRender)
		if err != nil {
			return err
		}
		// Run whatever command was supposed to happen after
		// the template was rendered out.
		out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
		log.Printf("Output from OnRender for %s: %s", t.Dest, string(out))
		if err != nil {
			log.Printf("Error running command: %s", err)
			return err
		}
	}

	return nil
}
