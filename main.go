package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"sync"

	rice "github.com/GeertJohan/go.rice"
	"github.com/fluxynet/gocipe/generators"
	"github.com/fluxynet/gocipe/generators/crud"
	"github.com/fluxynet/gocipe/util"
)

//go:generate rice embed-go

var (
	_recipeHash string
	_recipePath string
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	toolset := initToolset()
	noSkip := flag.Bool("noskip", false, "Do not skip overwriting existing files")
	flag.Parse()

	work := util.GenerationWork{
		Waitgroup: new(sync.WaitGroup),
		Done:      make(chan util.GeneratedCode),
	}

	recipe, err := loadRecipe()

	if err != nil {
		log.Fatalln(err)
	}

	util.SetTemplates(rice.MustFindBox("templates"))

	work.Waitgroup.Add(5)

	entities := make(map[string]util.Entity)
	for _, entity := range recipe.Entities {
		entities[entity.Name] = entity
	}

	go generators.GenerateBootstrap(work, recipe.Bootstrap, recipe.HTTP)
	go generators.GenerateHTTP(work, recipe.HTTP)
	go crud.Generate(work, recipe.Crud, entities)
	// go generators.GenerateREST(work, recipe.Rest, recipe.Entities)
	go generators.GenerateSchema(work, recipe.Schema, recipe.Entities)
	go generators.GenerateVuetify(work, recipe.Rest, recipe.Vuetify, recipe.Entities)

	var wg sync.WaitGroup
	wg.Add(1)

	go processOutput(&wg, work, toolset, *noSkip)

	work.Waitgroup.Wait()
	close(work.Done)
	wg.Wait()
}

func loadRecipe() (*util.Recipe, error) {
	var recipe util.Recipe

	if len(os.Args) == 1 {
		return nil, errors.New("no recipe provided. usage: gocipe gocipe.json")
	}

	recipePath, err := util.GetAbsPath(os.Args[len(os.Args)-1])
	if err != nil {
		return nil, err
	}

	if !util.FileExists(recipePath) {
		return nil, fmt.Errorf("file not found: %s", recipePath)
	}

	recipeContent, err := ioutil.ReadFile(recipePath)

	if err != nil {
		return nil, err
	}

	_recipeHash = fmt.Sprintf("%x", sha256.Sum256([]byte(recipeContent)))
	_recipePath = recipePath

	err = json.Unmarshal(recipeContent, &recipe)
	if err != nil {
		return nil, fmt.Errorf("recipe decoding failed: %s", err)
	}

	return &recipe, nil
}

func processOutput(waitgroup *sync.WaitGroup, work util.GenerationWork, toolset util.Toolset, noSkip bool) {

	var (
		outlog, output, gofiles  []string
		written, skipped, failed int
		err                      error
	)

	aggregates := make(map[string][]util.GeneratedCode)
	outlog = append(outlog, "[Recipe Hash] "+_recipeHash)

	for generated := range work.Done {
		if generated.Error == util.ErrorSkip {
			outlog = append(outlog, fmt.Sprintf("[Skipped] Generation skipped [%s]", generated.Generator))
			skipped++
		} else if generated.Error != nil {
			outlog = append(outlog, fmt.Sprintf("[Error] Generation failed [%s]: %s", generated.Generator, generated.Error))
			failed++
		} else if generated.Aggregate {
			a := aggregates[generated.Filename]
			aggregates[generated.Filename] = append(a, generated)
		} else {
			fname, l, err := saveGenerated(generated, noSkip)
			outlog = append(outlog, l)

			if err == nil {
				if strings.HasSuffix(fname, ".go") {
					gofiles = append(gofiles, fname)
				}

				written++
			} else if err == util.ErrorSkip {
				skipped++
			} else {
				failed++
			}
		}
		work.Waitgroup.Done()
	}

	for _, generated := range aggregates {
		fname, l, err := saveAggregate(generated, noSkip)
		outlog = append(outlog, l)

		if err == nil {
			if strings.HasSuffix(fname, ".go") {
				gofiles = append(gofiles, fname)
			}

			written++
		} else if err == util.ErrorSkip {
			skipped++
		} else {
			failed++
		}
	}

	err = ioutil.WriteFile(_recipePath+".log", []byte(strings.Join(outlog, "\n")), os.FileMode(0755))
	if err != nil {
		fmt.Printf("failed to write file log file %s.log: %s", _recipePath, err)
		return
	}

	if skipped > 0 {
		output = append(output, fmt.Sprintf("Skipped %d files.", skipped))
	}

	if written > 0 {
		output = append(output, fmt.Sprintf("Wrote %d files.", written))
	}

	if failed > 0 {
		output = append(output, fmt.Sprintf("%d errors occurred during recipe generation.", failed))
	}

	if len(gofiles) > 0 {
		postProcessGoFiles(toolset, gofiles)
	}

	output = append(output, fmt.Sprintf("See log file %s.log for details.", _recipePath))
	fmt.Println(strings.Join(output, " "))
	waitgroup.Done()
}

// saveGenerated saves a generated file and returns absolute filename, log entry and error
func saveGenerated(generated util.GeneratedCode, noSkip bool) (string, string, error) {
	filename, err := util.GetAbsPath(generated.Filename)
	if err != nil {
		return "", fmt.Sprintf("[WriteError] cannot resolve path [%s] %s: %s", generated.Generator, generated.Filename, err), err
	}

	if !noSkip && util.FileExists(filename) && generated.NoOverwrite {
		return "", fmt.Sprintf("[Skip] skipping existing file [%s] %s", generated.Generator, generated.Filename), util.ErrorSkip
	}

	var mode os.FileMode = 0755
	if err = os.MkdirAll(path.Dir(filename), mode); err != nil {
		return "", fmt.Sprintf("[WriteError] directory creation failed [%s] %s: %s", generated.Generator, generated.Filename, err), err
	}

	var code []byte
	if generated.NoOverwrite {
		code = []byte(generated.Code)
	} else {
		var generatedHeaderFormat string
		if generated.GeneratedHeaderFormat == "" {
			generatedHeaderFormat = "// %s"
		} else {
			generatedHeaderFormat = generated.GeneratedHeaderFormat
		}

		generatedHeaderFormat = fmt.Sprintf(generatedHeaderFormat, `generated by gocipe `+_recipeHash+`; DO NOT EDIT`)

		code = []byte(generatedHeaderFormat + "\n\n" + generated.Code)
	}

	err = ioutil.WriteFile(filename, code, mode)
	if err != nil {
		return "", fmt.Sprintf("[WriteError] failed to write file [%s] %s: %s", generated.Generator, generated.Filename, err), err
	}

	return filename, fmt.Sprintf("[Wrote] %s", filename), nil
}

// saveAggregate saves aggregated files and returns absolute filename, log entry and error
func saveAggregate(aggregate []util.GeneratedCode, noSkip bool) (string, string, error) {
	var generated util.GeneratedCode

	generated.Filename = aggregate[0].Filename
	generated.Generator = aggregate[0].Generator
	generated.GeneratedHeaderFormat = aggregate[0].GeneratedHeaderFormat

	for _, g := range aggregate {
		generated.NoOverwrite = generated.NoOverwrite || g.NoOverwrite
		generated.Code += g.Code + "\n"
	}

	return saveGenerated(generated, noSkip)
}

func initToolset() util.Toolset {
	var (
		err error
		ok  = true
	)

	goimports, err := exec.LookPath("goimports")
	if err != nil {
		fmt.Println(err)
		ok = false
	}

	gofmt, err := exec.LookPath("gofmt")
	if err != nil {
		fmt.Println(err)
		ok = false
	}

	if !ok {
		log.Fatalln("Required tools missing: goimports and gofmt")
	}

	return util.Toolset{GoFmt: gofmt, GoImports: goimports}
}

// postProcessGoFiles executes goimports and gofmt on go files that have been generated
func postProcessGoFiles(toolset util.Toolset, gofiles []string) {
	var wg sync.WaitGroup
	wg.Add(len(gofiles))

	for _, file := range gofiles {
		go func(file string) {
			cmd := exec.Command(toolset.GoImports, "-w", file)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()

			if err != nil {
				fmt.Printf("Error running %s on %s: %s\n", toolset.GoImports, file, err)
				wg.Done()
				return
			}

			cmd = exec.Command(toolset.GoFmt, "-w", file)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()

			if err != nil {
				fmt.Printf("Error running %s on %s: %s\n", toolset.GoFmt, file, err)
			}

			wg.Done()
		}(file)
	}

	wg.Wait()
}
