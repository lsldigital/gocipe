package crud

import (
	"fmt"
	"strings"

	"github.com/fluxynet/gocipe/util"
)

type entityCrud struct {
	Imports      []string
	Structure    string
	Get          string
	List         string
	DeleteSingle string
	DeleteMany   string
	Save         string
	Insert       string
	Update       string
	Merge        string
	SaveRelated  []string
	LoadRelated  []string
}

type relationship struct {
	Table        string
	ThisID       string
	ThatID       string
	PropertyName string
}

// Generate returns generated code to run an http server
func Generate(work util.GenerationWork, opts util.CrudOpts, entityList []util.Entity) {
	entities, err := preprocessEntities(entityList)

	if err != nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateCRUD", Error: err}
		return
	}

	work.Waitgroup.Add(len(entityList) * 2) //2 jobs to be waited upon for each thread - entity.go and entity_crud_hooks.go generation

	for _, entity := range entities {
		go func(entity util.Entity) {
			var (
				code entityCrud
				crud string
				err  error
			)

			if entity.Crud == nil {
				entity.Crud = &opts
			}

			if entity.PrimaryKey == "" {
				entity.PrimaryKey = util.PrimaryKeySerial
			}

			code, err = generateCrud(entity, entities)
			if err == nil {
				crud, err = util.ExecuteTemplate("crud/crud.go.tmpl", code)
			}

			if err == nil {
				work.Done <- util.GeneratedCode{Generator: "GenerateCRUD", Code: crud, Filename: fmt.Sprintf("models/%s_repo.gocipe.go", strings.ToLower(entity.Name))}
			} else {
				work.Done <- util.GeneratedCode{Generator: "GenerateCRUD", Error: fmt.Errorf("failed to execute template: %s", err)}
			}

			hasHooks := entity.Crud.Hooks.PreSave ||
				entity.Crud.Hooks.PostSave ||
				entity.Crud.Hooks.PreRead ||
				entity.Crud.Hooks.PostRead ||
				entity.Crud.Hooks.PreList ||
				entity.Crud.Hooks.PostList ||
				entity.Crud.Hooks.PreDeleteSingle ||
				entity.Crud.Hooks.PostDeleteSingle ||
				entity.Crud.Hooks.PreDeleteMany ||
				entity.Crud.Hooks.PostDeleteMany

			if hasHooks {
				hooks, err := util.ExecuteTemplate("crud/hooks.go.tmpl", struct {
					Hooks  util.CrudHooks
					Entity util.Entity
				}{entity.Crud.Hooks, entity})

				if err == nil {
					work.Done <- util.GeneratedCode{Generator: "GenerateCRUDHooks", Code: hooks, Filename: fmt.Sprintf("models/%s_crud_hooks.gocipe.go", strings.ToLower(entity.Name)), NoOverwrite: true}
				} else {
					work.Done <- util.GeneratedCode{Generator: "GenerateCRUDHooks", Error: err}
				}
			} else {
				work.Done <- util.GeneratedCode{Generator: "GenerateCRUDHooks", Error: util.ErrorSkip}
			}
		}(entity)
	}

	work.Waitgroup.Add(1)
	proto, err := generateProtobuf(entities)
	if err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateProto", Code: proto, Filename: "proto/entities.proto"}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateProto", Error: fmt.Errorf("failed to load execute template: %s", err)}
	}

	models, err := util.ExecuteTemplate("crud/models.go.tmpl", struct {
		Entities []util.Entity
	}{entityList})
	if err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateCRUDModels", Code: models, Filename: "models/models.gocipe.go"}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateCRUDModels", Error: fmt.Errorf("failed to load execute template: %s", err)}
	}
}

func generateCrud(entity util.Entity, entities map[string]util.Entity) (entityCrud, error) {
	var (
		code entityCrud
		err  error
	)

	if err == nil && entity.Crud.Create {
		code.Insert, err = generateInsert(entities, entity)
	}

	if err == nil && entity.Crud.Read {
		code.Get, err = generateGet(entities, entity)
	}

	if err == nil && entity.Crud.ReadList {
		code.List, err = generateList(entities, entity)
	}

	if err == nil && entity.Crud.Update {
		code.Update, err = generateUpdate(entities, entity)
	}

	if err == nil && entity.Crud.Delete {
		code.DeleteMany, err = generateDeleteMany(entities, entity)
		if err == nil {
			code.DeleteSingle, err = generateDeleteSingle(entities, entity)
		}
	}

	if err == nil && entity.Crud.Merge {
		code.Merge, err = generateMerge(entities, entity)
	}

	if err == nil && entity.Crud.Create && entity.Crud.Update && entity.Crud.Merge {
		code.Save, err = generateSave(entities, entity)
	}

	if err == nil {
		for _, rel := range entity.Relationships {
			switch rel.Type {
			case util.RelationshipTypeManyMany:
				c, err := generateLoadRelatedManyMany(entities, entity, rel)
				if err != nil {
					return code, err
				}
				code.LoadRelated = append(code.LoadRelated, c)
			case util.RelationshipTypeOneMany:
				c, err := generateLoadRelatedOneMany(entities, entity, rel)
				if err != nil {
					return code, err
				}
				code.LoadRelated = append(code.LoadRelated, c)
			case util.RelationshipTypeManyOne:
				if rel.Full {
					c, err := generateLoadRelatedManyOne(entities, entity, rel)
					if err != nil {
						return code, err
					}
					code.LoadRelated = append(code.LoadRelated, c)
				}
			}
		}
	}

	return code, err
}

// generateDeleteSingle produces code for database deletion of single entity (DELETE WHERE id)
func generateDeleteSingle(entities map[string]util.Entity, entity util.Entity) (string, error) {
	return util.ExecuteTemplate("crud/partials/delete_single.go.tmpl", struct {
		EntityName  string
		PrimaryKey  string
		Table       string
		HasPreHook  bool
		HasPostHook bool
	}{
		EntityName:  entity.Name,
		PrimaryKey:  entity.PrimaryKey,
		Table:       entity.Table,
		HasPreHook:  entity.Crud.Hooks.PreDeleteSingle,
		HasPostHook: entity.Crud.Hooks.PostDeleteSingle,
	})
}

// generateDeleteMany produces code for database deletion of entity via filters
func generateDeleteMany(entities map[string]util.Entity, entity util.Entity) (string, error) {
	return util.ExecuteTemplate("crud/partials/delete_many.go.tmpl", struct {
		EntityName  string
		PrimaryKey  string
		Table       string
		HasPreHook  bool
		HasPostHook bool
	}{
		EntityName:  entity.Name,
		PrimaryKey:  entity.PrimaryKey,
		Table:       entity.Table,
		HasPreHook:  entity.Crud.Hooks.PreDeleteMany,
		HasPostHook: entity.Crud.Hooks.PostDeleteMany,
	})
}

// generateSave produces code for database saving of entity
func generateSave(entities map[string]util.Entity, entity util.Entity) (string, error) {
	return util.ExecuteTemplate("crud/partials/save.go.tmpl", struct {
		EntityName string
		PrimaryKey string
	}{
		EntityName: entity.Name,
		PrimaryKey: entity.PrimaryKey,
	})
}

// generateSaveRelated produces code for database saving of related entities
func generateSaveRelated(entities map[string]util.Entity, entity util.Entity) (string, error) {
	return util.ExecuteTemplate("crud/partials/saverelated.go.tmpl", struct {
		PropertyName string
		PrimaryKey   string
		PropertyType string
		EntityName   string
		Table        string
		ThisID       string
		ThatID       string
	}{
		PropertyName: "",
		PrimaryKey:   "",
		PropertyType: "",
		EntityName:   "",
		Table:        "",
		ThisID:       "",
		ThatID:       "",
	})
}
