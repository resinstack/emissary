package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/the-maldridge/emissary/pkg/secret"
	"github.com/the-maldridge/emissary/pkg/tmpl"

	_ "github.com/the-maldridge/emissary/pkg/secret/awssm"
	_ "github.com/the-maldridge/emissary/pkg/secret/insecure"
)

func doTemplate(path string, wg *sync.WaitGroup) {
	log.Println("Launching template worker for", path)
	defer wg.Done()

	t, err := tmpl.Parse(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}

	if err := t.Render(); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}
}

func main() {
	secret.InitializeProviders()

	basepath := os.Getenv("EMISSARY_TPL_DIR")
	if basepath == "" {
		basepath, _ = os.Getwd()
	}
	log.Printf("Searching for templates in %s", basepath)

	// The only errors that filepath can throw are related to
	// parsing the glob.  Its hardcoded here and this is known to
	// work.
	tmpls, _ := filepath.Glob(filepath.Join(basepath, "*.tpl"))

	var wg sync.WaitGroup
	for _, t := range tmpls {
		wg.Add(1)
		go doTemplate(t, &wg)
	}
	wg.Wait()
}
