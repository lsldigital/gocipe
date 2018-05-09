package generators

import (
	"fmt"

	"github.com/fluxynet/gocipe/util"
)

// GenerateHTTP returns generated code to run an http server
func GenerateHTTP(work util.GenerationWork, opts util.HTTPOpts) error {
	if !opts.Generate {
		work.Done <- util.GeneratedCode{Generator: "GenerateHTTP", Error: util.ErrorSkip}
		return nil
	}

	code, err := util.ExecuteTemplate("http.go.tmpl", opts)

	if err != nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateHTTP", Error: fmt.Errorf("failed to load execute template: %s", err)}
		return err
	}

	work.Done <- util.GeneratedCode{Generator: "GenerateHTTP", Filename: "app/http.go", Code: code}
	return nil
}
