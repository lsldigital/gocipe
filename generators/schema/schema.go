package schema

import (
	"github.com/fluxynet/gocipe/output"
	"github.com/fluxynet/gocipe/util"
)

// Generate returns generated database schema creation code
func Generate(work util.GenerationWork, opts util.SchemaOpts, entities map[string]util.Entity) error {
	output.GenerateAndSave("schema", "schema/schema.sql.tmpl", "schema/schema.gocipe.sql", entities, false)
	work.Waitgroup.Done()
	return nil
}
