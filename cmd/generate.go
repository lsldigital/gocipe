package cmd

import (
	"log"

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
		output.SetVerbose(verbose)

		rcp, err := recipe.Load()

		if err != nil {
			log.Fatalln("[loadRecipe]", err)
		}

		var outpt output.Output

		//scaffold application layout - synchronously before launching generators
		application.Generate(outpt, rcp, noSkip)

		if generateBootstrap {
			bootstrap.Generate(outpt, rcp)
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

		outpt.ProcessProto()
		outpt.PostProcessGoFiles()
	},
}
