package crud

import (
	"fmt"
	"strings"

	"github.com/fluxynet/gocipe/util"
)

// generateGet produces code for database retrieval of single entity (SELECT WHERE id)
func generateGet(entities map[string]util.Entity, entity util.Entity) (string, error) {
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

	for _, rel := range entity.Relationships {
		related := entities[rel.Entity]
		if rel.Type == util.RelationshipTypeManyOne {
			structfields = append(structfields, fmt.Sprintf("&entity.%s", rel.Name+"ID"))
			sqlfields = append(sqlfields, fmt.Sprintf("%s", strings.ToLower(related.Name)+"_id"))
		}
	}

	return util.ExecuteTemplate("crud/partials/get.go.tmpl", struct {
		EntityName   string
		SQLFields    string
		Table        string
		StructFields string
		PrimaryKey   string
		Before       []string
		After        []string
		HasPreHook   bool
		HasPostHook  bool
	}{
		EntityName:   entity.Name,
		Table:        entity.Table,
		SQLFields:    strings.Join(sqlfields, ", "),
		StructFields: strings.Join(structfields, ", "),
		PrimaryKey:   entity.PrimaryKey,
		Before:       before,
		After:        after,
		HasPreHook:   entity.Crud.Hooks.PreRead,
		HasPostHook:  entity.Crud.Hooks.PostRead,
	})
}
