package recipe

import (
	"fmt"
	"strings"

	"github.com/fluxynet/gocipe/util"
	"github.com/jinzhu/inflection"
)

// Preprocess does some preprocessing (checking, etc)
func Preprocess(recipe *util.Recipe) (map[string]util.Entity, error) {
	var (
		err error
	)
	fieldLabelWhiteList := []string{"name", "title", "description", "summary", "banner_type"}

	entities := make(map[string]util.Entity)
	for i, entity := range recipe.Entities {
		if entity.Name == "" {
			return nil, fmt.Errorf("entity #%d name cannot be blank", i)
		}

		entity.Table = inflection.Plural(strings.ToLower(entity.Name))

		for i, field := range entity.Fields {
			if field.Schema.Field == "" {
				field.Schema.Field = strings.ToLower(field.Property.Name)
			}

			switch field.Schema.Field {
			case "status", "id":
				return nil, fmt.Errorf("%s is a reserved fieldname", field.Schema.Field)
			}
			entity.Fields[i] = field
		}

		entity.Fields = append(entity.Fields, util.Field{
			Label:    "Status",
			Property: util.FieldProperty{Name: "Status", Type: "string"},
			Schema:   util.FieldSchema{Field: "status", Type: "CHAR(1)", Default: "'D'"},
			EditWidget: util.EditWidgetOpts{
				Type: util.WidgetTypeStatus,
				Options: []util.EditWidgetOption{
					util.EditWidgetOption{Text: "Draft", Value: "D"},
					util.EditWidgetOption{Text: "Saved", Value: "S"},
					util.EditWidgetOption{Text: "Published", Value: "P"},
					util.EditWidgetOption{Text: "Unpublished", Value: "U"},
				},
			},
			ListWidget: util.ListWidgetOpts{
				Type: util.WidgetTypeSelect,
			},
		})

		if entity.Slug != "" {
			var slugValid bool
			for _, field := range entity.Fields {
				fieldSchemaName := strings.ToLower(field.Schema.Field)
				propertyType := strings.ToLower(field.Property.Type)
				if entity.Slug == fieldSchemaName && propertyType == "string" {
					slugValid = true
					break
				}
			}
			if slugValid {
				entity.Fields = append(entity.Fields, util.Field{
					Label:    "Slug",
					Property: util.FieldProperty{Name: "Slug", Type: "string"},
					Schema:   util.FieldSchema{Field: "slug", Type: "VARCHAR(255)"},
					EditWidget: util.EditWidgetOpts{
						Hide: true,
					},
					ListWidget: util.ListWidgetOpts{
						Type: util.WidgetTypeTextField,
					},
				})
			} else {
				entity.Slug = ""
			}
		}

		if entity.DefaultSort == "" {
			entity.DefaultSort = `t."id" DESC`
		}

		if entity.CrudHooks == nil {
			entity.CrudHooks = &recipe.Crud.Hooks
		}

		if entity.Admin == nil {
			entity.Admin = &recipe.Admin
		}

		if entity.PrimaryKey == "" {
			entity.PrimaryKey = util.PrimaryKeySerial
		}

		if entity.Vuetify.Icon == "" {
			entity.Vuetify.Icon = "dashboard"
		}

		if entity.LabelField != "" {
			var fieldValid bool
			for _, field := range entity.Fields {
				fieldSchemaName := strings.ToLower(field.Schema.Field)
				propertyType := strings.ToLower(field.Property.Type)
				if entity.LabelField == fieldSchemaName && propertyType == "string" {
					fieldValid = true
					break
				}
			}
			if !fieldValid {
				entity.LabelField = ""
			}
		}

		if entity.LabelField == "" {
			var defaultLabelField string
			var firstStringField string
			for _, field := range entity.Fields {
				if defaultLabelField != "" {
					break
				}
				fieldSchemaName := strings.ToLower(field.Schema.Field)
				propertyType := strings.ToLower(field.Property.Type)
				for _, fieldName := range fieldLabelWhiteList {
					if fieldName == fieldSchemaName && propertyType == "string" {
						defaultLabelField = field.Schema.Field
						break
					}
				}
				if firstStringField == "" && propertyType == "string" {
					firstStringField = field.Schema.Field
				}
			}
			if defaultLabelField == "" {
				defaultLabelField = firstStringField
			}
			entity.LabelField = defaultLabelField
		}

		entities[entity.Name] = entity
	}

	for _, entity := range entities {
		for r := range entity.Relationships {
			var isMany bool
			rel := &entities[entity.Name].Relationships[r]

			if _, ok := entities[rel.Entity]; rel.Entity == "" || !ok {
				return nil, fmt.Errorf("relationship %s invalid in entity %s", rel.Entity, entity.Name)
			}

			switch rel.Type {
			default:
				return nil, fmt.Errorf("invalid relationship type %s for entity %s", rel.Type, entity.Name)
			case util.RelationshipTypeOneOne:
				if rel.ThatID == "" {
					rel.ThatID = "id"
				}

				if rel.ThisID == "" {
					rel.ThisID = "id"
				}
			case util.RelationshipTypeOneMany:
				rel.ThisID = "id"
				isMany = true
				rel.JoinTable = entities[rel.Entity].Table
				if rel.ThatID == "" {
					rel.ThatID = strings.ToLower(entity.Name) + "_id"
				}
			case util.RelationshipTypeManyOne:
				rel.ThatID = "id"
				if rel.ThisID == "" {
					rel.ThisID = strings.ToLower(rel.Entity) + "_id"
				}
			case util.RelationshipTypeManyMany, util.RelationshipTypeManyManyOwner, util.RelationshipTypeManyManyInverse:
				isMany = true
				if rel.ThatID == "" {
					rel.ThatID = strings.ToLower(entity.Name) + "_id"
				}

				if rel.ThisID == "" {
					rel.ThisID = strings.ToLower(rel.Entity) + "_id"
				}

				if rel.JoinTable == "" {
					if strings.Compare(entity.Table, entities[rel.Entity].Table) == -1 {
						rel.JoinTable = entity.Table + "_" + entities[rel.Entity].Table
					} else {
						rel.JoinTable = entities[rel.Entity].Table + "_" + entity.Table
					}
				}
			}

			// if rel.Name == "" {
			if isMany {
				rel.Name = inflection.Plural(strings.Title(strings.ToLower(rel.Entity)))
			} else {
				rel.Name = strings.Title(rel.Entity)
			}
			// }
		}
		entities[entity.Name] = entity
	}

	return entities, err
}
