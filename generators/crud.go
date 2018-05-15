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

					SQLFieldsSelectGet    string
					SQLFieldsSelectList   string
					SQLFieldsUpdate       string
					SQLFieldsInsert       string
					SQLFieldsMergeInsert  string
					SQLFieldsMergeUpdate  string
					SQLPlaceholdersInsert string
					SQLPlaceholdersMerge  string

					StructFieldsSelectGet  string
					StructFieldsSelectList string
					StructFieldsUpdate     string
					StructFieldsInsert     string
					StructFieldsMerge      string

					Joins         string
					JoinVarsDecl  []string
					JoinVarsAssgn []string

					BeforeUpdate []string
					BeforeInsert []string
					BeforeMerge  []string

					HasRelationshipManyMany bool
					ManyManyFields          []util.Field
					Imports                 []string
				}

				sqlfieldsSelectGet    []string
				sqlfieldsSelectList   []string
				sqlfieldsUpdate       []string
				sqlfieldsInsert       []string
				sqlPlaceholdersInsert []string
				sqlPlaceholdersMerge  []string
				sqlFieldsMergeInsert  []string
				sqlFieldsMergeUpdate  []string

				structFieldsSelectGet  []string
				structFieldsSelectList []string
				structFieldsUpdate     []string
				structFieldsInsert     []string
				structFieldsMerge      []string
				sqlPlaceholderCount    = 1

				joins     []string
				joinCount int
			)

			if entity.Crud == nil {
				entity.Crud = &opts
			}

			if !entity.Crud.Create && !entity.Crud.Read && !entity.Crud.ReadList && !entity.Crud.Update && !entity.Crud.Delete {
				work.Done <- util.GeneratedCode{Generator: "GenerateCrud", Error: util.ErrorSkip}
			}

			if entity.PrimaryKey == "" {
				entity.PrimaryKey = util.PrimaryKeySerial
			}

			sqlFieldsMergeInsert = append(sqlFieldsMergeInsert, "id")
			structFieldsMerge = append(structFieldsMerge, "*entity.ID")
			sqlPlaceholdersMerge = append(sqlPlaceholdersMerge, fmt.Sprintf("$%d", sqlPlaceholderCount))

			if entity.PrimaryKey != util.PrimaryKeySerial {
				sqlPlaceholdersInsert = append(sqlPlaceholdersInsert, fmt.Sprintf("$%d", sqlPlaceholderCount))
				sqlfieldsInsert = append(sqlfieldsInsert, "id")
				structFieldsInsert = append(structFieldsInsert, "*entity.ID")

				sqlPlaceholderCount++
			}

			for _, field := range entity.Fields {
				if field.Relationship.Type == "" {
					sqlfieldsSelectGet = append(sqlfieldsSelectGet, fmt.Sprintf("t.%s", field.Schema.Field))
					sqlfieldsSelectList = append(sqlfieldsSelectList, fmt.Sprintf("t.%s", field.Schema.Field))
					structFieldsSelectGet = append(structFieldsSelectGet, fmt.Sprintf("entity.%s", field.Property.Name))
					structFieldsSelectList = append(structFieldsSelectList, fmt.Sprintf("entity.%s", field.Property.Name))

					sqlPlaceholdersInsert = append(sqlPlaceholdersInsert, fmt.Sprintf("$%d", sqlPlaceholderCount))
					sqlPlaceholdersMerge = append(sqlPlaceholdersMerge, fmt.Sprintf("$%d", sqlPlaceholderCount))

					sqlfieldsUpdate = append(sqlfieldsUpdate, fmt.Sprintf("%s = $%d", field.Schema.Field, sqlPlaceholderCount))
					sqlFieldsMergeUpdate = append(sqlFieldsMergeUpdate, fmt.Sprintf("%s = $%d", field.Schema.Field, sqlPlaceholderCount))
					sqlfieldsInsert = append(sqlfieldsInsert, fmt.Sprintf("%s", field.Schema.Field))
					sqlFieldsMergeInsert = append(sqlFieldsMergeInsert, fmt.Sprintf("%s", field.Schema.Field))

					structFieldsInsert = append(structFieldsInsert, fmt.Sprintf("*entity.%s", field.Property.Name))
					structFieldsUpdate = append(structFieldsUpdate, fmt.Sprintf("*entity.%s", field.Property.Name))
					structFieldsMerge = append(structFieldsMerge, fmt.Sprintf("*entity.%s", field.Property.Name))
					sqlPlaceholderCount++

					if field.Property.Name == "CreatedAt" {
						data.BeforeInsert = append(data.BeforeInsert, "*entity.CreatedAt = time.Now()")
						data.BeforeMerge = append(data.BeforeMerge, "*entity.CreatedAt = time.Now()")
					} else if field.Property.Name == "UpdatedAt" {
						data.BeforeInsert = append(data.BeforeInsert, "*entity.UpdatedAt = time.Now()")
						data.BeforeUpdate = append(data.BeforeUpdate, "*entity.UpdatedAt = time.Now()")
						data.BeforeMerge = append(data.BeforeMerge, "*entity.UpdatedAt = time.Now()")
					}
				} else {
					joins = append(joins,
						fmt.Sprintf("%s jt%d ON (t.%s = jt%d.%s)",
							field.Relationship.Target.Table,
							joinCount,
							field.Relationship.Target.ThisID,
							joinCount,
							field.Relationship.Target.ThatID))
					data.JoinVarsDecl = append(data.JoinVarsDecl, fmt.Sprintf("j%d %s", joinCount, strings.TrimPrefix(field.Property.Type, "[]")))
					data.JoinVarsAssgn = append(data.JoinVarsAssgn, fmt.Sprintf("*entity.%s = append(*entity.%s, j%d)", field.Property.Name, field.Property.Name, joinCount))
					sqlfieldsSelectGet = append(sqlfieldsSelectGet, fmt.Sprintf("jt%d.%s", joinCount, field.Relationship.Target.ThatID))
					structFieldsSelectGet = append(structFieldsSelectGet, fmt.Sprintf("&j%d, ", joinCount))
					joinCount++

					if field.Relationship.Type == util.RelationshipTypeManyMany {
						data.HasRelationshipManyMany = true
						data.ManyManyFields = append(data.ManyManyFields, field)
					}
				}
			}

			data.Entity = entity
			data.Package = strings.ToLower(entity.Name)
			data.SQLFieldsSelectGet = strings.Join(sqlfieldsSelectGet, ", ")
			data.SQLFieldsSelectList = strings.Join(sqlfieldsSelectList, ", ")
			data.SQLFieldsUpdate = strings.Join(sqlfieldsUpdate, ", ")
			data.SQLFieldsInsert = strings.Join(sqlfieldsInsert, ", ")
			data.SQLFieldsMergeInsert = strings.Join(sqlFieldsMergeInsert, ", ")
			data.SQLFieldsMergeUpdate = strings.Join(sqlFieldsMergeUpdate, ", ")
			data.SQLPlaceholdersInsert = strings.Join(sqlPlaceholdersInsert, ", ")
			data.SQLPlaceholdersMerge = strings.Join(sqlPlaceholdersMerge, ", ")

			data.StructFieldsSelectGet = strings.Join(structFieldsSelectGet, ", ")
			data.StructFieldsSelectList = strings.Join(structFieldsSelectList, ", ")
			data.StructFieldsUpdate = strings.Join(structFieldsUpdate, ", ")
			data.StructFieldsInsert = strings.Join(structFieldsInsert, ", ")
			data.StructFieldsMerge = strings.Join(structFieldsMerge, ", ")

			if entity.PrimaryKey == util.PrimaryKeyUUID {
				data.Imports = append(data.Imports, `"github.com/satori/go.uuid"`)
			}

			if joinCount > 0 {
				data.Joins = "INNER JOIN " + strings.Join(joins, " INNER JOIN ") + " "
			}

			structure, err := util.ExecuteTemplate("crud_structure.go.tmpl", struct {
				Entity  util.Entity
				Package string
			}{entity, data.Package})
			if err == nil {
				work.Done <- util.GeneratedCode{Generator: "GenerateCRUDModel", Code: structure, Filename: fmt.Sprintf("models/%s/%s.gocipe.go", data.Package, data.Package)}
			} else {
				work.Done <- util.GeneratedCode{Generator: "GenerateCRUDModel", Error: fmt.Errorf("failed to load execute template: %s", err)}
			}

			code, err := util.ExecuteTemplate("crud.go.tmpl", data)
			if err == nil {
				work.Done <- util.GeneratedCode{Generator: "GenerateCRUD", Code: code, Filename: fmt.Sprintf("models/%s/%s_crud.gocipe.go", data.Package, data.Package)}
			} else {
				work.Done <- util.GeneratedCode{Generator: "GenerateCRUD", Error: fmt.Errorf("failed to load execute template: %s", err)}
			}
			if entity.Crud.Hooks.PreSave || entity.Crud.Hooks.PostSave || entity.Crud.Hooks.PreRead || entity.Crud.Hooks.PostRead || entity.Crud.Hooks.PreList || entity.Crud.Hooks.PostList || entity.Crud.Hooks.PreDelete || entity.Crud.Hooks.PostDelete {
				hooks, e := util.ExecuteTemplate("crud_hooks.go.tmpl", struct {
					Hooks   util.CrudHooks
					Entity  util.Entity
					Package string
				}{entity.Crud.Hooks, entity, data.Package})

				if e == nil {
					work.Done <- util.GeneratedCode{Generator: "GenerateCRUDHooks", Code: hooks, Filename: fmt.Sprintf("models/%s/%s_crud_hooks.gocipe.go", data.Package, data.Package), NoOverwrite: true}
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
		work.Done <- util.GeneratedCode{Generator: "GenerateCRUD", Code: code, Filename: "models/filters.gocipe.go"}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateCRUD", Error: fmt.Errorf("failed to load execute template: %s", err)}
	}

	return err
}
