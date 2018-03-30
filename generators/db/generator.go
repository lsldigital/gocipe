package db

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/fluxynet/gocipe/generators"
)

//Command the name of the command to start this generator
const Command = "db"

//Generator represent arguments accepted by this generator
type Generator struct {
	FlagSet   *flag.FlagSet
	Filename  string
	Structure string
	Output    string
}

//Description the description of the command when used on cli
const Description = "Generate database schema for a structure"

//NewGenerator returns a new Generator
func NewGenerator() *Generator {
	arguments := new(Generator)
	arguments.FlagSet = flag.NewFlagSet("crud", flag.ExitOnError)

	arguments.FlagSet.StringVar(&arguments.Filename, "file", "", "Filename where struct is located")
	arguments.FlagSet.StringVar(&arguments.Structure, "struct", "", "Name of the structure to use")
	arguments.FlagSet.StringVar(&arguments.Output, "output", "", "File to output for the schema definition")

	return arguments
}

//Generate produce the generated code according to options
func Generate(generator Generator) string {
	if len(generator.Structure) == 0 || len(generator.Filename) == 0 || len(generator.Output) == 0 {
		fmt.Fprintln(os.Stderr, "Missing arguments: file, struct, output")
		generator.FlagSet.PrintDefaults()
		os.Exit(1)
	}

	structInfo, err := generators.NewStructureInfo(generator.Filename, generator.Structure)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(structInfo)

	generated, err := GenerateDatabase(*structInfo)
	if err != nil {
		log.Fatalln(err)
	}

	targetFilename := filepath.Dir(generator.Output)
	ioutil.WriteFile(targetFilename, []byte(generated), 0644)

	return generated
}
