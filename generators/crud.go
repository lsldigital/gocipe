package generators

import (
	"fmt"
	"strings"
)

// GenerateCrud returns generated code to run an http server
func GenerateCrud(work GenerationWork, opts CrudOpts, entities []Entity) error {
	if !opts.Create && !opts.Read && !opts.ReadList && !opts.Update && !opts.Delete {
		work.Done <- GeneratedCode{Generator: "GenerateCrud", Error: ErrorSkip}
		return nil
	}

	work.Waitgroup.Add(len(entities))
	for _, entity := range entities {
		go func(entity Entity) {
			var (
				data struct {
					Package string
					Entity  Entity

					SQLFieldsSelect string
					SQLFieldsUpdate string
					SQLFieldsInsert string
					SQLPlaceholders string

					StructFieldsSelect string
					StructFieldsUpdate string
					StructFieldsInsert string

					Joins         string
					JoinVarsDecl  []string
					JoinVarsAssgn []string

					BeforeUpdate []string
					BeforeInsert []string

					HasRelationshipManyMany bool
					ManyManyFields          []Field

					Hooks CrudHooks
				}

				sqlfieldsSelect []string
				sqlfieldsUpdate []string
				sqlfieldsInsert []string
				sqlPlaceholders []string

				structFieldsSelect  []string
				structFieldsUpdate  []string
				structFieldsInsert  []string
				sqlPlaceholderCount = 1

				joins     []string
				joinCount int
			)

			for _, field := range entity.Fields {
				if field.Relationship.Type == "" {
					sqlfieldsSelect = append(sqlfieldsSelect, fmt.Sprintf("t.%s", field.Schema.Field))
					structFieldsSelect = append(structFieldsSelect, fmt.Sprintf("entity.%s", field.Property.Name))

					if field.Property.Name != "ID" {
						sqlPlaceholders = append(sqlPlaceholders, fmt.Sprintf("$%d", sqlPlaceholderCount))
						sqlfieldsUpdate = append(sqlfieldsUpdate, fmt.Sprintf("%s = $%d", field.Schema.Field, sqlPlaceholderCount))
						sqlfieldsInsert = append(sqlfieldsInsert, fmt.Sprintf("$%d", sqlPlaceholderCount))

						structFieldsInsert = append(structFieldsInsert, fmt.Sprintf("*entity.%s", field.Property.Name))
						structFieldsUpdate = append(structFieldsUpdate, fmt.Sprintf("*entity.%s", field.Property.Name))
					}

					if field.Property.Name == "CreatedAt" {
						data.BeforeInsert = append(data.BeforeInsert, "*entity.CreatedAt = time.Now()")
					} else if field.Property.Name == "UpdatedAt" {
						data.BeforeInsert = append(data.BeforeInsert, "*entity.UpdatedAt = time.Now()")
						data.BeforeUpdate = append(data.BeforeUpdate, "*entity.UpdatedAt = time.Now()")
					}

					sqlPlaceholderCount++
				} else {
					joins = append(joins,
						fmt.Sprintf("%s jt%d ON (t.%s = jt%d.%s)",
							field.Relationship.Target.Table,
							joinCount,
							field.Relationship.Target.ThisID,
							joinCount,
							field.Relationship.Target.ThatID))

					data.JoinVarsDecl = append(data.JoinVarsDecl, fmt.Sprintf("j%d int64", joinCount))
					data.JoinVarsAssgn = append(data.JoinVarsAssgn, fmt.Sprintf("entity.%s = append(entity.%s, j%d)", field.Property.Name, field.Property.Name, joinCount))
					sqlfieldsSelect = append(sqlfieldsSelect, fmt.Sprintf("jt%d.%s", joinCount, field.Relationship.Target.ThatID))
					structFieldsSelect = append(structFieldsSelect, fmt.Sprintf("&j%d, ", joinCount))
					joinCount++

					if field.Relationship.Type == "many-many" {
						data.HasRelationshipManyMany = true
						data.ManyManyFields = append(data.ManyManyFields, field)
					}
				}
			}

			data.Entity = entity
			data.Package = strings.ToLower(entity.Name)
			data.SQLFieldsSelect = strings.Join(sqlfieldsSelect, ", ")
			data.SQLFieldsUpdate = strings.Join(sqlfieldsUpdate, ", ")
			data.SQLFieldsInsert = strings.Join(sqlfieldsInsert, ", ")
			data.SQLPlaceholders = strings.Join(sqlPlaceholders, ", ")

			data.StructFieldsSelect = strings.Join(structFieldsSelect, ", ")
			data.StructFieldsUpdate = strings.Join(structFieldsUpdate, ", ")
			data.StructFieldsInsert = strings.Join(structFieldsInsert, ", ")

			if joinCount > 0 {
				data.Joins = "INNER JOIN " + strings.Join(joins, " INNER JOIN ") + " "
			}

			code, err := ExecuteTemplate("crud.go.tmpl", data)

			if err != nil {
				work.Done <- GeneratedCode{Generator: "GenerateCRUD", Error: fmt.Errorf("failed to load execute template: %s", err)}
			}

			work.Done <- GeneratedCode{Generator: "GeneratorCRUD", Code: code, Filename: fmt.Sprintf("models/%s/%s_crud.go", data.Package, data.Package)}
		}(entity)
	}

	work.Waitgroup.Done()
	return nil
}
