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
	"github.com/fluxynet/gocipe/util"

	"github.com/spf13/cobra"
)

var (
	overwrite         bool
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
		out := output.New(verbose)

		rcp, err := util.LoadRecipe()

		if err != nil {
			log.Fatalln("[loadRecipe]", err)
		}

		//scaffold application layout - synchronously before launching generators
		application.Generate(out, rcp, overwrite)

		if generateBootstrap {
			bootstrap.Generate(out, rcp)
		}

		if generateSchema {
			schema.Generate(out, rcp)
		}

		if generateCrud {
			crud.Generate(out, rcp)
		}

		if generateAdmin {
			admin.Generate(out, rcp)
		}

		if generateAuth {
			auth.Generate(out)
		}

		if generateUtils {
			utils.Generate(out, rcp)
		}

		if generateVuetify {
			vuetify.Generate(out, rcp)
		}

		out.ProcessProto()
		out.PostProcessGoFiles()
		out.Write("gocipe.log")
	},
}
