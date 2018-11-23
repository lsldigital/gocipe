package data

import (
	"github.com/fluxynet/gocipe/output"
	"github.com/fluxynet/gocipe/util"
)

// Generate returns generated code for a Data service
func Generate(out *output.Output, r *util.Recipe) {

	// TODO: define in recipe / gocipe.json
	// if !r.Data.Generate {
	// 	return
	// }

	data := struct {
		Entities       []util.Entity
		ImportPath     string
		DecksGenerated bool
	}{
		Entities:       r.Entities,
		ImportPath:     r.ImportPath,
		DecksGenerated: r.Decks.Generate,
	}

	out.GenerateAndOverwrite("Generate ServiceData Proto", "data/service_data.proto.tmpl", "proto/service_data.proto", output.WithHeader, data)
	out.GenerateAndOverwrite("Generate ServiceData", "data/service_data.go.tmpl", "services/data/service_data.gocipe.go", output.WithHeader, data)
}
