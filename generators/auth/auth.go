package auth

import (
	"github.com/lsldigital/gocipe/output"
	"github.com/lsldigital/gocipe/util"
)

// Generate auth code
func Generate(out *output.Output, r *util.Recipe) {
	if !r.Admin.Auth.Generate {
		return
	}

	out.GenerateAndOverwrite("GenerateAuth", "auth/auth.go.tmpl", "auth/auth.gocipe.go", output.WithHeader, nil)
}
