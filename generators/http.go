package generators

import "fmt"

// GenerateHTTP returns generated code to run an http server
func GenerateHTTP(work GenerationWork, opts HTTPOpts) error {
	if !opts.Generate {
		work.Done <- GeneratedCode{Generator: "GenerateHTTP", Error: ErrorSkip}
		return nil
	}

	code, err := ExecuteTemplate("http.go.tmpl", opts)

	if err != nil {
		work.Done <- GeneratedCode{Generator: "GenerateHTTP", Error: fmt.Errorf("failed to load execute template: %s", err)}
		return err
	}

	work.Done <- GeneratedCode{Generator: "GenerateHTTP", Filename: "http.go", Code: code}
	return nil
}
