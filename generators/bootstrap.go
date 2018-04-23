package generators

import "fmt"

// GenerateBootstrap returns bootstrap generated code
func GenerateBootstrap(work GenerationWork, opts BootstrapOpts) error {
	if !opts.Generate {
		work.Done <- GeneratedCode{Generator: "GenerateBootstrap", Error: ErrorSkip}
	}

	code, err := ExecuteTemplate("bootstrap.go.tmpl", opts)

	if err != nil {
		work.Done <- GeneratedCode{Generator: "GenerateBootstrap", Error: fmt.Errorf("failed to execute template: %s", err)}
		return err
	}

	work.Done <- GeneratedCode{Generator: "GenerateBootstrap", Filename: "bootstrap.go", Code: code}
	return nil
}
