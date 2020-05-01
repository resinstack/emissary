package main

import (
	"fmt"

	"github.com/the-maldridge/emissary/pkg/secret"
	"github.com/the-maldridge/emissary/pkg/tmpl"

	_ "github.com/the-maldridge/emissary/pkg/secret/insecure"
)

func main() {
	secret.InitializeProviders()

	t, err := tmpl.Parse("test.tmpl")
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := t.Render(); err != nil {
		fmt.Println(err)
		return
	}
}
