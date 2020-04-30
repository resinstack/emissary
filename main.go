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
	fmt.Printf("%v\n", t)
}
