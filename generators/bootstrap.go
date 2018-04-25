package generators

import (
	"fmt"

	"github.com/fluxynet/gocipe/util"
)

// GenerateBootstrap returns bootstrap generated code
func GenerateBootstrap(work util.GenerationWork, opts util.BootstrapOpts) error {
	if !opts.Generate {
		work.Done <- util.GeneratedCode{Generator: "GenerateBootstrap", Error: util.ErrorSkip}
	}

	code, err := util.ExecuteTemplate("bootstrap.go.tmpl", opts)

	if err != nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateBootstrap", Error: fmt.Errorf("failed to execute template: %s", err)}
		return err
	}

	work.Done <- util.GeneratedCode{Generator: "GenerateBootstrap", Filename: "bootstrap.go", Code: code}
	return nil
}
