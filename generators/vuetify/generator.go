package vuetify

import (
	"errors"
	"flag"
	"io/ioutil"
	"os"

	"github.com/fluxynet/gocipe/generators"
	"github.com/jinzhu/inflection"
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
	flagset.StringVar(&g.Output, "output", "", "Folder where to output generated vuetify components")

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
		path   string
		name   string
	)

	structInfo, err := generators.NewStructureInfo(g.Filename, g.Structure)
	if err != nil {
		return "", err
	}

	path, err = generators.GetAbsPath(g.Output)
	if err != nil {
		return "", nil
	}

	name = inflection.Plural(structInfo.Name)
	path = path + string(os.PathSeparator) + name

	if g.GenerateEditor {
		output, err = GenerateEditor(*structInfo)

		if err != nil {
			return "", err
		}

		err = ioutil.WriteFile(path+"Edit.vue", []byte(output), 0644)
		if err != nil {
			return "", err
		}
	}

	if g.GenerateListing {
		output, err = GenerateList(*structInfo)

		if err != nil {
			return "", err
		}

		err = ioutil.WriteFile(path+"List.vue", []byte(output), 0644)

		if err != nil {
			return "", err
		}
	}

	return output, err
}
