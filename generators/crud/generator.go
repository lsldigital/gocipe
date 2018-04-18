package crud

import (
	"errors"
	"flag"
	"fmt"
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
	FlagSet   *flag.FlagSet
	Filename  string
	Structure string

	// If true, will print the output (model, functions/methods and hooks). False by default
	Verbose bool

	GenerateDelete     bool
	GeneratePreDelete  bool
	GeneratePostDelete bool

	GenerateGet     bool
	GeneratePreGet  bool
	GeneratePostGet bool

	GenerateSave     bool
	GeneratePreSave  bool
	GeneratePostSave bool

	GenerateList     bool
	GeneratePreList  bool
	GeneratePostList bool
}

func factory(args []string) (generators.Command, error) {
	var g generator
	flagset := flag.NewFlagSet("crud", flag.ExitOnError)

	flagset.StringVar(&g.Filename, "file", "", "Filename where struct is located")
	flagset.StringVar(&g.Structure, "struct", "", "Name of the structure to use")

	flagset.BoolVar(&g.Verbose, "v", false, "Prints the generated models, functions/methods and hooks")

	flagset.BoolVar(&g.GenerateDelete, "d", true, "Generate Delete")
	flagset.BoolVar(&g.GeneratePreDelete, "hd", false, "Generate Delete pre-execution hook")
	flagset.BoolVar(&g.GeneratePostDelete, "dh", false, "Generate Delete post-execution hook")

	flagset.BoolVar(&g.GenerateGet, "g", true, "Generate Get")
	flagset.BoolVar(&g.GeneratePreGet, "hg", false, "Generate Get pre-execution hook")
	flagset.BoolVar(&g.GeneratePostGet, "gh", false, "Generate Get post-execution hook")

	flagset.BoolVar(&g.GenerateSave, "s", true, "Generate Save")
	flagset.BoolVar(&g.GeneratePreSave, "hs", false, "Generate Save pre-execution hook")
	flagset.BoolVar(&g.GeneratePostSave, "sh", false, "Generate Save post-execution hook")

	flagset.BoolVar(&g.GenerateList, "l", true, "Generate List")
	flagset.BoolVar(&g.GeneratePreList, "hl", false, "Generate List pre-execution hook")
	flagset.BoolVar(&g.GeneratePostList, "lh", false, "Generate List post-execution hook")

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

	var (
		generated []string
		hooks     []string
	)

	generated = append(generated, "package "+structInfo.Package+"\n")

	{
		segment, err := GenerateModel(*structInfo)
		if err != nil {
			log.Fatalf("An error occured during GenerateModel: %s\n", err)
		}
		generated = append(generated, segment)
	}

	if g.GenerateGet {
		segment, err := GenerateGet(*structInfo, g.GeneratePreGet, g.GeneratePostGet)
		if err != nil {
			log.Fatalf("An error occured during GenerateGet: %s\n", err)
		}
		generated = append(generated, segment)

		if g.GeneratePreGet || g.GeneratePostGet {
			segment, err := GenerateGetHook(*structInfo, g.GeneratePreGet, g.GeneratePostGet)
			if err != nil {
				log.Fatalf("An error occured during GenerateGetHook: %s\n", err)
			}
			hooks = append(hooks, segment)
		}
	}

	if g.GenerateList {
		segment, err := GenerateList(*structInfo, g.GeneratePreList, g.GeneratePostList)
		if err != nil {
			log.Fatalf("An error occured during GenerateList: %s\n", err)
		}
		generated = append(generated, segment)

		if g.GeneratePreList || g.GeneratePostList {
			segment, err := GenerateListHook(*structInfo, g.GeneratePreList, g.GeneratePostList)
			if err != nil {
				log.Fatalf("An error occured during GenerateListHook: %s\n", err)
			}
			hooks = append(hooks, segment)
		}
	}

	if g.GenerateDelete {
		segment, err := GenerateDelete(*structInfo, g.GeneratePreDelete, g.GeneratePostDelete)
		if err != nil {
			log.Fatalf("An error occured during GenerateDelete: %s\n", err)
		}
		generated = append(generated, segment)

		if g.GeneratePreDelete || g.GeneratePostDelete {
			segment, err := GenerateDeleteHook(g.GeneratePreDelete, g.GeneratePostDelete)
			if err != nil {
				log.Fatalf("An error occured during GenerateDeleteHook: %s\n", err)
			}
			hooks = append(hooks, segment)
		}
	}

	if g.GenerateSave {
		var (
			segment string
			err     error
		)
		segment, err = GenerateInsert(*structInfo, g.GeneratePreSave, g.GeneratePostSave)

		if err != nil {
			log.Fatalf("An error occured during GenerateInsert: %s\n", err)
		}
		generated = append(generated, segment)

		segment, err = GenerateUpdate(*structInfo, g.GeneratePreSave, g.GeneratePostSave)
		if err != nil {
			log.Fatalf("An error occured during GenerateUpdate: %s\n", err)
		}
		generated = append(generated, segment)

		segment, err = GenerateSave(*structInfo)
		if err != nil {
			log.Fatalf("An error occured during GenerateSave: %s\n", err)
		}
		generated = append(generated, segment)

		if g.GeneratePreSave || g.GeneratePostSave {
			segment, err := GenerateSaveHook(*structInfo, g.GeneratePreSave, g.GeneratePostSave)
			if err != nil {
				log.Fatalf("An error occured during GenerateSaveHook: %s\n", err)
			}
			hooks = append(hooks, segment)
		}
	}

	targetFilename, err := generators.GetAbsPath(filepath.Dir(g.Filename) + "/" + strings.ToLower(structInfo.Name) + "_crud.go")
	if err != nil {
		log.Fatalf("Failed to get absolute path (crud): %s\n", err)
	}
	output := strings.Join(generated, "\n")
	err = ioutil.WriteFile(targetFilename, []byte(output), 0644)
	if err != nil {
		log.Fatalf("Failed to write output to %s: %s", targetFilename, err)
	}

	var hookOutput string
	targetFilename, err = generators.GetAbsPath(filepath.Dir(g.Filename) + "/" + strings.ToLower(structInfo.Name) + "_crud_hooks.go")
	if err != nil {
		log.Fatalf("Failed to get absolute path (crud hooks): %s\n", err)
	}

	if len(hooks) > 0 && !generators.FileExists(targetFilename) {
		hookOutput = "package " + structInfo.Package + " \n" + strings.Join(hooks, "\n")
		err = ioutil.WriteFile(targetFilename, []byte(hookOutput), 0644)
		if err != nil {
			log.Fatalf("Failed to write output to %s: %s", targetFilename, err)
		}
	}

	targetFilename = filepath.Dir(g.Filename) + "/../filters.go"
	err = ioutil.WriteFile(targetFilename, []byte(generateModels()), 0644)
	if err != nil {
		log.Fatalf("Failed to write output to %s: %s", targetFilename, err)
	}

	// if err, out := exec.Command("goimports -w " + targetFilename).Output(); err != nil {
	// 	log.Fatalf("An error occurred during goimports: %s\nOutput:\n%s", err, out)
	// } else {
	// 	fmt.Println(out)
	// }

	if g.Verbose {
		fmt.Println(output+"\n//--- HOOKS ---\n\n"+hookOutput, err)
		fmt.Println(structInfo)
	}

	return "", err
}
