package admin

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/fluxynet/gocipe/util"
)

type fileField struct {
	Entity string
	Field  string
}

// Generate returns generated code for a Admin service
func Generate(work util.GenerationWork, entities map[string]util.Entity) error {
	var (
		fileFields   []fileField
		generateAuth bool
	)
	entitiesActions := []string{"Create", "Edit", "View", "List", "Delete", "Lookup"}
	entitiesFileFields := make(map[string][]fileField)
	entitiesLabelField := make(map[string]string)
	for key, entity := range entities {
		if !entity.Admin.Generate {
			continue
		}

		if !generateAuth && entity.Admin.Auth.Generate {
			generateAuth = true
		}

		for _, field := range entity.Fields {
			switch field.EditWidget.Type {
			case util.WidgetTypeFile, util.WidgetTypeImage:
				fileFields = append(fileFields, fileField{Entity: entity.Name, Field: field.Property.Name})
				entitiesFileFields[key] = append(entitiesFileFields[key], fileField{Entity: entity.Name, Field: field.Property.Name})
			}
		}

		entitiesLabelField[key] = entity.LabelField

		if hasHook(entity) {
			hooks, err := util.ExecuteTemplate("admin/service_admin_hooks.go.tmpl", struct {
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
				ImportPath string
			}{
				Entity:     entity,
				PreRead:    entity.Admin.Hooks.PreRead,
				PostRead:   entity.Admin.Hooks.PostRead,
				PreList:    entity.Admin.Hooks.PreList,
				PostList:   entity.Admin.Hooks.PostList,
				PreCreate:  entity.Admin.Hooks.PreCreate,
				PostCreate: entity.Admin.Hooks.PostCreate,
				PreUpdate:  entity.Admin.Hooks.PreUpdate,
				PostUpdate: entity.Admin.Hooks.PostUpdate,
				PreDelete:  entity.Admin.Hooks.PreDelete,
				PostDelete: entity.Admin.Hooks.PostDelete,
				ImportPath: util.AppImportPath,
			})

			work.Waitgroup.Add(1)
			if err == nil {
				work.Done <- util.GeneratedCode{
					Generator: "GenerateAdminHooks:" + entity.Name,
					Code:      hooks,
					Filename: fmt.Sprintf(
						"services/admin/%s_hooks.gocipe.go",
						strings.ToLower(entity.Name),
					),
					NoOverwrite: true,
				}
			} else {
				work.Done <- util.GeneratedCode{Generator: "GenerateAdminHooks", Error: err}
			}
		}
	}

	hasFileFields := len(fileFields) > 0
	if hasFileFields {
		code, err := util.ExecuteTemplate("admin/admin_config_upload.go.tmpl", struct {
			FileFields []fileField
		}{
			FileFields: fileFields,
		})

		work.Waitgroup.Add(1)
		if err == nil {
			work.Done <- util.GeneratedCode{Generator: "GenerateAdmin Upload", Code: code, Filename: "services/admin/service_admin_config_upload.gocipe.go", NoOverwrite: true}
		} else {
			work.Done <- util.GeneratedCode{Generator: "GenerateAdmin Upload", Error: fmt.Errorf("failed to execute template: %s", err)}
		}
	}

	// generate admin_helpers.gocipe.go
	helpers, err := util.ExecuteTemplate("admin/admin_helpers.go.tmpl", struct {
		FileFields []fileField
		ImportPath string
	}{
		FileFields: fileFields,
		ImportPath: util.AppImportPath,
	})

	work.Waitgroup.Add(1)
	if err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateAdmin Helpers", Code: helpers, Filename: "services/admin/admin_helpers.gocipe.go"}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateAdmin Helpers", Error: fmt.Errorf("failed to execute template: %s", err)}
	}

	// generate admin.proto
	proto, err := util.ExecuteTemplate("admin/service_admin.proto.tmpl", struct {
		AppImportPath string
		Entities      map[string]util.Entity
	}{util.AppImportPath, entities})

	work.Waitgroup.Add(1)
	if err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateAdminProto", Code: proto, Filename: "proto/service_admin.proto"}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateAdminProto", Error: err}
	}

	// generate admin_permissions.go
	genUTF8List := func() []string {
		var (
			strList []string
			r       rune
		)
		for i := 0; i < 500; i++ {
			r = 'A' + rune(i)
			if unicode.IsPrint(r) &&
				!unicode.IsSymbol(r) &&
				!unicode.IsSpace(r) &&
				!unicode.IsControl(r) &&
				r != '\\' {
				strList = append(strList, string(r))
			}
		}
		return strList
	}
	permissions, err := util.ExecuteTemplate("admin/admin_permissions.go.tmpl", struct {
		ImportPath      string
		Entities        map[string]util.Entity
		EntitiesActions []string
		UTF8List        []string
	}{util.AppImportPath, entities, entitiesActions, genUTF8List()})

	work.Waitgroup.Add(1)
	if err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateAdminProto", Code: permissions, Filename: "services/admin/service_admin_permissions.go"}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateAdminProto", Error: err}
	}

	code, err := util.ExecuteTemplate("admin/service_admin.go.tmpl", struct {
		Entities           map[string]util.Entity
		EntitiesFileFields map[string][]fileField
		EntitiesLabelField map[string]string
		GenerateAuth       bool
		ImportPath         string
	}{entities, entitiesFileFields, entitiesLabelField, generateAuth, util.AppImportPath})
	work.Waitgroup.Add(1)
	if err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateAdmin", Code: code, Filename: "services/admin/service_admin.gocipe.go"}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateAdmin", Error: fmt.Errorf("failed to execute template: %s", err)}
	}

	work.Waitgroup.Done()
	return nil
}

func hasHook(entity util.Entity) bool {
	switch true {
	case
		entity.Admin.Hooks.PreCreate,
		entity.Admin.Hooks.PostCreate,
		entity.Admin.Hooks.PreRead,
		entity.Admin.Hooks.PostRead,
		entity.Admin.Hooks.PreList,
		entity.Admin.Hooks.PostList,
		entity.Admin.Hooks.PreUpdate,
		entity.Admin.Hooks.PostUpdate,
		entity.Admin.Hooks.PreDelete,
		entity.Admin.Hooks.PostDelete:
		return true
	}
	return false
}
