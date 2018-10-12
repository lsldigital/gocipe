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
func Generate(work util.GenerationWork, r *util.Recipe) error {
	var (
		fileFields   []fileField
		generateAuth bool
	)
	entitiesActions := []string{"Create", "Edit", "View", "List", "Delete", "Lookup"}
	for _, entity := range r.Entities {
		if !entity.Admin.Generate {
			continue
		}

		if !generateAuth && entity.Admin.Auth.Generate {
			generateAuth = true
		}

		if hasHook(entity) {
			hooks, err := util.ExecuteTemplate("admin/service_admin_hooks.go.tmpl", struct {
				Entity     util.Entity
				ImportPath string
			}{
				Entity:     entity,
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
			Entities []util.Entity
		}{r.Entities})

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
		Entities      []util.Entity
	}{util.AppImportPath, r.Entities})

	work.Waitgroup.Add(1)
	if err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateAdminProto", Code: proto, Filename: "proto/service_admin.proto"}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateAdminProto", Error: err}
	}

	// generate admin_permissions.go
	permissions, err := util.ExecuteTemplate("admin/admin_permissions.go.tmpl", struct {
		ImportPath      string
		Entities        []util.Entity
		EntitiesActions []string
		UTF8List        []string
	}{util.AppImportPath, r.Entities, entitiesActions, genUTF8List()})

	work.Waitgroup.Add(1)
	if err == nil {
		work.Done <- util.GeneratedCode{Generator: "GenerateAdminProto", Code: permissions, Filename: "services/admin/service_admin_permissions.go"}
	} else {
		work.Done <- util.GeneratedCode{Generator: "GenerateAdminProto", Error: err}
	}

	code, err := util.ExecuteTemplate("admin/service_admin.go.tmpl", struct {
		Entities     []util.Entity
		GenerateAuth bool
		ImportPath   string
	}{r.Entities, generateAuth, util.AppImportPath})
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

func genUTF8List() []string {
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
