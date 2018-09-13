package cmd

import (
	"log"
	"runtime"
	"sync"

	"github.com/fluxynet/gocipe/generators/application"
	"github.com/fluxynet/gocipe/generators/bootstrap"
	"github.com/fluxynet/gocipe/generators/bread"
	"github.com/fluxynet/gocipe/generators/crud"
	"github.com/fluxynet/gocipe/generators/schema"
	utils "github.com/fluxynet/gocipe/generators/util"
	"github.com/fluxynet/gocipe/generators/vuetify"
	"github.com/fluxynet/gocipe/output"
	"github.com/fluxynet/gocipe/recipe"
	"github.com/fluxynet/gocipe/util"
	"github.com/spf13/cobra"
)

var (
	noSkip            bool
	generateBootstrap bool
	generateSchema    bool
	generateCrud      bool
	generateBread     bool
	generateUtils     bool
	generateVuetify   bool
	verbose           bool
)

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

		if generateBootstrap {
			work.Waitgroup.Add(1)
		}

		if generateSchema {
			work.Waitgroup.Add(1)
		}

		if generateCrud {
			work.Waitgroup.Add(1)
		}

		if generateBread {
			work.Waitgroup.Add(1)
		}

		if generateUtils {
			work.Waitgroup.Add(1)
		}

		if generateVuetify {
			work.Waitgroup.Add(1)
		}

		entities, err := recipe.Preprocess(rcp)
		if err != nil {
			log.Fatalln("preprocessRecipe", err)
		}

		//scaffold application layout - synchronously before launching generators
		application.Generate(rcp, noSkip)

		if generateBootstrap {
			go bootstrap.Generate(work, rcp.Bootstrap)
		}

		if generateSchema {
			go schema.Generate(work, rcp.Schema, entities)
		}

		if generateCrud {
			go crud.Generate(work, rcp.Crud, entities)
		}

		if generateBread {
			go bread.Generate(work, entities)
		}

		if generateUtils {
			go utils.Generate(work)
		}

		if generateVuetify {
			go vuetify.Generate(work, rcp, entities)
		}
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
