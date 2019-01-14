package application

import (
	"os"
	"strings"

	"github.com/lsldigital/gocipe/output"
	"github.com/lsldigital/gocipe/util"
)

// Generate common utility functions
func Generate(out *output.Output, r *util.Recipe, noSkip bool) {
	out.GenerateAndOverwrite("Scaffold Folder", "", "web/app/dist/.gitkeep", output.WithoutHeader, "Front-end production files must compile or be placed here. Delete this file when done.")
	out.GenerateAndOverwrite("Scaffold Folder", "", "web/app/src/services/.gitkeep", output.WithoutHeader, "Generated client code will be here.")
	out.GenerateAndOverwrite("Scaffold Folder", "", "services/.gitkeep", output.WithoutHeader, "Generated server code and implementation will be here.")
	out.GenerateAndOverwrite("Scaffold Folder", "", "assets/.gitkeep", output.WithoutHeader, "Place assets in this folder.")
	out.GenerateAndOverwrite("Scaffold Folder", "", "assets/templates/.gitkeep", output.WithoutHeader, "Place templates in this folder.")
	out.GenerateAndOverwrite("Scaffold Folder", "", "assets/web/app/.gitkeep", output.WithoutHeader, "Place web assets in this folder.")
	out.GenerateAndOverwrite("Scaffold GenService", "application/gen-service.sh.tmpl", "gen-service.sh", output.WithoutHeader, struct{ GeneratePath string }{GeneratePath: "$GOPATH" + strings.TrimPrefix(util.WorkingDir, os.Getenv("GOPATH")) + "/services"})

	if noSkip {
		out.GenerateAndOverwrite("Scaffold Makefile", "application/makefile.tmpl", "Makefile", output.WithHeader, struct{ AppName string }{util.AppName})
		out.GenerateAndOverwrite("Scaffold Main", "application/main.go.tmpl", "main.go", output.WithHeader, struct{ Recipe *util.Recipe }{r})
		out.GenerateAndOverwrite("Scaffold Route", "application/route.go.tmpl", "route.go", output.WithHeader, struct{ Recipe *util.Recipe }{r})
	} else {
		out.GenerateAndSave("Scaffold Makefile", "application/makefile.tmpl", "Makefile", output.WithHeader, struct{ AppName string }{util.AppName})
		out.GenerateAndSave("Scaffold Main", "application/main.go.tmpl", "main.go", output.WithHeader, struct{ Recipe *util.Recipe }{r})
		out.GenerateAndSave("Scaffold Route", "application/route.go.tmpl", "route.go", output.WithHeader, struct{ Recipe *util.Recipe }{r})
	}

	adminFolder, _ := util.GetAbsPath("web/admin")
	if !util.FileExists(adminFolder) {
		output.GenerateAndSave("Scaffold", "application/preset.json.tmpl", "web/preset.json", output.WithHeader, struct{}{})
		output.CreateVue("admin")
	}
}
