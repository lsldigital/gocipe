package crud

import (
	"flag"
	"fmt"
	"os"
	"log"
	"projects/gocipe/generators"
	"strings"
)

//Command the name of the command to start this generator
const Command = "crud"

//Generator represent arguments accepted by this generator
type Generator struct {
	FlagSet        *flag.FlagSet
	Filename       string
	Structure      string
	GenerateDelete bool
	GenerateGet    bool
	GenerateSave   bool
	GenerateList   bool
}

//Command the description of the command when used on cli
const Description = "Generate CRUD functions and methods for an entity"

//NewGenerator returns a new Generator
func NewGenerator() *Generator {
	arguments := new(Generator)
	arguments.FlagSet = flag.NewFlagSet("crud", flag.ExitOnError)

	arguments.FlagSet.StringVar(&arguments.Filename, "file", "", "Filename where struct is located")
	arguments.FlagSet.StringVar(&arguments.Structure, "struct", "", "Name of the structure to use")
	arguments.FlagSet.BoolVar(&arguments.GenerateGet, "g", true, "Generate Get")
	arguments.FlagSet.BoolVar(&arguments.GenerateList, "l", true, "Generate List")
	arguments.FlagSet.BoolVar(&arguments.GenerateDelete, "d", true, "Generate Delete")
	arguments.FlagSet.BoolVar(&arguments.GenerateSave, "s", true, "Generate Save")

	return arguments
}

func Generate(generator Generator) string {
	if len(generator.Structure) == 0 || len(generator.Filename) == 0 {
		fmt.Fprintln(os.Stderr, "Missing arguments: file, struct\n")
		generator.FlagSet.PrintDefaults()
		os.Exit(1)
	}

	structInfo, err := generators.ProcessFile(generator.Filename, generator.Structure)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(structInfo)

	generated := make([]string, 0)

	if generator.GenerateGet {
		segment, err := GenerateGet(*structInfo)

		if err != nil {
			log.Fatalf("An error occured during GenerateGet: %s\n", err)
		}

		generated = append(generated, segment)
	}

	if generator.GenerateList {
		segment, err := GenerateList(*structInfo)

		if err != nil {
			log.Fatalf("An error occured during GenerateList: %s\n", err)
		}

		generated = append(generated, segment)
	}

	if generator.GenerateDelete {
		segment, err := GenerateDelete(*structInfo)

		if err != nil {
			log.Fatalf("An error occured during GenerateDelete: %s\n", err)
		}

		generated = append(generated, segment)
	}

	if generator.GenerateSave {
		segment, err := GenerateSave(*structInfo)

		if err != nil {
			log.Fatalf("An error occured during GenerateSave: %s\n", err)
		}

		generated = append(generated, segment)
	}

	fmt.Println(generated)

	return strings.Join(generated, "\n")
}