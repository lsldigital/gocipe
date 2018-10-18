package application

import (
	"os"
	"strings"

	"github.com/fluxynet/gocipe/output"
	"github.com/fluxynet/gocipe/util"
)

// Generate common utility functions
func Generate(out *output.Output, recipe *util.Recipe, noSkip bool) {
	out.GenerateAndOverwrite("Scaffold Folder", "", "web/app/dist/.gitkeep", "Front-end production files must compile or be placed here. Delete this file when done.")
	out.GenerateAndOverwrite("Scaffold Folder", "", "web/app/src/services/.gitkeep", "Generated client code will be here.")
	out.GenerateAndOverwrite("Scaffold Folder", "", "services/.gitkeep", "Generated server code and implementation will be here.")
	out.GenerateAndOverwrite("Scaffold Folder", "", "assets/.gitkeep", "Place assets in this folder.")
	out.GenerateAndOverwrite("Scaffold Folder", "", "assets/templates/.gitkeep", "Place templates in this folder.")
	out.GenerateAndOverwrite("Scaffold Folder", "", "assets/web/app/.gitkeep", "Place web assets in this folder.")
	out.GenerateAndOverwrite("Scaffold GenService", "application/gen-service.sh.tmpl", "gen-service.sh", struct{ GeneratePath string }{GeneratePath: "$GOPATH" + strings.TrimPrefix(util.WorkingDir, os.Getenv("GOPATH")) + "/services"})

	if noSkip {
		out.GenerateAndOverwrite("Scaffold Makefile", "application/makefile.tmpl", "Makefile", struct{ AppName string }{util.AppName})
		out.GenerateAndOverwrite("Scaffold Main", "application/main.go.tmpl", "main.go",
			struct {
				Recipe     *util.Recipe
				ImportPath string
			}{
				Recipe:     recipe,
				ImportPath: util.AppImportPath,
			})
		out.GenerateAndOverwrite("Scaffold Route", "application/route.go.tmpl", "route.go", struct {
			Bootstrap util.BootstrapOpts
			Admin     util.AdminOpts
		}{recipe.Bootstrap, recipe.Admin})
	} else {
		out.GenerateAndSave("Scaffold Makefile", "application/makefile.tmpl", "Makefile", struct{ AppName string }{util.AppName})
		out.GenerateAndSave("Scaffold Main", "application/main.go.tmpl", "main.go",
			struct {
				Recipe     *util.Recipe
				ImportPath string
			}{
				Recipe:     recipe,
				ImportPath: util.AppImportPath,
			})
		out.GenerateAndSave("Scaffold Route", "application/route.go.tmpl", "route.go", struct {
			Bootstrap util.BootstrapOpts
			Admin     util.AdminOpts
		}{recipe.Bootstrap, recipe.Admin})
	}
}
