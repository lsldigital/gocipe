package admin

import (
	"fmt"
	"strings"

	"github.com/fluxynet/gocipe/util"

	"github.com/fluxynet/gocipe/output"
)

type fileField struct {
	Entity string
	Field  string
}

// Generate returns generated code for a Admin service
func Generate(out output.Output, r *util.Recipe) error {
	if !r.Admin.Generate {
		return nil
	}

	for _, entity := range r.Entities {
		if !entity.Admin.Generate {
			continue
		}

		if entity.HasCrudHooks() {

			out.GenerateAndOverwrite("GenerateAdminHooks"+entity.Name, "admin/service_admin_hooks.go.tmpl", fmt.Sprintf(
				"services/admin/%s_hooks.gocipe.go", strings.ToLower(entity.Name)), struct {
				Entity     util.Entity
				ImportPath string
			}{
				Entity:     entity,
				ImportPath: util.AppImportPath,
			})

		}

	}

	if r.HasFileFields() {
		// code, err := util.ExecuteTemplate("admin/admin_config_upload.go.tmpl", struct {
		// 	Entities []util.Entity
		// }{r.Entities})

		out.GenerateAndOverwrite("GenerateAdmin Upload", "admin/admin_config_upload.go.tmpl", "services/admin/service_admin_config_upload.gocipe.go", struct {
			Entities []util.Entity
		}{r.Entities})

		// work.Waitgroup.Add(1)
		// if err == nil {
		// 	work.Done <- util.GeneratedCode{Generator: "GenerateAdmin Upload", Code: code, Filename: "services/admin/service_admin_config_upload.gocipe.go", NoOverwrite: true}
		// } else {
		// 	work.Done <- util.GeneratedCode{Generator: "GenerateAdmin Upload", Error: fmt.Errorf("failed to execute template: %s", err)}
		// }
	}

	// generate admin_helpers.gocipe.go
	// helpers, err := util.ExecuteTemplate("admin/admin_helpers.go.tmpl", struct {
	// 	ImportPath string
	// }{util.AppImportPath})

	out.GenerateAndOverwrite("GenerateAdmin Helpers", "admin/admin_helpers.go.tmpl", "services/admin/admin_helpers.gocipe.go", struct {
		ImportPath string
	}{util.AppImportPath})
	// work.Waitgroup.Add(1)
	// if err == nil {
	// 	work.Done <- util.GeneratedCode{Generator: "GenerateAdmin Helpers", Code: helpers, Filename: "services/admin/admin_helpers.gocipe.go"}
	// } else {
	// 	work.Done <- util.GeneratedCode{Generator: "GenerateAdmin Helpers", Error: fmt.Errorf("failed to execute template: %s", err)}
	// }

	// generate admin.proto
	// proto, err := util.ExecuteTemplate("admin/service_admin.proto.tmpl", struct {
	// 	ImportPath string
	// 	Entities   []util.Entity
	// }{util.AppImportPath, r.Entities})
	out.GenerateAndOverwrite("GenerateAdminProto", "admin/service_admin.proto.tmpl", "proto/service_admin.proto", struct {
		ImportPath string
	}{util.AppImportPath})
	// work.Waitgroup.Add(1)
	// if err == nil {
	// 	work.Done <- util.GeneratedCode{Generator: "GenerateAdminProto", Code: proto, Filename: "proto/service_admin.proto"}
	// } else {
	// 	work.Done <- util.GeneratedCode{Generator: "GenerateAdminProto", Error: err}
	// }

	// generate admin_permissions.go
	// permissions, err := util.ExecuteTemplate("admin/admin_permissions.go.tmpl", struct {
	// 	ImportPath  string
	// 	Permissions []util.Permission
	// }{util.AppImportPath, r.GetPermissions()})

	out.GenerateAndOverwrite("GenerateAdminProto", "admin/admin_permissions.go.tmpl", "services/admin/service_admin_permissions.go", struct {
		ImportPath  string
		Permissions []util.Permission
	}{util.AppImportPath, r.GetPermissions()})
	// work.Waitgroup.Add(1)
	// if err == nil {
	// 	work.Done <- util.GeneratedCode{Generator: "GenerateAdminProto", Code: permissions, Filename: "services/admin/service_admin_permissions.go"}
	// } else {
	// 	work.Done <- util.GeneratedCode{Generator: "GenerateAdminProto", Error: err}
	// }

	// code, err := util.ExecuteTemplate("admin/service_admin.go.tmpl", struct {
	// 	Entities     []util.Entity
	// 	GenerateAuth bool
	// 	ImportPath   string
	// }{r.Entities, r.HasAuth(), util.AppImportPath})

	out.GenerateAndOverwrite("GenerateAdmin", "admin/service_admin.go.tmpl", "services/admin/service_admin.gocipe.go", struct {
		Entities     []util.Entity
		GenerateAuth bool
		ImportPath   string
	}{r.Entities, r.HasAuth(), util.AppImportPath})
	// work.Waitgroup.Add(1)
	// if err == nil {
	// 	work.Done <- util.GeneratedCode{Generator: "GenerateAdmin", Code: code, Filename: "services/admin/service_admin.gocipe.go"}
	// } else {
	// 	work.Done <- util.GeneratedCode{Generator: "GenerateAdmin", Error: fmt.Errorf("failed to execute template: %s", err)}
	// }

	// work.Waitgroup.Done()
	return nil
}
