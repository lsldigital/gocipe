package util

import (
	"fmt"

	"github.com/fluxynet/gocipe/util"
)

// Generate common utility functions
func Generate(work util.GenerationWork) {
	models, err := util.ExecuteTemplate("util/util.go.tmpl", struct{}{})
	if err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateUtil", Code: models, Filename: "util/util.gocipe.go"}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateUtil", Error: fmt.Errorf("failed to load execute template: %s", err)}
	}
}
