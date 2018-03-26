package rest

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
const Command = "rest"

//Generator represent arguments accepted by this generator
type Generator struct {
	FlagSet        *flag.FlagSet
	Filename       string
	Structure      string
	GenerateDelete bool
	GenerateGet    bool
	GenerateCreate bool
	GenerateUpdate bool
	GenerateList   bool
}

//Description the description of the command when used on cli
const Description = "Generate REST endpoint functions and methods for an entity"

//NewGenerator returns a new Generator
func NewGenerator() *Generator {
	arguments := new(Generator)
	arguments.FlagSet = flag.NewFlagSet("rest", flag.ExitOnError)

	arguments.FlagSet.StringVar(&arguments.Filename, "file", "", "Filename where struct is located")
	arguments.FlagSet.StringVar(&arguments.Structure, "struct", "", "Name of the structure to use")
	arguments.FlagSet.BoolVar(&arguments.GenerateGet, "g", true, "Generate Get")
	arguments.FlagSet.BoolVar(&arguments.GenerateList, "l", true, "Generate List")
	arguments.FlagSet.BoolVar(&arguments.GenerateDelete, "d", true, "Generate Delete")
	arguments.FlagSet.BoolVar(&arguments.GenerateCreate, "s", true, "Generate Create")
	arguments.FlagSet.BoolVar(&arguments.GenerateUpdate, "u", true, "Generate Update")

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

	var generated []string
	generated = append(generated, "package "+structInfo.Package+"\n")

	if segment, err := GenerateStructures(*structInfo); err == nil {
		generated = append(generated, segment)
	} else {
		log.Fatalf("An error occured during GenerateStructures: %s\n", err)
	}

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

	if generator.GenerateCreate {
		segment, err := GenerateCreate(*structInfo)

		if err != nil {
			log.Fatalf("An error occured during GenerateCreate: %s\n", err)
		}

		generated = append(generated, segment)
	}

	if generator.GenerateUpdate {
		segment, err := GenerateUpdate(*structInfo)

		if err != nil {
			log.Fatalf("An error occured during GenerateUpdate: %s\n", err)
		}

		generated = append(generated, segment)
	}

	targetFilename := filepath.Dir(generator.Filename) + "/" + strings.ToLower(structInfo.Name) + "_rest.go"
	output := strings.Join(generated, "\n")
	ioutil.WriteFile(targetFilename, []byte(output), 0644)

	return output
}
