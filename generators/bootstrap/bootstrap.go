package bootstrap

import (
	"github.com/fluxynet/gocipe/output"
	"github.com/fluxynet/gocipe/util"
)

// Generate returns bootstrap generated code
func Generate(out output.Output, r *util.Recipe) error {

	if !r.Bootstrap.Generate {
		util.DeleteIfExists("core/bootstrap.gocipe.go")
	}

	if r.Bootstrap.HTTPPort == "" {
		r.Bootstrap.HTTPPort = "7000"
	}

	if r.Bootstrap.GRPCPort == "" {
		r.Bootstrap.GRPCPort = "4000"
	}

	out.GenerateAndSave("GenerateBootstrap", "bootstrap/bootstrap.go.tmpl", "core/bootstrap.gocipe.go", struct {
		Bootstrap util.BootstrapOpts
	}{r.Bootstrap})

	out.GenerateAndSave("GenerateBootstrap Env", "bootstrap/env.tmpl", ".env.dist", struct {
		Bootstrap util.BootstrapOpts
	}{r.Bootstrap})

	return nil
}
