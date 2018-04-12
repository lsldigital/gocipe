package rest

import (
	"errors"
	"flag"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/fluxynet/gocipe/generators"
	"log"
)

func init() {
	generators.AddCommand("rest", "Generate REST endpoint functions and methods for an entity", factory)
}

type generator struct {
	FlagSet   *flag.FlagSet
	Filename  string
	Structure string

	// If true, will print the output (model, functions/methods and hooks). False by default
	Verbose bool

	GenerateDelete             bool
	GenerateDeletePreExecHook  bool
	GenerateDeletePostExecHook bool

	GenerateGet             bool
	GenerateGetPreExecHook  bool
	GenerateGetPostExecHook bool

	GenerateCreate             bool
	GenerateCreatePreExecHook  bool
	GenerateCreatePostExecHook bool

	GenerateUpdate             bool
	GenerateUpdatePreExecHook  bool
	GenerateUpdatePostExecHook bool

	GenerateList             bool
	GenerateListPreExecHook  bool
	GenerateListPostExecHook bool
}

func factory(args []string) (generators.Command, error) {
	var g generator
	flagset := flag.NewFlagSet("rest", flag.ExitOnError)

	flagset.StringVar(&g.Filename, "file", "", "Filename where struct is located")
	flagset.StringVar(&g.Structure, "struct", "", "Name of the structure to use")

	flagset.BoolVar(&g.Verbose, "v", false, "Prints the generated models, functions/methods and hooks")

	flagset.BoolVar(&g.GenerateDelete, "d", true, "Generate Delete")
	flagset.BoolVar(&g.GenerateDeletePreExecHook, "hd", false, "Generate Delete pre-execution hook")
	flagset.BoolVar(&g.GenerateDeletePostExecHook, "dh", false, "Generate Delete post-execution hook")

	flagset.BoolVar(&g.GenerateGet, "g", true, "Generate Get")
	flagset.BoolVar(&g.GenerateGetPreExecHook, "hg", false, "Generate Get pre-execution hook")
	flagset.BoolVar(&g.GenerateGetPostExecHook, "gh", false, "Generate Get post-execution hook")

	flagset.BoolVar(&g.GenerateCreate, "c", true, "Generate Create")
	flagset.BoolVar(&g.GenerateCreatePreExecHook, "hc", false, "Generate Create pre-execution hook")
	flagset.BoolVar(&g.GenerateCreatePostExecHook, "ch", false, "Generate Create post-execution hook")

	flagset.BoolVar(&g.GenerateCreate, "u", true, "Generate Update")
	flagset.BoolVar(&g.GenerateCreatePreExecHook, "hu", false, "Generate Update pre-execution hook")
	flagset.BoolVar(&g.GenerateCreatePostExecHook, "uh", false, "Generate Update post-execution hook")

	flagset.BoolVar(&g.GenerateList, "l", true, "Generate List")
	flagset.BoolVar(&g.GenerateListPreExecHook, "hl", false, "Generate List pre-execution hook")
	flagset.BoolVar(&g.GenerateListPostExecHook, "lh", false, "Generate List post-execution hook")

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

	var (
		generated []string
		hooks     []string
	)

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
		segment, err := GenerateGet(*structInfo, g.GenerateGetPreExecHook, g.GenerateGetPostExecHook)
		if err != nil {
			return "", err
		}
		generated = append(generated, segment)

		if g.GenerateGetPreExecHook || g.GenerateGetPostExecHook {
			segment, err := GenerateGetHook(g.GenerateGetPreExecHook, g.GenerateGetPostExecHook)
			if err != nil {
				log.Fatalf("An error occured during GenerateGetHook: %s\n", err)
			}
			hooks = append(hooks, segment)
		}
	}

	if g.GenerateList {
		segment, err := GenerateList(*structInfo, g.GenerateListPreExecHook, g.GenerateListPostExecHook)
		if err != nil {
			return "", err
		}
		generated = append(generated, segment)

		if g.GenerateListPreExecHook || g.GenerateListPostExecHook {
			segment, err := GenerateListHook(g.GenerateListPreExecHook, g.GenerateListPostExecHook)
			if err != nil {
				log.Fatalf("An error occured during GenerateListHook: %s\n", err)
			}
			hooks = append(hooks, segment)
		}
	}

	if g.GenerateDelete {
		segment, err := GenerateDelete(*structInfo, g.GenerateDeletePreExecHook, g.GenerateDeletePostExecHook)
		if err != nil {
			return "", err
		}
		generated = append(generated, segment)

		if g.GenerateDeletePreExecHook || g.GenerateDeletePostExecHook {
			segment, err := GenerateDeleteHook(g.GenerateDeletePreExecHook, g.GenerateDeletePostExecHook)
			if err != nil {
				log.Fatalf("An error occured during GenerateDeleteHook: %s\n", err)
			}
			hooks = append(hooks, segment)
		}
	}

	if g.GenerateCreate {
		segment, err := GenerateCreate(*structInfo, g.GenerateCreatePreExecHook, g.GenerateCreatePostExecHook)
		if err != nil {
			return "", err
		}
		generated = append(generated, segment)

		if g.GenerateCreatePreExecHook || g.GenerateCreatePostExecHook {
			segment, err := GenerateCreateHook(g.GenerateCreatePreExecHook, g.GenerateCreatePostExecHook)
			if err != nil {
				log.Fatalf("An error occured during GenerateCreateHook: %s\n", err)
			}
			hooks = append(hooks, segment)
		}
	}

	if g.GenerateUpdate {
		segment, err := GenerateUpdate(*structInfo, g.GenerateUpdatePreExecHook, g.GenerateUpdatePostExecHook)
		if err != nil {
			return "", err
		}
		generated = append(generated, segment)

		if g.GenerateUpdatePreExecHook || g.GenerateUpdatePostExecHook {
			segment, err := GenerateUpdateHook(g.GenerateUpdatePreExecHook, g.GenerateUpdatePostExecHook)
			if err != nil {
				log.Fatalf("An error occured during GenerateUpdateHook: %s\n", err)
			}
			hooks = append(hooks, segment)
		}
	}

	targetFilename := filepath.Dir(g.Filename) + "/" + strings.ToLower(structInfo.Name) + "_rest.go"
	output := strings.Join(generated, "\n")
	err = ioutil.WriteFile(targetFilename, []byte(output), 0644)

	var hookOutput string
	targetFilename, err = generators.GetAbsPath(filepath.Dir(g.Filename) + "/" + strings.ToLower(structInfo.Name) + "_rest_hooks.go")
	if err != nil {
		log.Fatalf("Failed to get absolute path (rest hooks): %s\n", err)
	}
	if len(hooks) > 0 && !generators.FileExists(targetFilename) {
		hookOutput = "package " + structInfo.Package + " \n" + strings.Join(hooks, "\n")
		err = ioutil.WriteFile(targetFilename, []byte(hookOutput), 0644)
		if err != nil {
			log.Fatalf("Failed to write output to %s: %s", targetFilename, err)
		}
	}

	if g.Verbose {
		return output + "\n//--- HOOKS ---\n\n" + hookOutput, err
	}

	return "", err
}
