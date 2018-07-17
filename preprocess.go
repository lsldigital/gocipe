package main

import (
	"fmt"
	"strings"

	"github.com/fluxynet/gocipe/util"
	"github.com/jinzhu/inflection"
)

func preprocessRecipe(recipe *util.Recipe) (map[string]util.Entity, error) {
	var (
		err error
	)

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
			entity.Fields[i] = field
		}

		if entity.Crud == nil {
			entity.Crud = &recipe.Crud
		}

		if entity.Bread == nil {
			entity.Bread = &recipe.Bread
		}

		if entity.PrimaryKey == "" {
			entity.PrimaryKey = util.PrimaryKeySerial
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
				if rel.ThatID == "" {
					rel.ThatID = strings.ToLower(entity.Name) + "_id"
				}
			case util.RelationshipTypeManyOne:
				rel.ThatID = "id"
				if rel.ThisID == "" {
					rel.ThisID = strings.ToLower(rel.Entity) + "_id"
				}
			case util.RelationshipTypeManyMany:
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

			if rel.Name == "" {
				if isMany {
					rel.Name = inflection.Plural(strings.Title(strings.ToLower(rel.Entity)))
				} else {
					rel.Name = strings.Title(rel.Entity)
				}
			}
		}
		entities[entity.Name] = entity
	}

	return entities, err
}
