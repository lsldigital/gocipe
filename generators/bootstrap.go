package generators

import (
	"fmt"

	"github.com/fluxynet/gocipe/util"
)

// GenerateBootstrap returns bootstrap generated code
func GenerateBootstrap(work util.GenerationWork, opts util.BootstrapOpts, httpOpts util.HTTPOpts) error {
	if !opts.Generate {
		work.Done <- util.GeneratedCode{Generator: "GenerateBootstrap", Error: util.ErrorSkip}
	}

	if httpOpts.Port == "" {
		httpOpts.Port = "8080"
	}

	code, err := util.ExecuteTemplate("bootstrap.go.tmpl", struct {
		Bootstrap util.BootstrapOpts
		HTTP      util.HTTPOpts
	}{opts, httpOpts})

	if err != nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateBootstrap", Error: fmt.Errorf("failed to execute template: %s", err)}
		return err
	}

	work.Done <- util.GeneratedCode{Generator: "GenerateBootstrap", Filename: "bootstrap.go", Code: code}
	return nil
}
