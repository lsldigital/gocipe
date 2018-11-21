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

	out.GenerateAndOverwrite(
		"Generate ServiceData Proto", "data/service_data.proto.tmpl", "proto/service_data.proto", output.WithHeader,
		struct {
			Entities   []util.Entity
			ImportPath string
		}{
			Entities:   r.Entities,
			ImportPath: r.ImportPath,
		},
	)

	out.GenerateAndOverwrite(
		"Generate ServiceData", "data/service_data.go.tmpl", "services/data/service_data.gocipe.go", output.WithHeader,
		struct {
			Entities   []util.Entity
			ImportPath string
		}{
			Entities:   r.Entities,
			ImportPath: r.ImportPath,
		},
	)
}
