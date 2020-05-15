package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/resinstack/emissary/pkg/secret"
	"github.com/resinstack/emissary/pkg/tmpl"

	_ "github.com/resinstack/emissary/pkg/secret/awssm"
	_ "github.com/resinstack/emissary/pkg/secret/insecure"
)

func doTemplate(path string, wg *sync.WaitGroup) {
	log.Println("Launching template worker for", path)
	defer wg.Done()

	t, err := tmpl.Parse(path)
	if err != nil {
		log.Printf("Error parsing template at %s: %s", path, err)
		return
	}

	if err := t.Render(); err != nil {
		log.Printf("Error rendering template at %s: %s", path, err)
		return
	}
	log.Printf("Template worker for %s is terminating.", path)
}

func startupDelay() error {
	if delay := os.Getenv("EMISSARY_STARTUP_DELAY"); delay != "" {
		d, err := time.ParseDuration(delay)
		if err != nil {
			// This is really the cleanest solution to
			// just jump out of this case.
			return err
		}
		log.Printf("Initial startup delay of %s", d)
		time.Sleep(d)
	}
	return nil
}

func main() {
	log.Println("Emissary is starting")
	if err := startupDelay(); err != nil {
		log.Printf("Startup delay error: %s", err)
	}
	log.Println("Startup delay phase is complete")

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
	log.Println("All template workers terminated")
	log.Println("Emissary is shutting down")
}
