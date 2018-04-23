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
)

//go:generate rice embed-go

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var recipe *generators.Recipe
	done := make(chan generators.GeneratedCode)

	work := generators.GenerationWork{
		Waitgroup: new(sync.WaitGroup),
		Done:      done,
	}

	if len(os.Args) == 1 {
		log.Fatalln("Usage: gocipe gocipe.json")
	}

	recipePath, err := generators.GetAbsPath(os.Args[1])
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

	generators.SetTemplates(rice.MustFindBox("templates"))

	work.Waitgroup.Add(4)

	go generators.GenerateBootstrap(work, recipe.Bootstrap)
	go generators.GenerateHTTP(work, recipe.HTTP)
	go generators.GenerateCrud(work, recipe.Crud, recipe.Entities)
	go generators.GenerateREST(work, recipe.Rest, recipe.Entities)

	go func() {
		for generated := range done {
			if generated.Error == nil {
				fmt.Println("Filename: ", generated.Filename)
				fmt.Println(generated.Code)
			} else {
				fmt.Println("Error: ", generated.Error)
			}
			work.Waitgroup.Done()
		}
	}()

	work.Waitgroup.Wait()
	close(done)
}
