package application

import (
	"os"
	"strings"

	"github.com/fluxynet/gocipe/util"
)

// Generate common utility functions
func Generate(bootOpts util.BootstrapOpts) {
	util.GenerateAndSave("", "web/dist/.gitkeep", "Front-end production files must compile or be placed here. Delete this file when done.", true)
	util.GenerateAndSave("", "web/src/services/.gitkeep", "Generated client code will be here.", true)
	util.GenerateAndSave("", "services/.gitkeep", "Generated server code and implementation will be here.", true)
	util.GenerateAndSave("", "assets/.gitkeep", "Place assets in this folder.", true)
	util.GenerateAndSave("", "assets/templates/.gitkeep", "Place templates in this folder.", true)
	util.GenerateAndSave("application/gen-service.sh.tmpl", "gen-service.sh", struct{ GeneratePath string }{GeneratePath: "$GOPATH" + strings.TrimPrefix(util.WorkingDir, os.Getenv("GOPATH")) + "/services"}, true)
	util.GenerateAndSave("application/main.go.tmpl", "main.go", struct{ Bootstrap util.BootstrapOpts }{bootOpts}, true)
}
