package rest

import (
	"errors"
	"flag"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/fluxynet/gocipe/generators"
)

func init() {
	generators.AddCommand("rest", "Generate REST endpoint functions and methods for an entity", factory)
}

type generator struct {
	FlagSet        *flag.FlagSet
	Filename       string
	Structure      string
	GenerateDelete bool
	GenerateGet    bool
	GenerateCreate bool
	GenerateUpdate bool
	GenerateList   bool
}

func factory(args []string) (generators.Command, error) {
	var g generator
	flagset := flag.NewFlagSet("rest", flag.ExitOnError)

	flagset.StringVar(&g.Filename, "file", "", "Filename where struct is located")
	flagset.StringVar(&g.Structure, "struct", "", "Name of the structure to use")
	flagset.BoolVar(&g.GenerateGet, "g", true, "Generate Get")
	flagset.BoolVar(&g.GenerateList, "l", true, "Generate List")
	flagset.BoolVar(&g.GenerateDelete, "d", true, "Generate Delete")
	flagset.BoolVar(&g.GenerateCreate, "s", true, "Generate Create")
	flagset.BoolVar(&g.GenerateUpdate, "u", true, "Generate Update")

	flagset.Parse(args)

	if len(g.Structure) == 0 || len(g.Filename) == 0 {
		flagset.PrintDefaults()
		return nil, errors.New("Missing arguments: file, struct")
	}

	return g, nil
}

func (g generator) Generate() (string, error) {
	structInfo, err := generators.NewStructureInfo(g.Filename, g.Structure)
	if err != nil {
		return "", err
	}

	var generated []string
	generated = append(generated, "package "+structInfo.Package+"\n")

	if segment, err := GenerateStructures(*structInfo); err == nil {
		generated = append(generated, segment)
	} else {
		return "", err
	}

	if segment, err := GenerateRoutes(*structInfo, g); err == nil {
		generated = append(generated, segment)
	} else {
		return "", err
	}

	if g.GenerateGet {
		segment, err := GenerateGet(*structInfo)

		if err != nil {
			return "", err
		}

		generated = append(generated, segment)
	}

	if g.GenerateList {
		segment, err := GenerateList(*structInfo)

		if err != nil {
			return "", err
		}

		generated = append(generated, segment)
	}

	if g.GenerateDelete {
		segment, err := GenerateDelete(*structInfo)

		if err != nil {
			return "", err
		}

		generated = append(generated, segment)
	}

	if g.GenerateCreate {
		segment, err := GenerateCreate(*structInfo)

		if err != nil {
			return "", err
		}

		generated = append(generated, segment)
	}

	if g.GenerateUpdate {
		segment, err := GenerateUpdate(*structInfo)

		if err != nil {
			return "", err
		}

		generated = append(generated, segment)
	}

	targetFilename := filepath.Dir(g.Filename) + "/" + strings.ToLower(structInfo.Name) + "_rest.go"
	output := strings.Join(generated, "\n")
	err = ioutil.WriteFile(targetFilename, []byte(output), 0644)

	return output, err
}
