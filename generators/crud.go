package generators

import (
	"fmt"
	"strings"

	"github.com/fluxynet/gocipe/util"
)

// GenerateCrud returns generated code to run an http server
func GenerateCrud(work util.GenerationWork, opts util.CrudOpts, entities []util.Entity) error {
	work.Waitgroup.Add(len(entities) * 3) //2 jobs to be waited upon for each thread - entity.go,  entity_crud.go and entity_crud_hooks.go generation

	for _, entity := range entities {
		go func(entity util.Entity) {
			var (
				data struct {
					Package string
					Entity  util.Entity

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
					ManyManyFields          []util.Field
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

			if entity.Crud == nil {
				entity.Crud = &opts
			}

			if !entity.Crud.Create && !entity.Crud.Read && !entity.Crud.ReadList && !entity.Crud.Update && !entity.Crud.Delete {
				work.Done <- util.GeneratedCode{Generator: "GenerateCrud", Error: util.ErrorSkip}
			}

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

					if field.Relationship.Type == util.RelationshipTypeManyMany {
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

			structure, err := util.ExecuteTemplate("crud_structure.go.tmpl", struct {
				Entity  util.Entity
				Package string
			}{entity, data.Package})
			if err == nil {
				work.Done <- util.GeneratedCode{Generator: "GenerateCRUDModel", Code: structure, Filename: fmt.Sprintf("models/%s/%s.go", data.Package, data.Package)}
			} else {
				work.Done <- util.GeneratedCode{Generator: "GenerateCRUDModel", Error: fmt.Errorf("failed to load execute template: %s", err)}
			}

			code, err := util.ExecuteTemplate("crud.go.tmpl", data)
			if err == nil {
				work.Done <- util.GeneratedCode{Generator: "GenerateCRUD", Code: code, Filename: fmt.Sprintf("models/%s/%s_crud.go", data.Package, data.Package)}
			} else {
				work.Done <- util.GeneratedCode{Generator: "GenerateCRUD", Error: fmt.Errorf("failed to load execute template: %s", err)}
			}

			if entity.Crud.Hooks.PreCreate || entity.Crud.Hooks.PostCreate || entity.Crud.Hooks.PreRead || entity.Crud.Hooks.PostRead || entity.Crud.Hooks.PreList || entity.Crud.Hooks.PostList || entity.Crud.Hooks.PreUpdate || entity.Crud.Hooks.PostUpdate || entity.Crud.Hooks.PreDelete || entity.Crud.Hooks.PostDelete {
				hooks, e := util.ExecuteTemplate("crud_hooks.go.tmpl", struct {
					Hooks   util.CrudHooks
					Name    string
					Package string
				}{entity.Crud.Hooks, entity.Name, data.Package})

				if e == nil {
					work.Done <- util.GeneratedCode{Generator: "GenerateCRUDHooks", Code: hooks, Filename: fmt.Sprintf("models/%s/%s_crud_hooks.go", data.Package, data.Package), NoOverwrite: true}
				} else {
					work.Done <- util.GeneratedCode{Generator: "GenerateCRUDHooks", Error: e}
				}
			} else {
				work.Done <- util.GeneratedCode{Generator: "GenerateCRUDHooks", Error: util.ErrorSkip}
			}
		}(entity)
	}

	code, err := util.ExecuteTemplate("crud_filters.go.tmpl", struct{}{})
	if err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateCRUD", Code: code, Filename: "models/filters.go"}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateCRUD", Error: fmt.Errorf("failed to load execute template: %s", err)}
	}

	return err
}
