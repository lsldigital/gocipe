package main

import (
	"fmt"
	"os"

	"github.com/fluxynet/gocipe/generators/crud"
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
	default:
		listCommands()
	}
}

func listCommands() {
	commands := map[string]string{
		crud.Command: crud.Description,
	}

	for name, command := range commands {
		fmt.Println("Please enter a valid command. Available Commands:")
		fmt.Printf("\t%10s\t%s\n", name, command)
	}
	os.Exit(1)
}
