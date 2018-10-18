package auth

import (
	"github.com/fluxynet/gocipe/output"
)

// Generate auth code
func Generate(out *output.Output) {
	out.GenerateAndOverwrite("Auth", "auth/auth.go.tmpl", "auth/auth.gocipe.go", output.WithHeader, nil)
}
