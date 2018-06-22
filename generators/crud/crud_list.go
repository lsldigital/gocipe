package crud

import (
	"fmt"
	"strings"

	"github.com/fluxynet/gocipe/util"
)

// generateList produces code for database retrieval of list of entities with optional filters
func generateList(entities map[string]util.Entity, entity util.Entity) (string, error) {
	var sqlfields, structfields, before, after []string

	sqlfields = append(sqlfields, fmt.Sprintf("%s", "id"))
	structfields = append(structfields, fmt.Sprintf("&entity.%s", "ID"))

	for _, field := range entity.Fields {
		if field.Property.Type == "time" {
			prop := strings.ToLower(field.Property.Name)
			before = append(before, fmt.Sprintf("var %s time.Time", prop))
			structfields = append(structfields, fmt.Sprintf("&%s", prop))
			after = append(after, fmt.Sprintf("entity.%s, _ = ptypes.TimestampProto(%s)", field.Property.Name, prop))
		} else {
			structfields = append(structfields, fmt.Sprintf("&entity.%s", field.Property.Name))
		}
		sqlfields = append(sqlfields, fmt.Sprintf("%s", field.Schema.Field))
	}

	return util.ExecuteTemplate("crud/partials/list.go.tmpl", struct {
		EntityName   string
		SQLFields    string
		StructFields string
		Table        string
		Before       []string
		After        []string
		HasPreHook   bool
		HasPostHook  bool
	}{
		EntityName:   entity.Name,
		Table:        entity.Table,
		SQLFields:    strings.Join(sqlfields, ", "),
		StructFields: strings.Join(structfields, ", "),
		Before:       before,
		After:        after,
		HasPreHook:   entity.Crud.Hooks.PreList,
		HasPostHook:  entity.Crud.Hooks.PostList,
	})
}
