package auth

import (
	"github.com/fluxynet/gocipe/output"
	"github.com/fluxynet/gocipe/util"
)

// Generate auth code
func Generate(work util.GenerationWork) {
	output.GenerateAndSave("Auth", "auth/auth.go.tmpl", "auth/auth.gocipe.go", nil, false)
	work.Waitgroup.Done()
}
