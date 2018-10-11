package schema

import (
	"github.com/fluxynet/gocipe/output"
	"github.com/fluxynet/gocipe/util"
)

// Generate returns generated database schema creation code
func Generate(work util.GenerationWork, r *util.Recipe) error {
	output.GenerateAndSave("schema", "schema/schema.sql.tmpl", "schema/schema.gocipe.sql", r, false)
	work.Waitgroup.Done()
	return nil
}
