package cmd

import (
	"log"

	"github.com/lsldigital/gocipe/generators/admin"
	"github.com/lsldigital/gocipe/generators/application"
	"github.com/lsldigital/gocipe/generators/auth"
	"github.com/lsldigital/gocipe/generators/bootstrap"
	"github.com/lsldigital/gocipe/generators/crud"
	"github.com/lsldigital/gocipe/generators/data"
	"github.com/lsldigital/gocipe/generators/schema"
	utils "github.com/lsldigital/gocipe/generators/util"
	"github.com/lsldigital/gocipe/generators/vuetify"
	"github.com/lsldigital/gocipe/output"
	"github.com/lsldigital/gocipe/util"

	"github.com/spf13/cobra"
)

var (
	overwrite         bool
	generateBootstrap bool
	generateSchema    bool
	generateCrud      bool
	generateAdmin     bool
	generateData      bool
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

		r, err := util.LoadRecipe("gocipe.json")

		if err != nil {
			log.Fatalln("[loadRecipe]", err)
		}

		application.Generate(out, r, overwrite)

		if generateBootstrap {
			bootstrap.Generate(out, r)
		}

		if generateData {
			data.Generate(out, r)
		}

		if generateSchema {
			schema.Generate(out, r)
		}

		if generateCrud {
			crud.Generate(out, r)
		}

		if generateAdmin {
			admin.Generate(out, r)
		}

		if generateAuth {
			auth.Generate(out, r)
		}

		if generateUtils {
			utils.Generate(out, r)
		}

		if generateVuetify {
			vuetify.Generate(out, r)
		}

		out.ProcessProto()
		out.PostProcessGoFiles(r)
		out.Write("gocipe.log")
	},
}
