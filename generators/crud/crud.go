package crud

import (
	"fmt"
	"strings"

	"github.com/fluxynet/gocipe/util"
	"github.com/jinzhu/inflection"
)

// // RepoCodes that must be analysed for interface generation
// var (
// 	RepoCodes []RepoCode
// 	repoChan  = make(chan RepoCode)
// )

// // RepoCode represents generated code for a crud repo, which will be used to extract an interface
// type RepoCode struct {
// 	SourceFile string
// 	TargetFile string
// }

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
func Generate(work util.GenerationWork, opts util.CrudOpts, entities map[string]util.Entity) {
	generateAny := false
	work.Waitgroup.Add(len(entities) * 2) //2 threads per entities. for models and models_hooks

	for _, entity := range entities {
		generateForEntity := entity.Crud.Create || entity.Crud.Read || entity.Crud.ReadList ||
			entity.Crud.Update || entity.Crud.Delete || entity.Crud.Merge

		generateAny = generateAny || generateForEntity

		if !generateForEntity {
			util.DeleteIfExists(fmt.Sprintf("models/%s_repo.gocipe.go", strings.ToLower(entity.Name)))
			util.DeleteIfExists(fmt.Sprintf("models/%s_crud_hooks.gocipe.go", strings.ToLower(entity.Name)))
			work.Done <- util.GeneratedCode{Generator: fmt.Sprintf("GenerateCRUD[%s]", entity.Name), Error: util.ErrorSkip}
			work.Done <- util.GeneratedCode{Generator: fmt.Sprintf("GenerateCRUDHooks[%s]", entity.Name), Error: util.ErrorSkip}
			continue
		}

		go func(entity util.Entity) {
			var (
				code entityCrud
				crud string
				err  error
			)

			code, err = generateCrud(entity, entities)
			if err == nil {
				crud, err = util.ExecuteTemplate("crud/crud.go.tmpl", code)
			}

			if err == nil {
				fname := fmt.Sprintf("models/%s_repo.gocipe.go", strings.ToLower(entity.Name))
				work.Done <- util.GeneratedCode{Generator: fmt.Sprintf("GenerateCRUD[%s]", entity.Name), Code: crud, Filename: fname}
			} else {
				work.Done <- util.GeneratedCode{Generator: fmt.Sprintf("GenerateCRUD[%s]", entity.Name), Error: fmt.Errorf("failed to execute template: %s", err)}
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
					work.Done <- util.GeneratedCode{Generator: fmt.Sprintf("GenerateCRUDHooks[%s]", entity.Name), Code: hooks, Filename: fmt.Sprintf("models/%s_crud_hooks.gocipe.go", strings.ToLower(entity.Name)), NoOverwrite: true}
				} else {
					work.Done <- util.GeneratedCode{Generator: fmt.Sprintf("GenerateCRUDHooks[%s]", entity.Name), Error: err}
				}
			} else {
				work.Done <- util.GeneratedCode{Generator: fmt.Sprintf("GenerateCRUDHooks[%s]", entity.Name), Error: util.ErrorSkip}
			}
		}(entity)
	}

	work.Waitgroup.Add(1)
	if generateAny {
		models, err := util.ExecuteTemplate("crud/models.go.tmpl", struct {
			Entities map[string]util.Entity
		}{entities})
		if err == nil {
			work.Done <- util.GeneratedCode{Generator: "GenerateCRUDModels", Code: models, Filename: "models/models.gocipe.go"}
		} else {
			work.Done <- util.GeneratedCode{Generator: "GenerateCRUDModels", Error: fmt.Errorf("failed to load execute template: %s", err)}
		}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateCRUDModels", Error: util.ErrorSkip}
	}

	work.Waitgroup.Add(1)
	if generateAny {
		models, err := util.ExecuteTemplate("crud/moderrors.go.tmpl", struct {
			Entities map[string]util.Entity
		}{entities})
		if err == nil {
			work.Done <- util.GeneratedCode{Generator: "GenerateCRUDModelErrors", Code: models, Filename: "models/moderrors/errors.gocipe.go"}
		} else {
			work.Done <- util.GeneratedCode{Generator: "GenerateCRUDModelErrors", Error: fmt.Errorf("failed to load execute template: %s", err)}
		}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateCRUDModelErrors", Error: util.ErrorSkip}
	}

	proto, err := generateProtobuf(entities)
	if err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateProto", Code: proto, Filename: "proto/models.proto"}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateProto", Error: fmt.Errorf("failed to load execute template: %s", err)}
	}
}

func generateCrud(entity util.Entity, entities map[string]util.Entity) (entityCrud, error) {
	var (
		code       entityCrud
		importUUID bool
		err        error
	)

	if err == nil && entity.Crud.Create {
		code.Insert, err = generateInsert(entities, entity)

		if entity.PrimaryKey == util.PrimaryKeyUUID {
			importUUID = true
		}
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

				c, err = generateSaveRelated(entities, entity, rel)
				if err != nil {
					return code, err
				}
				code.SaveRelated = append(code.SaveRelated, c)
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

	if importUUID {
		code.Imports = append(code.Imports, `uuid "github.com/satori/go.uuid"`)
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
func generateSaveRelated(entities map[string]util.Entity, entity util.Entity, rel util.Relationship) (string, error) {
	var table string

	if strings.Compare(entity.Table, entities[rel.Name].Table) > 0 {
		table = inflection.Plural(entities[rel.Name].Table) + "_" + inflection.Plural(entity.Table)
	} else {
		table = inflection.Plural(inflection.Plural(entity.Table) + "_" + entities[rel.Name].Table)
	}

	return util.ExecuteTemplate("crud/partials/saverelated.go.tmpl", struct {
		PropertyName string
		PrimaryKey   string
		PropertyType string
		EntityName   string
		Table        string
		Funcname     string
	}{
		PropertyName: rel.Name,
		PrimaryKey:   entity.PrimaryKey,
		PropertyType: entities[rel.Name].PrimaryKey,
		EntityName:   entity.Name,
		Table:        table,
		Funcname:     util.RelFuncName(rel),
	})
}
