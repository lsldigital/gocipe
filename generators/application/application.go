package application

import (
	"os"
	"strings"

	"github.com/fluxynet/gocipe/output"
	"github.com/fluxynet/gocipe/util"
)

// Generate common utility functions
func Generate(recipe *util.Recipe, noSkip bool) {
	output.GenerateAndSave("Scaffold", "", "web/app/dist/.gitkeep", "Front-end production files must compile or be placed here. Delete this file when done.", false)
	output.GenerateAndSave("Scaffold", "", "web/app/src/services/.gitkeep", "Generated client code will be here.", false)
	output.GenerateAndSave("Scaffold", "", "services/.gitkeep", "Generated server code and implementation will be here.", false)
	output.GenerateAndSave("Scaffold", "", "assets/.gitkeep", "Place assets in this folder.", false)
	output.GenerateAndSave("Scaffold", "", "assets/templates/.gitkeep", "Place templates in this folder.", false)
	output.GenerateAndSave("Scaffold", "", "assets/web/app/.gitkeep", "Place web assets in this folder.", false)
	output.GenerateAndSave("Scaffold", "application/gen-service.sh.tmpl", "gen-service.sh", struct{ GeneratePath string }{GeneratePath: "$GOPATH" + strings.TrimPrefix(util.WorkingDir, os.Getenv("GOPATH")) + "/services"}, false)
	output.GenerateAndSave("Scaffold", "application/makefile.tmpl", "Makefile", struct{ AppName string }{util.AppName}, !noSkip)
	output.GenerateAndSave("Scaffold", "application/main.go.tmpl", "main.go",
		struct {
			Recipe     *util.Recipe
			ImportPath string
		}{
			Recipe:     recipe,
			ImportPath: util.AppImportPath,
		}, !noSkip)
	output.GenerateAndSave("Scaffold", "application/route.go.tmpl", "route.go", struct {
		Bootstrap util.BootstrapOpts
		Admin     util.AdminOpts
	}{recipe.Bootstrap, recipe.Admin}, !noSkip)
}
