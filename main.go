package main

import (
	"fmt"
	"os"

	"github.com/fluxynet/gocipe/generators/crud"
	"github.com/fluxynet/gocipe/generators/http"
	"github.com/fluxynet/gocipe/generators/rest"
)

func main() {
	if len(os.Args) < 2 {
		listCommands()
	}

	switch os.Args[1] {
	case crud.Command:
		crudGenerator := crud.NewGenerator()
		crudGenerator.FlagSet.Parse(os.Args[2:])
		crud.Generate(*crudGenerator)
	case rest.Command:
		restGenerator := rest.NewGenerator()
		restGenerator.FlagSet.Parse(os.Args[2:])
		rest.Generate(*restGenerator)
	case http.Command:
		httpGenerator := http.NewGenerator()
		httpGenerator.FlagSet.Parse(os.Args[2:])
		http.Generate(*httpGenerator)

	default:
		listCommands()
	}
}

func listCommands() {
	commands := map[string]string{
		http.Command: http.Description,
		crud.Command: crud.Description,
		rest.Command: rest.Description,
	}

	for name, command := range commands {
		fmt.Println("Please enter a valid command. Available Commands:")
		fmt.Printf("\t%10s\t%s\n", name, command)
	}
	os.Exit(1)
}
