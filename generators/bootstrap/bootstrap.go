package bootstrap

import (
	"github.com/fluxynet/gocipe/output"
	"github.com/fluxynet/gocipe/util"
)

// Generate returns bootstrap generated code
func Generate(out output.Output, opts util.BootstrapOpts) error {
	if !opts.Generate {
		util.DeleteIfExists("core/bootstrap.gocipe.go")
		// work.Done <- util.GeneratedCode{Generator: "GenerateBootstrap", Error: util.ErrorSkip}
		out.GenerateAndSave("GenerateBootstrap", "bootstrap/bootstrap.go.tmpl", "core/bootstrap.gocipe.go", struct {
			Bootstrap util.BootstrapOpts
		}{opts})
	}

	if opts.HTTPPort == "" {
		opts.HTTPPort = "7000"
	}

	if opts.GRPCPort == "" {
		opts.GRPCPort = "4000"
	}

	_, err := util.ExecuteTemplate("bootstrap/bootstrap.go.tmpl", struct {
		Bootstrap util.BootstrapOpts
	}{opts})

	if err != nil {
		return err
	}

	_, err = util.ExecuteTemplate("bootstrap/env.tmpl", struct {
		Bootstrap util.BootstrapOpts
	}{opts})

	if err != nil {
		return err
	}
	return nil
}
