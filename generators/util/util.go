package util

import (
	"fmt"

	"github.com/fluxynet/gocipe/util"
)

// Generate common utility functions
func Generate(work util.GenerationWork) {
	work.Waitgroup.Add(1)
	models, err := util.ExecuteTemplate("util/rice.go.tmpl", struct{}{})
	if err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateUtilRice", Code: models, Filename: "util/rice.gocipe.go"}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateUtilRice", Error: fmt.Errorf("failed to load execute template: %s", err)}
	}

	models, err = util.ExecuteTemplate("util/util.go.tmpl", struct{}{})
	if err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateUtil", Code: models, Filename: "util/util.gocipe.go"}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateUtil", Error: fmt.Errorf("failed to load execute template: %s", err)}
	}
}
