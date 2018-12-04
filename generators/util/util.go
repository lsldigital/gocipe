package util

import (
	"github.com/lsldigital/gocipe/output"
	utils "github.com/lsldigital/gocipe/util"
)

// Generate common utility functions
func Generate(out *output.Output, r *utils.Recipe) {
	data := struct{ AppImportPath string }{r.ImportPath}

	out.GenerateAndOverwrite("GenerateUtil Rice", "util/rice.go.tmpl", "util/rice.gocipe.go", output.WithHeader, nil)
	out.GenerateAndOverwrite("GenerateUtil Credentials", "util/credentials.go.tmpl", "util/credentials/credentials.gocipe.go", output.WithHeader, nil)
	out.GenerateAndOverwrite("GenerateUtil", "util/util.go.tmpl", "util/util.gocipe.go", output.WithHeader, nil)
	out.GenerateAndOverwrite("GenerateUtil Ws", "util/ws.go.tmpl", "util/web/ws.gocipe.go", output.WithHeader, nil)
	out.GenerateAndOverwrite("GenerateUtil Fileupload", "util/files.go.tmpl", "util/files/files.gocipe.go", output.WithHeader, data)
	out.GenerateAndOverwrite("GenerateUtil Imagist", "util/imagist.go.tmpl", "util/imagist/imagist.gocipe.go", output.WithHeader, data)
	out.GenerateAndOverwrite("GenerateUtil Web", "util/web.go.tmpl", "util/web/web.gocipe.go", output.WithHeader, struct {
		ImportPath string
	}{r.ImportPath})

	if r.Decks.Generate {
		out.GenerateAndOverwrite("GenerateUtil Decks", "util/decks.go.tmpl", "util/decks/decks.gocipe.go", output.WithHeader, struct {
			AppImportPath string
			Decks         []utils.DeckOpts
		}{r.ImportPath, r.Decks.Decks})
	}
}
