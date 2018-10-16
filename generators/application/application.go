package application

import (
	"os"
	"strings"

	"github.com/fluxynet/gocipe/output"
	"github.com/fluxynet/gocipe/util"
)

// Generate common utility functions
func Generate(out output.Output, recipe *util.Recipe, noSkip bool) {
	out.GenerateAndOverwrite("Scaffold", "", "web/app/dist/.gitkeep", "Front-end production files must compile or be placed here. Delete this file when done.")
	out.GenerateAndOverwrite("Scaffold", "", "web/app/src/services/.gitkeep", "Generated client code will be here.")
	out.GenerateAndOverwrite("Scaffold", "", "services/.gitkeep", "Generated server code and implementation will be here.")
	out.GenerateAndOverwrite("Scaffold", "", "assets/.gitkeep", "Place assets in this folder.")
	out.GenerateAndOverwrite("Scaffold", "", "assets/templates/.gitkeep", "Place templates in this folder.")
	out.GenerateAndOverwrite("Scaffold", "", "assets/web/app/.gitkeep", "Place web assets in this folder.")
	out.GenerateAndOverwrite("Scaffold", "application/gen-service.sh.tmpl", "gen-service.sh", struct{ GeneratePath string }{GeneratePath: "$GOPATH" + strings.TrimPrefix(util.WorkingDir, os.Getenv("GOPATH")) + "/services"})
	out.GenerateAndOverwrite("Scaffold", "application/makefile.tmpl", "Makefile", struct{ AppName string }{util.AppName})
	out.GenerateAndOverwrite("Scaffold", "application/main.go.tmpl", "main.go",
		struct {
			Recipe     *util.Recipe
			ImportPath string
		}{
			Recipe:     recipe,
			ImportPath: util.AppImportPath,
		})
	out.GenerateAndOverwrite("Scaffold", "application/route.go.tmpl", "route.go", struct {
		Bootstrap util.BootstrapOpts
		Admin     util.AdminOpts
	}{recipe.Bootstrap, recipe.Admin})
}
