package vuetify

import (
	"errors"
	"flag"
	"io/ioutil"
	"strings"

	"github.com/fluxynet/gocipe/generators"
)

type generator struct {
	Filename        string
	Structure       string
	Output          string
	GenerateEditor  bool
	GenerateListing bool
}

func init() {
	generators.AddCommand("vuetify", "Generate vuetify components - listing and editor views for entity", factory)
}

func factory(args []string) (generators.Command, error) {
	var g generator

	flagset := flag.NewFlagSet("rest", flag.ExitOnError)

	flagset.StringVar(&g.Filename, "file", "", "Filename where struct is located")
	flagset.StringVar(&g.Structure, "struct", "", "Name of the structure to use")
	flagset.BoolVar(&g.GenerateEditor, "editor", true, "Generate Editor component")
	flagset.BoolVar(&g.GenerateListing, "listing", true, "Generate Listing component")
	flagset.StringVar(&g.Output, "output", "", "File to output for the schema definition")

	flagset.Parse(args)

	if len(g.Structure) == 0 || len(g.Filename) == 0 {
		flagset.PrintDefaults()
		return nil, errors.New("Missing arguments: file, struct")
	}

	return g, nil
}

func (g generator) Generate() (string, error) {
	var (
		generated []string
		err       error
		segment   string
	)

	structInfo, err := generators.NewStructureInfo(g.Filename, g.Structure)
	if err != nil {
		return "", err
	}

	segment, err = GenerateEditor(*structInfo)
	generated = append(generated, segment)

	if err != nil {
		return "", err
	}

	targetFilename := g.Output
	output := strings.Join(generated, "\n")
	err = ioutil.WriteFile(targetFilename, []byte(output), 0644)

	return output, err
}
