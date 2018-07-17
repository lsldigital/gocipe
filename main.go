package main

import (
	"flag"
	"log"
	"runtime"
	"sync"

	rice "github.com/GeertJohan/go.rice"
	"github.com/fluxynet/gocipe/generators"
	"github.com/fluxynet/gocipe/generators/application"
	"github.com/fluxynet/gocipe/generators/bread"
	"github.com/fluxynet/gocipe/generators/crud"
	utils "github.com/fluxynet/gocipe/generators/util"
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
		log.Fatalln("[loadRecipe]", err)
	}

	util.SetTemplates(rice.MustFindBox("templates"))

	work.Waitgroup.Add(7)

	entities, err := preprocessRecipe(recipe)
	if err != nil {
		log.Fatalln("preprocessRecipe", err)
	}

	go generators.GenerateBootstrap(work, recipe.Bootstrap)
	// go generators.GenerateHTTP(work, recipe.HTTP)
	go crud.Generate(work, recipe.Crud, entities)
	go bread.Generate(work, entities)
	// go generators.GenerateREST(work, recipe.Rest, recipe.Entities)
	go generators.GenerateSchema(work, recipe.Schema, entities)
	go generators.GenerateVuetify(work, recipe.Rest, recipe.Vuetify, recipe.Entities)
	go application.Generate(work, recipe.Bootstrap)
	go utils.Generate(work)

	var wg sync.WaitGroup
	wg.Add(1)

	go processOutput(&wg, work, toolset, *noSkip)

	work.Waitgroup.Wait()
	close(work.Done)
	wg.Wait()

	postProcessProtofiles(toolset)
}
