package generators

import (
	"fmt"
	"log"
	"os"
)

//Command interface provides blueprint for an object which can be used to generate code from a list of arguments (cli)
type Command interface {
	Generate() (string, error)
}

type cmd struct {
	Description string
	Factory     func(args []string) (Command, error)
}

var commands = make(map[string]cmd)

//AddCommand registers a generator command in the bootstrapper
func AddCommand(name string, description string, factory func(args []string) (Command, error)) {
	commands[name] = cmd{description, factory}
}

//Execute will execute generator from command line arguments
func Execute() {
	if len(os.Args) < 2 {
		listCommands()
	}

	name := os.Args[1]
	if command, ok := commands[name]; ok {
		generator, err := command.Factory(os.Args[2:])

		if err != nil {
			log.Fatalln(err)
		}

		output, err := generator.Generate()
		if output != "" {
			fmt.Println(output)
		}

		if err != nil {
			log.Fatalln(err)
		}
	} else {
		listCommands()
	}
}

func listCommands() {
	fmt.Println("Please enter a valid command. Available Commands:")
	for name, command := range commands {
		fmt.Printf("\t%10s\t%s\n", name, command.Description)
	}
	os.Exit(1)
}
