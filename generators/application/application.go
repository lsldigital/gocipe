package application

import (
	"os"
	"strings"

	"github.com/fluxynet/gocipe/output"
	"github.com/fluxynet/gocipe/util"
)

// Generate common utility functions
func Generate(bootOpts util.BootstrapOpts, noSkip bool) {
	output.GenerateAndSave("Scaffold", "", "web/app/dist/.gitkeep", "Front-end production files must compile or be placed here. Delete this file when done.", true, noSkip)
	output.GenerateAndSave("Scaffold", "", "web/app/src/services/.gitkeep", "Generated client code will be here.", true, noSkip)
	output.GenerateAndSave("Scaffold", "", "services/.gitkeep", "Generated server code and implementation will be here.", true, noSkip)
	output.GenerateAndSave("Scaffold", "", "assets/.gitkeep", "Place assets in this folder.", true, noSkip)
	output.GenerateAndSave("Scaffold", "", "assets/templates/.gitkeep", "Place templates in this folder.", true, noSkip)
	output.GenerateAndSave("Scaffold", "", "assets/web/app/.gitkeep", "Place web assets in this folder.", true, noSkip)
	output.GenerateAndSave("Scaffold", "application/gen-service.sh.tmpl", "gen-service.sh", struct{ GeneratePath string }{GeneratePath: "$GOPATH" + strings.TrimPrefix(util.WorkingDir, os.Getenv("GOPATH")) + "/services"}, true, noSkip)
	output.GenerateAndSave("Scaffold", "application/main.go.tmpl", "main.go", struct{ Bootstrap util.BootstrapOpts }{bootOpts}, true, noSkip)
	output.AddGoFile("main.go")
}
