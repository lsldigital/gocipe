package db

import (
	"errors"
	"flag"
	"io/ioutil"

	"github.com/fluxynet/gocipe/generators"
)

func init() {
	generators.AddCommand("db", "Generate database schema for a structure", factory)
}

type generator struct {
	FlagSet   *flag.FlagSet
	Filename  string
	Structure string
	Output    string
}

func factory(args []string) (generators.Command, error) {
	var g generator
	flagset := flag.NewFlagSet("db", flag.ExitOnError)

	flagset.StringVar(&g.Filename, "file", "", "Filename where struct is located")
	flagset.StringVar(&g.Structure, "struct", "", "Name of the structure to use")
	flagset.StringVar(&g.Output, "output", "", "File to output for the schema definition")

	flagset.Parse(args)

	if len(g.Structure) == 0 || len(g.Filename) == 0 || len(g.Output) == 0 {
		flagset.PrintDefaults()
		return nil, errors.New("Missing arguments: file, struct, output")
	}

	return g, nil
}

//Generate produce the generated code according to options
func (g generator) Generate() (string, error) {

	structInfo, err := generators.NewStructureInfo(g.Filename, g.Structure)
	if err != nil {
		return "", err
	}

	generated, err := GenerateDatabase(*structInfo)
	if err != nil {
		return "", err
	}

	targetFilename := g.Output
	err = ioutil.WriteFile(targetFilename, []byte(generated), 0644)

	return generated, err
}
