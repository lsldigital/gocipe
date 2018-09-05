package bread

import (
	"fmt"
	"strings"

	"github.com/fluxynet/gocipe/util"
)

// Generate returns generated code for a BREAD service - Browse, Read, Edit, Add & Delete
func Generate(work util.GenerationWork, entities map[string]util.Entity) error {
	for _, entity := range entities {
		if !entity.Bread.Generate {
			continue
		}

		var fileFields []string
		for _, field := range entity.Fields {
			switch field.EditWidget.Type {
			case util.WidgetTypeFile, util.WidgetTypeImage:
				fileFields = append(fileFields, field.Property.Name)
			}
		}

		hasFileFields := len(fileFields) > 0
		packagename := strings.ToLower(entity.Name)

		proto, err := util.ExecuteTemplate("bread/service_bread_entity.proto.tmpl", struct {
			Entity        util.Entity
			AppImportPath string
			HasFileFields bool
		}{entity, util.AppImportPath, hasFileFields})

		// generate separate bread services per entity
		work.Waitgroup.Add(1)
		if err == nil {
			work.Done <- util.GeneratedCode{Generator: "GenerateBreadProto", Code: proto, Filename: fmt.Sprintf("proto/service_bread_%s.proto", strings.ToLower(entity.Name))}
		} else {
			work.Done <- util.GeneratedCode{Generator: "GenerateBreadProto", Error: err}
		}

		code, err := util.ExecuteTemplate("bread/service_bread.entity.gocipe.go.tmpl", struct {
			Entity        util.Entity
			FileFields    string
			HasFileFields bool
		}{
			Entity:        entity,
			FileFields:    `"` + strings.Join(fileFields, `","`) + `"`,
			HasFileFields: hasFileFields,
		})

		work.Waitgroup.Add(1)
		if err == nil {
			work.Done <- util.GeneratedCode{Generator: "GenerateEntityBread|" + packagename, Code: code, Filename: "services/bread/" + packagename + "/bread.gocipe.go"}
		} else {
			work.Done <- util.GeneratedCode{Generator: "GenerateEntityBread|" + packagename, Error: fmt.Errorf("failed to execute template: %s", err)}
		}

		if hasFileFields {
			code, err = util.ExecuteTemplate("bread/bread_config_upload.gocipe.go.tmpl", struct {
				Entity        util.Entity
				FileFields    []string
				HasFileFields bool
			}{
				Entity:        entity,
				FileFields:    fileFields,
				HasFileFields: hasFileFields,
			})

			work.Waitgroup.Add(1)
			if err == nil {
				work.Done <- util.GeneratedCode{Generator: "GenerateBread Upload", Code: code, Filename: "services/bread/" + packagename + "/service_bread_config_upload.gocipe.go"}
			} else {
				work.Done <- util.GeneratedCode{Generator: "GenerateBread Upload", Error: fmt.Errorf("failed to execute template: %s", err)}
			}
		}

		if hasHook(entity) {
			hooks, err := util.ExecuteTemplate("bread/service_bread_hooks.go.tmpl", struct {
				Entity     util.Entity
				PreRead    bool
				PostRead   bool
				PreList    bool
				PostList   bool
				PreCreate  bool
				PostCreate bool
				PreUpdate  bool
				PostUpdate bool
				PreDelete  bool
				PostDelete bool
			}{
				Entity:     entity,
				PreRead:    entity.Bread.Hooks.PreRead,
				PostRead:   entity.Bread.Hooks.PostRead,
				PreList:    entity.Bread.Hooks.PreList,
				PostList:   entity.Bread.Hooks.PostList,
				PreCreate:  entity.Bread.Hooks.PreCreate,
				PostCreate: entity.Bread.Hooks.PostCreate,
				PreUpdate:  entity.Bread.Hooks.PreUpdate,
				PostUpdate: entity.Bread.Hooks.PostUpdate,
				PreDelete:  entity.Bread.Hooks.PreDelete,
				PostDelete: entity.Bread.Hooks.PostDelete,
			})

			work.Waitgroup.Add(1)
			if err == nil {
				work.Done <- util.GeneratedCode{
					Generator: "GenerateBreadHooks:" + entity.Name,
					Code:      hooks,
					Filename: fmt.Sprintf(
						"services/bread/%s/%s_hooks.gocipe.go",
						strings.ToLower(entity.Name),
						strings.ToLower(entity.Name),
					),
					NoOverwrite: true,
				}
			} else {
				work.Done <- util.GeneratedCode{Generator: "GenerateBreadHooks", Error: err}
			}
		}

	}

	// generate bread.proto
	proto, err := util.ExecuteTemplate("bread/service_bread.proto.tmpl", struct {
		AppImportPath string
	}{util.AppImportPath})

	work.Waitgroup.Add(1)
	if err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateBreadProto", Code: proto, Filename: "proto/bread.proto"}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateBreadProto", Error: err}
	}

	code, err := util.ExecuteTemplate("bread/bread.gocipe.go.tmpl", struct{}{})

	work.Waitgroup.Add(1)
	if err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateBread", Code: code, Filename: "services/bread/bread.gocipe.go"}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateBread", Error: fmt.Errorf("failed to execute template: %s", err)}
	}

	work.Waitgroup.Done()
	return nil
}

func hasHook(entity util.Entity) bool {
	switch true {
	case
		entity.Bread.Hooks.PreCreate,
		entity.Bread.Hooks.PostCreate,
		entity.Bread.Hooks.PreRead,
		entity.Bread.Hooks.PostRead,
		entity.Bread.Hooks.PreList,
		entity.Bread.Hooks.PostList,
		entity.Bread.Hooks.PreUpdate,
		entity.Bread.Hooks.PostUpdate,
		entity.Bread.Hooks.PreDelete,
		entity.Bread.Hooks.PostDelete:
		return true
	}
	return false
}
