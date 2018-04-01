package main

import (
	"github.com/fluxynet/gocipe/generators"
	_ "github.com/fluxynet/gocipe/generators/crud"
	_ "github.com/fluxynet/gocipe/generators/db"
	_ "github.com/fluxynet/gocipe/generators/http"
	_ "github.com/fluxynet/gocipe/generators/rest"
)

func main() {
	generators.Execute()
}
