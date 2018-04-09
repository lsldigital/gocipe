package vuetify

import (
	"errors"
	"flag"
	"io/ioutil"

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
	flagset.StringVar(&g.Output, "output", "", "File to output for the vuetify component (suffix will be added automatically)")

	flagset.Parse(args)

	if len(g.Structure) == 0 || len(g.Filename) == 0 {
		flagset.PrintDefaults()
		return nil, errors.New("Missing arguments: file, struct")
	}

	return g, nil
}

func (g generator) Generate() (string, error) {
	var (
		err    error
		output string
	)

	structInfo, err := generators.NewStructureInfo(g.Filename, g.Structure)
	if err != nil {
		return "", err
	}

	output, err = GenerateEditor(*structInfo)

	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(g.Output+"Edit.vue", []byte(output), 0644)
	if err != nil {
		return "", err
	}

	output, err = GenerateList(*structInfo)

	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(g.Output+"List.vue", []byte(output), 0644)

	if err != nil {
		return "", err
	}

	return output, err
}
