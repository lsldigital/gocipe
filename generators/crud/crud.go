package crud

import (
	"fmt"
	"strings"

	"github.com/fluxynet/gocipe/util"
)

type entityCrud struct {
	Package      string
	Imports      []string
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

	work.Waitgroup.Add(len(entityList) * 3) //3 jobs to be waited upon for each thread - entity.go,  entity_crud.go and entity_crud_hooks.go generation

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
				work.Done <- util.GeneratedCode{Generator: "GenerateCRUD", Code: crud, Filename: fmt.Sprintf("models/%s/%s_crud.gocipe.go", code.Package, code.Package)}
			} else {
				work.Done <- util.GeneratedCode{Generator: "GenerateCRUD", Error: fmt.Errorf("failed to execute template: %s", err)}
			}

			structure, err := generateStructure(entities, entity)
			if err == nil {
				work.Done <- util.GeneratedCode{Generator: "GenerateCRUDModel", Code: structure, Filename: fmt.Sprintf("models/%s/%s.gocipe.go", code.Package, code.Package)}
			} else {
				work.Done <- util.GeneratedCode{Generator: "GenerateCRUDModel", Error: fmt.Errorf("failed to load execute template: %s", err)}
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
					Hooks   util.CrudHooks
					Entity  util.Entity
					Package string
				}{entity.Crud.Hooks, entity, code.Package})

				if err == nil {
					work.Done <- util.GeneratedCode{Generator: "GenerateCRUDHooks", Code: hooks, Filename: fmt.Sprintf("models/%s/%s_crud_hooks.gocipe.go", code.Package, code.Package), NoOverwrite: true}
				} else {
					work.Done <- util.GeneratedCode{Generator: "GenerateCRUDHooks", Error: err}
				}
			} else {
				work.Done <- util.GeneratedCode{Generator: "GenerateCRUDHooks", Error: util.ErrorSkip}
			}
		}(entity)
	}

	if filters, err := util.ExecuteTemplate("crud/filters.go.tmpl", struct{}{}); err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateCRUDFilters", Code: filters, Filename: "models/filters.gocipe.go"}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateCRUDFilters", Error: fmt.Errorf("failed to load execute template: %s", err)}
	}
}

func generateCrud(entity util.Entity, entities map[string]util.Entity) (entityCrud, error) {
	var (
		code entityCrud
		err  error
	)

	code.Package = strings.ToLower(entity.Name)

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

	return code, err
}

func generateStructure(entities map[string]util.Entity, entity util.Entity) (string, error) {
	var fields []struct {
		Name       string
		Type       string
		Serialized string
	}

	var initialization []struct {
		Name string
		Type string
	}

	for _, field := range entity.Fields {
		fields = append(fields, struct {
			Name       string
			Type       string
			Serialized string
		}{field.Property.Name, "*" + field.Property.Type, field.Serialized})

		initialization = append(initialization, struct {
			Name string
			Type string
		}{field.Property.Name, fmt.Sprintf("new(%s)", field.Property.Type)})
	}

	relation := func(rel util.Relationship, many, full bool) {
		var t string
		if full {
			t = fmt.Sprintf("%s.%s", strings.ToLower(rel.Entity), rel.Entity)
		} else {
			t, _ = util.GetPrimaryKeyDataType(entities[rel.Entity].PrimaryKey)
		}

		if many {
			t = "[]" + t
		} else {
			t = "*" + t
		}

		fields = append(fields, struct {
			Name       string
			Type       string
			Serialized string
		}{rel.Name, t, rel.Serialized})

	}

	for _, rel := range entity.Relationships {
		switch rel.Type {
		case util.RelationshipTypeOneOne:
			fallthrough
		case util.RelationshipTypeManyOne:
			relation(rel, false, rel.Full)
		case util.RelationshipTypeOneMany:
			fallthrough
		case util.RelationshipTypeManyMany:
			relation(rel, true, rel.Full)
		}
	}

	return util.ExecuteTemplate("crud/structure.go.tmpl", struct {
		Package     string
		Name        string
		Description string
		PrimaryKey  string
		Fields      []struct {
			Name       string
			Type       string
			Serialized string
		}
		Initialization []struct {
			Name string
			Type string
		}
	}{
		Package:        strings.ToLower(entity.Name),
		Name:           entity.Name,
		Description:    entity.Description,
		PrimaryKey:     entity.PrimaryKey,
		Fields:         fields,
		Initialization: initialization,
	})
}

// generateGet produces code for database retrieval of single entity (SELECT WHERE id)
func generateGet(entities map[string]util.Entity, entity util.Entity) (string, error) {
	var sqlfields, structfields []string //, after []string

	sqlfields = append(sqlfields, fmt.Sprintf("%s", "id"))
	structfields = append(structfields, fmt.Sprintf("entity.%s", "ID"))

	for _, field := range entity.Fields {
		sqlfields = append(sqlfields, fmt.Sprintf("%s", field.Schema.Field))
		structfields = append(structfields, fmt.Sprintf("entity.%s", field.Property.Name))
	}

	// for _, rel := range entity.Relationships {
	// 	switch rel.Type {
	// 	case util.RelationshipTypeOneOne:
	// 		if rel.Full {
	// 			after = append(after, fmt.Sprintf("entity.%s, err = %s.Get(ctx, )"))
	// 			after = append(after, fmt.Sprintf("entity.%s = %s.Get(ctx, )"))
	// 			after = append(after, fmt.Sprintf("entity.%s = %s.Get(ctx, )"))
	// 		} else {
	// 			sqlfields = append(sqlfields, fmt.Sprintf("%s", rel.ThisID))
	// 			structfields = append(structfields, fmt.Sprintf("entity.%s", rel.Name))
	// 		}

	// 	case util.RelationshipTypeOneMany:
	// 		if rel.Full {
	// 			after = append(after, fmt.Sprintf("entity.%s = %s.Get(ctx, )"))
	// 		} else {
	// 			sqlfields = append(sqlfields, fmt.Sprintf("%s", rel.ThisID))
	// 			structfields = append(structfields, fmt.Sprintf("entity.%s", rel.Name))
	// 		}

	// 	case util.RelationshipTypeManyOne:
	// 	case util.RelationshipTypeManyMany:
	// 	}
	// }

	return util.ExecuteTemplate("crud/partials/get.go.tmpl", struct {
		EntityName   string
		SQLFields    string
		Table        string
		StructFields string
		PrimaryKey   string
		HasPreHook   bool
		HasPostHook  bool
	}{
		EntityName:   entity.Name,
		Table:        entity.Table,
		SQLFields:    strings.Join(sqlfields, ", "),
		StructFields: strings.Join(structfields, ", "),
		PrimaryKey:   entity.PrimaryKey,
		HasPreHook:   entity.Crud.Hooks.PreRead,
		HasPostHook:  entity.Crud.Hooks.PostRead,
	})
}

// generateList produces code for database retrieval of list of entities with optional filters
func generateList(entities map[string]util.Entity, entity util.Entity) (string, error) {
	var sqlfields, structfields []string

	sqlfields = append(sqlfields, fmt.Sprintf("%s", "id"))
	structfields = append(structfields, fmt.Sprintf("entity.%s", "ID"))

	for _, field := range entity.Fields {
		sqlfields = append(sqlfields, fmt.Sprintf("%s", field.Schema.Field))
		structfields = append(structfields, fmt.Sprintf("entity.%s", field.Property.Name))
	}

	return util.ExecuteTemplate("crud/partials/list.go.tmpl", struct {
		EntityName   string
		SQLFields    string
		StructFields string
		Table        string
		HasPreHook   bool
		HasPostHook  bool
	}{
		EntityName:   entity.Name,
		Table:        entity.Table,
		SQLFields:    strings.Join(sqlfields, ", "),
		StructFields: strings.Join(structfields, ", "),
		HasPreHook:   entity.Crud.Hooks.PreList,
		HasPostHook:  entity.Crud.Hooks.PostList,
	})
}

// generateDeleteSingle produces code for database deletion of single entity (DELETE WHERE id)
func generateDeleteSingle(entities map[string]util.Entity, entity util.Entity) (string, error) {
	return util.ExecuteTemplate("crud/partials/delete_single.go.tmpl", struct {
		EntityName  string
		Table       string
		HasPreHook  bool
		HasPostHook bool
	}{
		EntityName:  entity.Name,
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

// generateInsert produces code for database insertion of entity (INSERT INTO)
func generateInsert(entities map[string]util.Entity, entity util.Entity) (string, error) {
	var (
		before, sqlPlaceholders, sqlfields, structFields []string
		count                                            int
	)

	if entity.PrimaryKey != util.PrimaryKeySerial {
		sqlPlaceholders = append(sqlPlaceholders, fmt.Sprintf("$%d", count))
		sqlfields = append(sqlfields, "id")
		structFields = append(structFields, "*entity.ID")

		count++
	}

	for _, field := range entity.Fields {
		sqlPlaceholders = append(sqlPlaceholders, fmt.Sprintf("$%d", count))
		sqlfields = append(sqlfields, fmt.Sprintf("%s", field.Schema.Field))
		structFields = append(structFields, fmt.Sprintf("*entity.%s", field.Property.Name))

		if field.Property.Name == "CreatedAt" {
			before = append(before, "*entity.CreatedAt = time.Now()")
		} else if field.Property.Name == "UpdatedAt" {
			before = append(before, "*entity.UpdatedAt = time.Now()")
		}
	}

	return util.ExecuteTemplate("crud/partials/insert.go.tmpl", struct {
		Before          []string
		EntityName      string
		PrimaryKey      string
		SQLFields       string
		SQLPlaceholders string
		StructFields    string
		Table           string
		HasPreHook      bool
		HasPostHook     bool
		Relationships   []relationship
	}{
		Before:          before,
		EntityName:      entity.Name,
		PrimaryKey:      entity.PrimaryKey,
		SQLFields:       strings.Join(sqlfields, ", "),
		SQLPlaceholders: strings.Join(sqlPlaceholders, ", "),
		StructFields:    strings.Join(structFields, ", "),
		Table:           entity.Table,
		HasPostHook:     entity.Crud.Hooks.PreSave,
		HasPreHook:      entity.Crud.Hooks.PostSave,
		Relationships:   nil,
	})
}

// generateUpdate produces code for database update of entity (UPDATE)
func generateUpdate(entities map[string]util.Entity, entity util.Entity) (string, error) {
	var (
		before, sqlfields, structfields []string
		count                           = 1
	)

	for _, field := range entity.Fields {
		sqlfields = append(sqlfields, fmt.Sprintf("%s = $%d", field.Schema.Field, count))
		structfields = append(structfields, fmt.Sprintf("*entity.%s", field.Property.Name))
		count++

		if field.Property.Name == "CreatedAt" {
			before = append(before, "*entity.CreatedAt = time.Now()")
		} else if field.Property.Name == "UpdatedAt" {
			before = append(before, "*entity.UpdatedAt = time.Now()")
		}
	}

	return util.ExecuteTemplate("crud/partials/update.go.tmpl", struct {
		Before        []string
		EntityName    string
		HasPostHook   bool
		HasPreHook    bool
		SQLFields     string
		StructFields  string
		Table         string
		Relationships []relationship
	}{
		EntityName:    entity.Name,
		Table:         entity.Table,
		Before:        before,
		SQLFields:     strings.Join(sqlfields, ", "),
		StructFields:  strings.Join(structfields, ", "),
		HasPreHook:    entity.Crud.Hooks.PreSave,
		HasPostHook:   entity.Crud.Hooks.PostSave,
		Relationships: nil,
	})
}

// generateMerge produces code for database merge of entity (INSERT/ON CONFLICT UPDATE)
func generateMerge(entities map[string]util.Entity, entity util.Entity) (string, error) {
	var (
		before          []string
		sqlfieldsInsert []string
		sqlfieldsUpdate []string
		sqlPlaceholders []string
		structFields    []string
		count           = 0
	)

	sqlfieldsInsert = append(sqlfieldsInsert, "id")
	structFields = append(structFields, "*entity.ID")
	sqlPlaceholders = append(sqlPlaceholders, fmt.Sprintf("$%d", count))

	for _, field := range entity.Fields {
		sqlPlaceholders = append(sqlPlaceholders, fmt.Sprintf("$%d", count))
		sqlfieldsUpdate = append(sqlfieldsUpdate, fmt.Sprintf("%s = $%d", field.Schema.Field, count))
		sqlfieldsInsert = append(sqlfieldsInsert, fmt.Sprintf("%s", field.Schema.Field))
		structFields = append(structFields, fmt.Sprintf("*entity.%s", field.Property.Name))

		if field.Property.Name == "CreatedAt" {
			before = append(before, "*entity.CreatedAt = time.Now()")
		} else if field.Property.Name == "UpdatedAt" {
			before = append(before, "*entity.UpdatedAt = time.Now()")
		}
		count++
	}

	return util.ExecuteTemplate("crud/partials/merge.go.tmpl", struct {
		EntityName      string
		Before          []string
		Table           string
		SQLFieldsInsert string
		SQLPlaceholders string
		SQLFieldsUpdate string
		StructFields    string
		HasPreHook      bool
		HasPostHook     bool
	}{
		EntityName:      entity.Name,
		Before:          before,
		Table:           entity.Table,
		SQLFieldsInsert: strings.Join(sqlfieldsInsert, ", "),
		SQLPlaceholders: strings.Join(sqlPlaceholders, ", "),
		SQLFieldsUpdate: strings.Join(sqlfieldsUpdate, ", "),
		StructFields:    strings.Join(structFields, ", "),
		HasPreHook:      entity.Crud.Hooks.PreSave,
		HasPostHook:     entity.Crud.Hooks.PostSave,
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