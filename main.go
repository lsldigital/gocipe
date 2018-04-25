package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
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
	files := make(map[string]string)

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

	recipeContent, err := ioutil.ReadFile(recipePath)
	if err != nil {
		log.Fatalln("could not read file: ", err)
	}

	err = json.Unmarshal(recipeContent, &recipe)
	if err != nil {
		log.Fatalln("recipe decoding failed: ", err)
	}

	util.SetTemplates(rice.MustFindBox("templates"))

	work.Waitgroup.Add(1)

	// go generators.GenerateBootstrap(work, recipe.Bootstrap)
	// go generators.GenerateHTTP(work, recipe.HTTP)
	// go generators.GenerateCrud(work, recipe.Crud, recipe.Entities)
	// go generators.GenerateREST(work, recipe.Rest, recipe.Entities)
	// go generators.GenerateSchema(work, recipe.Schema, recipe.Entities)
	go generators.GenerateVuetify(work, recipe.Rest, recipe.Vuetify, recipe.Entities)

	go func() {
		for generated := range done {
			if generated.Error == nil {
				files[generated.Filename] = generated.Code
			} else {
				fmt.Println(generated.Generator, " Error: ", generated.Error)
			}
			work.Waitgroup.Done()
		}
	}()

	work.Waitgroup.Wait()
	close(done)

	// for filename, code := range files {
	// 	fmt.Println(filename, code)
	// }
}
