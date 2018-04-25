package generators

import (
	"fmt"
	"strings"

	"github.com/fluxynet/gocipe/util"
)

// GenerateSchema returns generated database schema creation code
func GenerateSchema(work util.GenerationWork, opts util.SchemaOpts, entities []util.Entity) error {

	work.Waitgroup.Add(len(entities))
	for _, entity := range entities {
		go func(entity util.Entity) {
			var (
				data struct {
					Entity         util.Entity
					ManyManyFields []util.Field
				}
			)

			if entity.Schema == nil {
				entity.Schema = &opts
			}

			data.Entity = entity

			for _, field := range data.Entity.Fields {
				if field.Relationship.Type == util.RelationshipTypeManyMany &&
					strings.Compare(field.Relationship.Target.ThisID, field.Relationship.Target.ThatID) == 1 {
					data.ManyManyFields = append(data.ManyManyFields, field)
				}
			}

			code, err := util.ExecuteTemplate("schema.sql.tmpl", data)

			if err != nil {
				work.Done <- util.GeneratedCode{Generator: "GenerateSchema", Error: err}
			} else if entity.Schema.Aggregate {
				work.Done <- util.GeneratedCode{Generator: "GenerateSchema", Code: code, Filename: "schema.sql"}
			} else {
				work.Done <- util.GeneratedCode{Generator: "GenerateSchema", Code: code, Filename: fmt.Sprintf("schema_%s.sql", entity.Table)}
			}
		}(entity)
	}

	work.Waitgroup.Done()
	return nil
}
