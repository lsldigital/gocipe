package generators

import (
	"fmt"
	"strings"

	"github.com/jinzhu/inflection"

	"github.com/fluxynet/gocipe/util"
)

// RelatedTable represents a related table to be created by schema generation
type RelatedTable struct {
	Table    string
	ThisID   string
	ThisType string
	ThatID   string
	ThatType string
}

// RelatedField represents a related field to be added to the table during schema generation
type RelatedField struct {
	Name string
	Type string
}

// GenerateSchema returns generated database schema creation code
func GenerateSchema(work util.GenerationWork, opts util.SchemaOpts, entities map[string]util.Entity) error {

	work.Waitgroup.Add(len(entities))
	for _, entity := range entities {
		go func(entity util.Entity) {
			var (
				data struct {
					Entity        util.Entity
					RelatedFields []RelatedField
					RelatedTables []RelatedTable
				}
			)

			if entity.PrimaryKey == "" {
				entity.PrimaryKey = util.PrimaryKeySerial
			}

			if entity.Schema == nil {
				entity.Schema = &opts
			} else if entity.Schema.Path == "" {
				entity.Schema.Path = opts.Path
			}

			path := entity.Schema.Path
			if path == "" {
				path = "schema"
			}

			data.Entity = entity

			for _, rel := range entity.Relationships {
				related := entities[rel.Entity]
				if rel.Type == util.RelationshipTypeManyOne {
					n := strings.ToLower(related.Name) + "_id"
					t, _ := util.GetPrimaryKeyFieldType(related.PrimaryKey)
					data.RelatedFields = append(data.RelatedFields, RelatedField{Name: n, Type: t})
				} else if rel.Type == util.RelationshipTypeManyMany && strings.Compare(entity.Table, related.Table) > 0 {
					table := inflection.Plural(related.Table) + "_" + inflection.Plural(entity.Table)
					thisID := strings.ToLower(entity.Name) + "_id"
					thatID := strings.ToLower(related.Name) + "_id"
					thisType, _ := util.GetPrimaryKeyFieldType(entity.PrimaryKey)
					thatType, _ := util.GetPrimaryKeyFieldType(related.PrimaryKey)
					data.RelatedTables = append(data.RelatedTables,
						RelatedTable{Table: table, ThisID: thisID, ThisType: thisType, ThatID: thatID, ThatType: thatType},
					)
				}
			}

			code, err := util.ExecuteTemplate("schema.sql.tmpl", data)

			if err != nil {
				work.Done <- util.GeneratedCode{Generator: "GenerateSchema", Error: err}
			} else if entity.Schema.Aggregate {
				work.Done <- util.GeneratedCode{Generator: "GenerateSchema", Code: code, Filename: path + "/schema.gocipe.sql", Aggregate: true, GeneratedHeaderFormat: "-- %s"}
			} else {
				work.Done <- util.GeneratedCode{Generator: "GenerateSchema", Code: code, Filename: fmt.Sprintf("%s/schema_%s.gocipe.sql", path, entity.Table), GeneratedHeaderFormat: "-- %s"}
			}
		}(entity)
	}

	work.Waitgroup.Done()
	return nil
}
