package generators

import (
	"fmt"
	"strings"
)

// GenerateSchema returns generated database schema creation code
func GenerateSchema(work GenerationWork, opts SchemaOpts, entities []Entity) error {

	work.Waitgroup.Add(len(entities))
	for _, entity := range entities {
		go func(entity Entity) {
			var (
				data struct {
					Entity         Entity
					ManyManyFields []Field
				}
			)

			if entity.Schema == nil {
				entity.Schema = &opts
			}

			data.Entity = entity

			for _, field := range data.Entity.Fields {
				if field.Relationship.Type == RelationshipTypeManyMany &&
					strings.Compare(field.Relationship.Target.ThisID, field.Relationship.Target.ThatID) == 1 {
					data.ManyManyFields = append(data.ManyManyFields, field)
				}
			}

			code, err := ExecuteTemplate("schema.sql.tmpl", data)

			if err != nil {
				work.Done <- GeneratedCode{Generator: "GenerateSchema", Error: err}
			} else if entity.Schema.Aggregate {
				work.Done <- GeneratedCode{Generator: "GenerateSchema", Code: code, Filename: "schema.sql"}
			} else {
				work.Done <- GeneratedCode{Generator: "GenerateSchema", Code: code, Filename: fmt.Sprintf("schema_%s.sql", entity.Table)}
			}
		}(entity)
	}

	work.Waitgroup.Done()
	return nil
}
