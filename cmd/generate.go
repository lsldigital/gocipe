package cmd

import (
	"log"
	"runtime"
	// "sync"

	"github.com/fluxynet/gocipe/generators/admin"
	"github.com/fluxynet/gocipe/generators/application"
	"github.com/fluxynet/gocipe/generators/auth"
	"github.com/fluxynet/gocipe/generators/bootstrap"
	"github.com/fluxynet/gocipe/generators/crud"
	"github.com/fluxynet/gocipe/generators/schema"
	utils "github.com/fluxynet/gocipe/generators/util"
	"github.com/fluxynet/gocipe/generators/vuetify"
	"github.com/fluxynet/gocipe/output"
	"github.com/fluxynet/gocipe/recipe"
	// "github.com/fluxynet/gocipe/util"
	"github.com/spf13/cobra"
)

var (
	noSkip            bool
	generateBootstrap bool
	generateSchema    bool
	generateCrud      bool
	generateAdmin     bool
	generateAuth      bool
	generateUtils     bool
	generateVuetify   bool
	verbose           bool
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"init"},
	Run: func(cmd *cobra.Command, args []string) {
		runtime.GOMAXPROCS(runtime.NumCPU())
		// work := util.GenerationWork{
		// 	Waitgroup: new(sync.WaitGroup),
		// 	Done:      make(chan util.GeneratedCode),
		// }

		output.SetVerbose(verbose)

		rcp, err := recipe.Load()

		if err != nil {
			log.Fatalln("[loadRecipe]", err)
		}

		// if generateBootstrap {
		// 	work.Waitgroup.Add(1)
		// }

		// if generateSchema {
		// 	work.Waitgroup.Add(1)
		// }

		// if generateCrud {
		// 	work.Waitgroup.Add(1)
		// }

		// if generateAdmin {
		// 	work.Waitgroup.Add(1)
		// }

		// if generateAuth {
		// 	work.Waitgroup.Add(1)
		// }

		// if generateUtils {
		// 	work.Waitgroup.Add(1)
		// }

		// if generateVuetify {
		// 	work.Waitgroup.Add(1)
		// }

		var outpt output.Output

		//scaffold application layout - synchronously before launching generators
		application.Generate(outpt, rcp, noSkip)

		


		if generateBootstrap {
			bootstrap.Generate(outpt, rcp.Bootstrap)
		}

		if generateSchema {
		   schema.Generate(outpt, rcp)
		}

		if generateCrud {
			crud.Generate(outpt, rcp)
		}

		if generateAdmin {
			 admin.Generate(outpt, rcp)
		}

		if generateAdmin {
			 auth.Generate(outpt)
		}

		if generateUtils {
			 utils.Generate(outpt)
		}

		if generateVuetify {
			 vuetify.Generate(outpt, rcp)
		}
		// go generators.GenerateHTTP(work, recipe.HTTP)
		// go generators.GenerateREST(work, recipe.Rest, recipe.Entities)

		// var wg sync.WaitGroup
		// wg.Add(1)

		// go output.Process(work, noSkip)
		
		//ask yousuf for this one 
		//  outpt.ProcessGoFiles(noSkip)

		// work.Waitgroup.Wait()
		// close(work.Done)
		// wg.Wait()

		output.ProcessProto()
		outpt.PostProcessGoFiles()
		output.WriteLog()
	},
}
