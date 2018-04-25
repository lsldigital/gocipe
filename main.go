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

	var (
		recipe *util.Recipe
		files  []string
	)
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

	go func() {
		for generated := range done {
			if generated.Error == nil {
				files = append(files, generated.Filename)
				saveGenerated(generated)
			} else {
				fmt.Println(generated.Generator, " Error: ", generated.Error)
			}
			work.Waitgroup.Done()
		}
	}()

	work.Waitgroup.Wait()
	close(done)

	err = ioutil.WriteFile(recipePath+".log", []byte(strings.Join(files, "\n")), os.FileMode(0644))
	if err != nil {
		log.Fatalf("failed to write file log file %s.log: %s", recipePath, err)
	}
}

func saveGenerated(generated util.GeneratedCode) error {
	filename, err := util.GetAbsPath(generated.Filename)
	if err != nil {
		return fmt.Errorf("cannot resolve path [%s] %s: %s", generated.Generator, generated.Filename, err)
	}

	if util.FileExists(filename) && generated.NoOverwrite {
		return fmt.Errorf("skipping existing file [%s] %s", generated.Generator, generated.Filename)
	}

	var mode os.FileMode = 0644
	if err = os.MkdirAll(path.Dir(filename), mode); err != nil {
		return fmt.Errorf("directory creation failed [%s] %s: %s", generated.Generator, generated.Filename, err)
	}

	err = ioutil.WriteFile(filename, []byte(generated.Code), mode)
	if err != nil {
		return fmt.Errorf("failed to write file [%s] %s: %s", generated.Generator, generated.Filename, err)
	}

	return nil
}
