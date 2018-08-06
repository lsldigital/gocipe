package cmd

import (
	"log"
	"runtime"
	"sync"

	"github.com/fluxynet/gocipe/generators"
	"github.com/fluxynet/gocipe/generators/application"
	"github.com/fluxynet/gocipe/generators/bread"
	"github.com/fluxynet/gocipe/generators/crud"
	utils "github.com/fluxynet/gocipe/generators/util"
	"github.com/fluxynet/gocipe/generators/vuetify"
	"github.com/fluxynet/gocipe/output"
	"github.com/fluxynet/gocipe/recipe"
	"github.com/fluxynet/gocipe/util"
	"github.com/spf13/cobra"
)

var noSkip bool

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"init"},
	Run: func(cmd *cobra.Command, args []string) {
		runtime.GOMAXPROCS(runtime.NumCPU())
		work := util.GenerationWork{
			Waitgroup: new(sync.WaitGroup),
			Done:      make(chan util.GeneratedCode),
		}

		rcp, err := recipe.Load()

		if err != nil {
			log.Fatalln("[loadRecipe]", err)
		}

		work.Waitgroup.Add(6)

		entities, err := recipe.Preprocess(rcp)
		if err != nil {
			log.Fatalln("preprocessRecipe", err)
		}

		//scaffold application layout - synchronously before launching generators
		application.Generate(rcp.Bootstrap, noSkip)

		go generators.GenerateBootstrap(work, rcp.Bootstrap)
		go crud.Generate(work, rcp.Crud, entities)
		go bread.Generate(work, entities)
		go generators.GenerateSchema(work, rcp.Schema, entities)
		go utils.Generate(work, rcp.Bootstrap)
		go vuetify.Generate(work, rcp.Vuetify, rcp.Entities)
		// go generators.GenerateHTTP(work, recipe.HTTP)
		// go generators.GenerateREST(work, recipe.Rest, recipe.Entities)

		var wg sync.WaitGroup
		wg.Add(1)

		go output.Process(&wg, work, noSkip)

		work.Waitgroup.Wait()
		close(work.Done)
		wg.Wait()

		output.ProcessProto()
		output.PostProcessGoFiles()
		output.WriteLog()
	},
}
