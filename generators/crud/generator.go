package crud

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fluxynet/gocipe/generators"
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

//Description the description of the command when used on cli
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

//Generate produce the generated code according to options
func Generate(generator Generator) string {
	if len(generator.Structure) == 0 || len(generator.Filename) == 0 {
		fmt.Fprintln(os.Stderr, "Missing arguments: file, struct")
		generator.FlagSet.PrintDefaults()
		os.Exit(1)
	}

	structInfo, err := generators.NewStructureInfo(generator.Filename, generator.Structure)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(structInfo)

	generated := make([]string, 1)
	generated[0] = "package " + structInfo.Package + "\n"

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

	targetFilename := filepath.Dir(generator.Filename) + "/" + strings.ToLower(structInfo.Name) + "_crud.go"
	output := strings.Join(generated, "")
	ioutil.WriteFile(targetFilename, []byte(output), 0644)

	return output
}
