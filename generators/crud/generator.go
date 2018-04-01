package crud

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/fluxynet/gocipe/generators"
)

func init() {
	generators.AddCommand("crud", "Generate CRUD functions and methods for an entity", factory)
}

type generator struct {
	FlagSet        *flag.FlagSet
	Filename       string
	Structure      string
	GenerateDelete bool
	GenerateGet    bool
	GenerateSave   bool
	GenerateList   bool
}

func factory(args []string) (generators.Command, error) {
	var g generator
	flagset := flag.NewFlagSet("crud", flag.ExitOnError)

	flagset.StringVar(&g.Filename, "file", "", "Filename where struct is located")
	flagset.StringVar(&g.Structure, "struct", "", "Name of the structure to use")
	flagset.BoolVar(&g.GenerateGet, "g", true, "Generate Get")
	flagset.BoolVar(&g.GenerateList, "l", true, "Generate List")
	flagset.BoolVar(&g.GenerateDelete, "d", true, "Generate Delete")
	flagset.BoolVar(&g.GenerateSave, "s", true, "Generate Save")

	flagset.Parse(args)

	if len(g.Structure) == 0 || len(g.Filename) == 0 {
		flagset.PrintDefaults()
		return nil, errors.New("Missing arguments: file, struct")
	}

	return g, nil
}

//Generate produce the generated code according to options
func (g generator) Generate() (string, error) {
	structInfo, err := generators.NewStructureInfo(g.Filename, g.Structure)
	if err != nil {
		return "", err
	}

	var generated []string
	generated = append(generated, "package "+structInfo.Package+"\n")
	generated = append(generated, `
import "database/sql"

var db *sql.DB

// Inject allows injection of services into the package
func Inject(database *sql.DB) {
	db = database
}
`)

	if g.GenerateGet {
		segment, err := GenerateGet(*structInfo)

		if err != nil {
			log.Fatalf("An error occured during GenerateGet: %s\n", err)
		}

		generated = append(generated, segment)
	}

	if g.GenerateList {
		segment, err := GenerateList(*structInfo)

		if err != nil {
			log.Fatalf("An error occured during GenerateList: %s\n", err)
		}

		generated = append(generated, segment)
	}

	if g.GenerateDelete {
		segment, err := GenerateDelete(*structInfo)

		if err != nil {
			log.Fatalf("An error occured during GenerateDelete: %s\n", err)
		}

		generated = append(generated, segment)
	}

	if g.GenerateSave {
		var (
			segment string
			err     error
		)
		segment, err = GenerateInsert(*structInfo)

		if err != nil {
			log.Fatalf("An error occured during GenerateInsert: %s\n", err)
		}

		generated = append(generated, segment)

		segment, err = GenerateUpdate(*structInfo)
		if err != nil {
			log.Fatalf("An error occured during GenerateUpdate: %s\n", err)
		}

		generated = append(generated, segment)
		segment, err = GenerateSave(*structInfo)
		if err != nil {
			log.Fatalf("An error occured during GenerateSave: %s\n", err)
		}

		generated = append(generated, segment)
	}

	targetFilename := filepath.Dir(g.Filename) + "/" + strings.ToLower(structInfo.Name) + "_crud.go"
	output := strings.Join(generated, "\n")
	err = ioutil.WriteFile(targetFilename, []byte(output), 0644)

	return output, err
}
