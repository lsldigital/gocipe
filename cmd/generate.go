package cmd

import (
	"log"

	"github.com/fluxynet/gocipe/generators/admin"
	"github.com/fluxynet/gocipe/generators/application"
	"github.com/fluxynet/gocipe/generators/auth"
	"github.com/fluxynet/gocipe/generators/bootstrap"
	"github.com/fluxynet/gocipe/generators/crud"
	"github.com/fluxynet/gocipe/generators/data"
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
