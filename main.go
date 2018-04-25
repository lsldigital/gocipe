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

	go func(recipePath string) {
		var outlog []string
		outcome := true

		for generated := range done {
			if generated.Error == nil {
				l, ok := saveGenerated(generated)
				outlog = append(outlog, l)
				outcome = outcome && ok
			} else {
				fmt.Println(generated.Generator, " Error: ", generated.Error)
				outcome = false
			}
			work.Waitgroup.Done()
		}

		err = ioutil.WriteFile(recipePath+".log", []byte(strings.Join(outlog, "\n")), os.FileMode(0755))
		if err != nil {
			fmt.Printf("failed to write file log file %s.log: %s", recipePath, err)
		}

		if outcome {
			fmt.Printf("Generated %d files without error.\n", len(outlog))
		} else {
			fmt.Printf("Some errors occurred during recipe generation. See %s.log for details.\n", recipePath)
		}
	}(recipePath)

	work.Waitgroup.Wait()
	close(done)

}

func saveGenerated(generated util.GeneratedCode) (string, bool) {
	filename, err := util.GetAbsPath(generated.Filename)
	if err != nil {
		return fmt.Sprintf("[WriteError] cannot resolve path [%s] %s: %s", generated.Generator, generated.Filename, err), false
	}

	if util.FileExists(filename) && generated.NoOverwrite {
		return fmt.Sprintf("[WriteError] skipping existing file [%s] %s", generated.Generator, generated.Filename), false
	}

	var mode os.FileMode = 0755
	if err = os.MkdirAll(path.Dir(filename), mode); err != nil {
		return fmt.Sprintf("[WriteError] directory creation failed [%s] %s: %s", generated.Generator, generated.Filename, err), false
	}

	err = ioutil.WriteFile(filename, []byte(generated.Code), mode)
	if err != nil {
		return fmt.Sprintf("[WriteError] failed to write file [%s] %s: %s", generated.Generator, generated.Filename, err), false
	}

	return fmt.Sprintf("[Wrote] %s", filename), true
}
