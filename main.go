package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"

	rice "github.com/GeertJohan/go.rice"
	"github.com/fluxynet/gocipe/generators"
	"github.com/fluxynet/gocipe/util"
)

//go:generate rice embed-go

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var recipe *util.Recipe

	done := make(chan util.GeneratedCode)
	work := util.GenerationWork{
		Waitgroup: new(sync.WaitGroup),
		Done:      done,
	}

	if len(os.Args) == 1 {
		log.Fatalln("Usage: gocipe gocipe.json")
	}

	recipePath, err := util.GetAbsPath(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	if !util.FileExists(recipePath) {
		log.Fatalf("file not found: %s", recipePath)
	}

	recipeContent, err := ioutil.ReadFile(recipePath)
	if err != nil {
		log.Fatalln("could not read file: ", err)
	}

	err = json.Unmarshal(recipeContent, &recipe)
	if err != nil {
		log.Fatalln("recipe decoding failed: ", err)
	}

	util.SetTemplates(rice.MustFindBox("templates"))

	work.Waitgroup.Add(6)

	go generators.GenerateBootstrap(work, recipe.Bootstrap)
	go generators.GenerateHTTP(work, recipe.HTTP)
	go generators.GenerateCrud(work, recipe.Crud, recipe.Entities)
	go generators.GenerateREST(work, recipe.Rest, recipe.Entities)
	go generators.GenerateSchema(work, recipe.Schema, recipe.Entities)
	go generators.GenerateVuetify(work, recipe.Rest, recipe.Vuetify, recipe.Entities)

	var wg sync.WaitGroup
	wg.Add(1)
	go func(recipePath string) {
		var (
			outlog, output           []string
			written, skipped, failed int
		)

		aggregates := make(map[string][]util.GeneratedCode)

		for generated := range done {
			if generated.Error != nil {
				outlog = append(outlog, fmt.Sprintf("[Error] Generation failed [%s]: %s", generated.Generator, generated.Error))
				failed++
			} else if generated.Aggregate {
				a := aggregates[generated.Filename]
				aggregates[generated.Filename] = append(a, generated)
			} else {
				l, err := saveGenerated(generated)
				outlog = append(outlog, l)

				if err == nil {
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
			l, err := saveAggregate(generated)
			outlog = append(outlog, l)

			if err == nil {
				written++
			} else if err == util.ErrorSkip {
				skipped++
			} else {
				failed++
			}
		}

		err = ioutil.WriteFile(recipePath+".log", []byte(strings.Join(outlog, "\n")), os.FileMode(0755))
		if err != nil {
			fmt.Printf("failed to write file log file %s.log: %s", recipePath, err)
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

		output = append(output, fmt.Sprintf("See log file %s.log for details.", recipePath))
		fmt.Println(strings.Join(output, " "))
		wg.Done()
	}(recipePath)

	work.Waitgroup.Wait()
	close(done)
	wg.Wait()
}

func saveGenerated(generated util.GeneratedCode) (string, error) {
	filename, err := util.GetAbsPath(generated.Filename)
	if err != nil {
		return fmt.Sprintf("[WriteError] cannot resolve path [%s] %s: %s", generated.Generator, generated.Filename, err), err
	}

	if util.FileExists(filename) && generated.NoOverwrite {
		return fmt.Sprintf("[Skip] skipping existing file [%s] %s", generated.Generator, generated.Filename), util.ErrorSkip
	}

	var mode os.FileMode = 0755
	if err = os.MkdirAll(path.Dir(filename), mode); err != nil {
		return fmt.Sprintf("[WriteError] directory creation failed [%s] %s: %s", generated.Generator, generated.Filename, err), err
	}

	err = ioutil.WriteFile(filename, []byte(generated.Code), mode)
	if err != nil {
		return fmt.Sprintf("[WriteError] failed to write file [%s] %s: %s", generated.Generator, generated.Filename, err), err
	}

	return fmt.Sprintf("[Wrote] %s", filename), nil
}

func saveAggregate(aggregate []util.GeneratedCode) (string, error) {
	var generated util.GeneratedCode

	generated.Filename = aggregate[0].Filename
	generated.Generator = aggregate[0].Generator

	for _, g := range aggregate {
		generated.NoOverwrite = generated.NoOverwrite || g.NoOverwrite
		generated.Code += g.Code + "\n"
	}

	return saveGenerated(generated)
}
