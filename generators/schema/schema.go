package schema

import (
	"github.com/fluxynet/gocipe/output"
	"github.com/fluxynet/gocipe/util"
)

// Generate returns generated database schema creation code
func Generate(out output.Output, r *util.Recipe) {
	out.GenerateAndOverwrite("schema", "schema/schema.sql.tmpl", "schema/schema.gocipe.sql", r)
}
