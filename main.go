package main

import (
	"fmt"

	"github.com/the-maldridge/emissary/pkg/tmpl"
)

func main() {
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
