package admin

import (
	"fmt"
	"strings"

	"github.com/fluxynet/gocipe/output"
	"github.com/fluxynet/gocipe/util"
)

type fileField struct {
	Entity string
	Field  string
}

// Generate returns generated code for a Admin service
func Generate(out *output.Output, r *util.Recipe) {
	if !r.Admin.Generate {
		return
	}

	for _, entity := range r.Entities {
		if !entity.Admin.Generate {
			continue
		}

		if entity.HasAdminHooks() {
			out.GenerateAndOverwrite("GenerateAdmin Hooks "+entity.Name, "admin/service_admin_hooks.go.tmpl", fmt.Sprintf(
				"services/admin/%s_hooks.gocipe.go", strings.ToLower(entity.Name)), output.WithoutHeader, struct {
				Entity     util.Entity
				ImportPath string
			}{
				Entity:     entity,
				ImportPath: util.AppImportPath,
			})
		}

	}

	if r.HasFileFields() || r.HasContentFileUpload() {
		out.GenerateAndOverwrite("GenerateAdmin Upload", "admin/admin_config_upload.go.tmpl", "services/admin/service_admin_config_upload.gocipe.go", output.WithHeader, struct {
			Entities []util.Entity
		}{r.Entities})
	}

	out.GenerateAndOverwrite("GenerateAdmin Helpers", "admin/admin_helpers.go.tmpl", "services/admin/admin_helpers.gocipe.go", output.WithHeader, struct {
		ImportPath string
	}{util.AppImportPath})

	out.GenerateAndOverwrite("GenerateAdmin Proto", "admin/service_admin.proto.tmpl", "proto/service_admin.proto", output.WithHeader, struct {
		ImportPath string
		Entities   []util.Entity
	}{util.AppImportPath, r.Entities})

	out.GenerateAndOverwrite("GenerateAdmin Permissions", "admin/admin_permissions.go.tmpl", "services/admin/service_admin_permissions.go", output.WithHeader, struct {
		ImportPath  string
		Permissions []util.Permission
	}{util.AppImportPath, r.GetPermissions()})

	out.GenerateAndOverwrite("GenerateAdmin", "admin/service_admin.go.tmpl", "services/admin/service_admin.gocipe.go", output.WithHeader, struct {
		Entities     []util.Entity
		GenerateAuth bool
		ImportPath   string
	}{r.Entities, r.HasAuth(), util.AppImportPath})
}
