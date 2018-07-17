package application

import (
	"fmt"
	"os"
	"strings"

	"github.com/fluxynet/gocipe/util"
)

// Generate common utility functions
func Generate(work util.GenerationWork, bootOpts util.BootstrapOpts) {
	work.Waitgroup.Add(1)
	work.Done <- util.GeneratedCode{Generator: "GenerateDirs", Code: "Front-end production files must compile or be placed here. Delete this file when done.", Filename: "web/dist/.gitkeep", NoOverwrite: true}

	work.Waitgroup.Add(1)
	work.Done <- util.GeneratedCode{Generator: "GenerateDirs", Code: "Generated client code will be here.", Filename: "web/src/services/.gitkeep", NoOverwrite: true}

	work.Waitgroup.Add(1)
	work.Done <- util.GeneratedCode{Generator: "GenerateDirs", Code: "Generated server code and implementation will be here.", Filename: "services/.gitkeep", NoOverwrite: true}

	if bootOpts.Assets {
		work.Waitgroup.Add(1)
		work.Done <- util.GeneratedCode{Generator: "GenerateAssetsDir", Code: "Place assets in this folder.", Filename: "assets/.gitkeep", NoOverwrite: true}

		work.Waitgroup.Add(1)
		work.Done <- util.GeneratedCode{Generator: "GenerateTemplatesDir", Code: "Place templates in assets folder.", Filename: "assets/templates/.gitkeep", NoOverwrite: true}

		work.Waitgroup.Add(1)
		work.Done <- util.GeneratedCode{Generator: "GenerateWebDir", Code: "Place web in assets folder.", Filename: "assets/web/.gitkeep", NoOverwrite: true}

		work.Waitgroup.Add(1)
		work.Done <- util.GeneratedCode{Generator: "GenerateAppDir", Code: "Place app in web folder.", Filename: "assets/web/app/.gitkeep", NoOverwrite: true}
	}

	work.Waitgroup.Add(1)
	genservice, err := util.ExecuteTemplate("application/gen-service.sh.tmpl", struct {
		GeneratePath string
	}{
		GeneratePath: "$GOPATH" + strings.TrimPrefix(util.WorkingDir, os.Getenv("GOPATH")) + "/services",
	})
	if err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateGenService", Code: genservice, Filename: "gen-service.sh", NoOverwrite: false, GeneratedHeaderFormat: util.NoHeaderFormat}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateGenService", Error: fmt.Errorf("failed to load execute template: %s", err)}
	}

	work.Waitgroup.Add(1)
	models, err := util.ExecuteTemplate("application/main.go.tmpl", struct {
		Bootstrap util.BootstrapOpts
	}{bootOpts})
	if err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateMain", Code: models, Filename: "main.go", NoOverwrite: true}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateMain", Error: fmt.Errorf("failed to load execute template: %s", err)}
	}

	makefile, err := util.ExecuteTemplate("application/makefile.tmpl", struct {
		Bootstrap util.BootstrapOpts
	}{bootOpts})
	if err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateMakefile", Code: makefile, Filename: "Makefile", NoOverwrite: true}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateMakefile", Error: fmt.Errorf("failed to load execute template: %s", err)}
	}
}
