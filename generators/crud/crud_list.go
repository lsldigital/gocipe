package crud

import (
	"fmt"
	"strings"

	"github.com/fluxynet/gocipe/util"
)

// generateList produces code for database retrieval of list of entities with optional filters
func generateList(entities map[string]util.Entity, entity util.Entity) (string, error) {
	var sqlfields, structfields, before, after, related []string

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

	for _, rel := range entity.Relationships {
		switch rel.Type {
		case util.RelationshipTypeManyMany, util.RelationshipTypeOneMany, util.RelationshipTypeManyOne:
			related = append(related, fmt.Sprintf("err = repo.Load%s(ctx, entities...)", util.RelFuncName(rel)))
		}
	}

	return util.ExecuteTemplate("crud/partials/list.go.tmpl", struct {
		EntityName   string
		SQLFields    string
		StructFields string
		Table        string
		Before       []string
		After        []string
		Related      []string
		HasPreHook   bool
		HasPostHook  bool
	}{
		EntityName:   entity.Name,
		Table:        entity.Table,
		SQLFields:    strings.Join(sqlfields, ", "),
		StructFields: strings.Join(structfields, ", "),
		Before:       before,
		After:        after,
		Related:      related,
		HasPreHook:   entity.Crud.Hooks.PreList,
		HasPostHook:  entity.Crud.Hooks.PostList,
	})
}
