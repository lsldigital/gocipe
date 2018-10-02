package crud

import (
	"fmt"
	"strings"

	"github.com/fluxynet/gocipe/util"
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
func Generate(work util.GenerationWork, crud util.CrudOpts, entities map[string]util.Entity) {
	generateAny := false
	work.Waitgroup.Add(len(entities) * 2) //2 threads per entities. for models and models_hooks

	for _, entity := range entities {
		generateAny = generateAny || crud.Generate

		if !crud.Generate {
			util.DeleteIfExists(fmt.Sprintf("models/%s_repo.gocipe.go", strings.ToLower(entity.Name)))
			util.DeleteIfExists(fmt.Sprintf("models/%s_crud_hooks.gocipe.go", strings.ToLower(entity.Name)))
			work.Done <- util.GeneratedCode{Generator: fmt.Sprintf("GenerateCRUD[%s]", entity.Name), Error: util.ErrorSkip}
			work.Done <- util.GeneratedCode{Generator: fmt.Sprintf("GenerateCRUDHooks[%s]", entity.Name), Error: util.ErrorSkip}
			continue
		}

		go func(entity util.Entity) {
			var (
				code        entityCrud
				crudContent string
				err         error
			)

			if crud.Generate {
				code, err = generateCrud(entity, entities)
				if err == nil {
					crudContent, err = util.ExecuteTemplate("crud/crud.go.tmpl", code)
				}
			}

			if err == nil {
				fname := fmt.Sprintf("models/%s_repo.gocipe.go", strings.ToLower(entity.Name))
				work.Done <- util.GeneratedCode{Generator: fmt.Sprintf("GenerateCRUD[%s]", entity.Name), Code: crudContent, Filename: fname}
			} else {
				work.Done <- util.GeneratedCode{Generator: fmt.Sprintf("GenerateCRUD[%s]", entity.Name), Error: fmt.Errorf("failed to execute template: %s", err)}
			}

			hasHooks := entity.CrudHooks.PreSave ||
				entity.CrudHooks.PostSave ||
				entity.CrudHooks.PreRead ||
				entity.CrudHooks.PostRead ||
				entity.CrudHooks.PreList ||
				entity.CrudHooks.PostList ||
				entity.CrudHooks.PreDeleteSingle ||
				entity.CrudHooks.PostDeleteSingle ||
				entity.CrudHooks.PreDeleteMany ||
				entity.CrudHooks.PostDeleteMany

			if hasHooks {
				hooks, err := util.ExecuteTemplate("crud/hooks.go.tmpl", struct {
					Hooks  util.CrudHooks
					Entity util.Entity
				}{*entity.CrudHooks, entity})

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
			Crud     bool
			Entities map[string]util.Entity
		}{
			Crud:     crud.Generate,
			Entities: entities,
		})
		if err == nil {
			work.Done <- util.GeneratedCode{Generator: "GenerateCRUDModels", Code: models, Filename: "models/models.gocipe.go"}
		} else {
			work.Done <- util.GeneratedCode{Generator: "GenerateCRUDModels", Error: fmt.Errorf("failed to load execute template: %s", err)}
		}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateCRUDModels", Error: util.ErrorSkip}
	}

	work.Waitgroup.Add(1)
	models, err := util.ExecuteTemplate("crud/moderrors.go.tmpl", struct {
		Entities map[string]util.Entity
	}{entities})
	if err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateCRUDModelErrors", Code: models, Filename: "models/moderrors/errors.gocipe.go"}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateCRUDModelErrors", Error: fmt.Errorf("failed to load execute template: %s", err)}
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

	if err == nil {
		code.Insert, err = generateInsert(entities, entity)

		if entity.PrimaryKey == util.PrimaryKeyUUID {
			importUUID = true
		}
	}

	if err == nil {
		code.Get, err = generateGet(entities, entity)
	}

	if err == nil {
		code.List, err = generateList(entities, entity)
	}

	if err == nil {
		code.Update, err = generateUpdate(entities, entity)
	}

	if err == nil {
		code.DeleteMany, err = generateDeleteMany(entities, entity)
		if err == nil {
			code.DeleteSingle, err = generateDeleteSingle(entities, entity)
		}
	}

	if err == nil {
		code.Merge, err = generateMerge(entities, entity)
	}

	if err == nil {
		code.Save, err = generateSave(entities, entity)
	}

	if err == nil {
		for _, rel := range entity.Relationships {
			// No SaveRelated template generated:
			// RelationshipTypeManyManyInverse, RelationshipTypeOneMany, RelationshipTypeManyOne, RelationshipTypeOneOne
			switch rel.Type {
			case util.RelationshipTypeManyManyInverse:
				c, err := generateLoadRelatedManyMany(entities, entity, rel)
				if err != nil {
					return code, err
				}
				code.LoadRelated = append(code.LoadRelated, c)
			case util.RelationshipTypeManyManyOwner:
				c, err := generateLoadRelatedManyMany(entities, entity, rel)
				if err != nil {
					return code, err
				}
				code.LoadRelated = append(code.LoadRelated, c)

				c, err = generateSaveRelatedManyManyOwner(entities, entity, rel)
				if err != nil {
					return code, err
				}
				code.SaveRelated = append(code.SaveRelated, c)
			case util.RelationshipTypeManyMany:
				c, err := generateLoadRelatedManyMany(entities, entity, rel)
				if err != nil {
					return code, err
				}
				code.LoadRelated = append(code.LoadRelated, c)

				c, err = generateSaveRelatedManyMany(entities, entity, rel)
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
			case util.RelationshipTypeManyOne, util.RelationshipTypeOneOne:
				c, err := generateLoadRelatedManyOne(entities, entity, rel)
				if err != nil {
					return code, err
				}
				code.LoadRelated = append(code.LoadRelated, c)
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
	var post []string

	for _, rel := range entity.Relationships {
		if rel.Type == util.RelationshipTypeManyMany || rel.Type == util.RelationshipTypeManyManyOwner {
			post = append(post, fmt.Sprintf("repo.Save%s(ctx, tx, false, entity.ID)", util.RelFuncName(rel)))
		}
	}

	return util.ExecuteTemplate("crud/partials/delete_single.go.tmpl", struct {
		EntityName  string
		PrimaryKey  string
		Table       string
		HasPreHook  bool
		HasPostHook bool
		Post        []string
	}{
		EntityName:  entity.Name,
		PrimaryKey:  entity.PrimaryKey,
		Table:       entity.Table,
		HasPreHook:  entity.CrudHooks.PreDeleteSingle,
		HasPostHook: entity.CrudHooks.PostDeleteSingle,
		Post:        post,
	})
}

// generateDeleteMany produces code for database deletion of entity via filters
func generateDeleteMany(entities map[string]util.Entity, entity util.Entity) (string, error) {
	return util.ExecuteTemplate("crud/partials/delete_many.go.tmpl", struct {
		EntityName    string
		PrimaryKey    string
		Table         string
		HasPreHook    bool
		HasPostHook   bool
		Relationships []util.Relationship
	}{
		EntityName:    entity.Name,
		PrimaryKey:    entity.PrimaryKey,
		Table:         entity.Table,
		HasPreHook:    entity.CrudHooks.PreDeleteMany,
		HasPostHook:   entity.CrudHooks.PostDeleteMany,
		Relationships: entity.Relationships,
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

// generateSaveRelatedManyMany produces code for database saving of related entities
func generateSaveRelatedManyMany(entities map[string]util.Entity, entity util.Entity, rel util.Relationship) (string, error) {
	return util.ExecuteTemplate("crud/partials/saverelated_manymany.go.tmpl", struct {
		PropertyName   string
		PrimaryKey     string
		PropertyType   string
		EntityName     string
		Table          string
		Funcname       string
		ThisColumn     string
		ThatColumn     string
		ThatType       string
		ThatPrimaryKey string
	}{
		PropertyName:   rel.Name,
		PrimaryKey:     entity.PrimaryKey,
		PropertyType:   entities[rel.Entity].PrimaryKey,
		EntityName:     entity.Name,
		Table:          rel.JoinTable,
		Funcname:       util.RelFuncName(rel),
		ThisColumn:     rel.ThisID,
		ThatColumn:     rel.ThatID,
		ThatType:       "*" + entities[rel.Entity].Name,
		ThatPrimaryKey: entities[rel.Entity].PrimaryKey,
	})
}

// generateSaveRelatedManyManyOwner produces code for database saving of related entities
func generateSaveRelatedManyManyOwner(entities map[string]util.Entity, entity util.Entity, rel util.Relationship) (string, error) {
	return util.ExecuteTemplate("crud/partials/saverelated_manymanyowner.go.tmpl", struct {
		PropertyName   string
		PrimaryKey     string
		PropertyType   string
		EntityName     string
		Table          string
		Funcname       string
		ThisColumn     string
		ThatColumn     string
		ThatType       string
		ThatPrimaryKey string
	}{
		PropertyName:   rel.Name,
		PrimaryKey:     entity.PrimaryKey,
		PropertyType:   entities[rel.Entity].PrimaryKey,
		EntityName:     entity.Name,
		Table:          rel.JoinTable,
		Funcname:       util.RelFuncName(rel),
		ThisColumn:     rel.ThisID,
		ThatColumn:     rel.ThatID,
		ThatType:       "*" + entities[rel.Entity].Name,
		ThatPrimaryKey: entities[rel.Entity].PrimaryKey,
	})
}
