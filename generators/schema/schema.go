package schema

import (
	"github.com/lsldigital/gocipe/output"
	"github.com/lsldigital/gocipe/util"
)

// Generate returns generated database schema creation code
func Generate(out *output.Output, r *util.Recipe) {
	out.GenerateAndOverwrite("GenerateSchema", "schema/schema.sql.tmpl", "schema/schema.gocipe.sql", output.WithHeader, r)
}
