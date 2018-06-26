package crud

import (
	"fmt"
	"strings"

	"github.com/fluxynet/gocipe/util"
)

// generateUpdate produces code for database update of entity (UPDATE)
func generateUpdate(entities map[string]util.Entity, entity util.Entity) (string, error) {
	var (
		before, after, post, sqlfields, structfields []string
		count                                        = 2
	)
	structfields = append(structfields, fmt.Sprintf("&entity.%s", "ID"))

	for _, field := range entity.Fields {
		if field.Property.Name == "UpdatedAt" {
			before = append(before, "entity.UpdatedAt = ptypes.TimestampNow()")
		}

		sqlfields = append(sqlfields, fmt.Sprintf("%s = $%d", field.Schema.Field, count))

		if field.Property.Type == "time" {
			prop := strings.ToLower(field.Property.Name)
			before = append(before, fmt.Sprintf("%s, _ := ptypes.Timestamp(entity.%s)", prop, field.Property.Name))
			structfields = append(structfields, fmt.Sprintf("%s", prop))
		} else {
			structfields = append(structfields, fmt.Sprintf("entity.%s", field.Property.Name))
		}
		count++
	}

	for _, rel := range entity.Relationships {
		if rel.Type == util.RelationshipTypeManyMany {
			post = append(post, fmt.Sprintf("repo.Save%s(ctx, entity.ID, entity.%s, tx, false)", util.RelFuncName(rel), rel.Name))
		}
	}

	return util.ExecuteTemplate("crud/partials/update.go.tmpl", struct {
		Before        []string
		After         []string
		Post          []string
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
		Post:          post,
		After:         after,
		SQLFields:     strings.Join(sqlfields, ", "),
		StructFields:  strings.Join(structfields, ", "),
		HasPreHook:    entity.Crud.Hooks.PreSave,
		HasPostHook:   entity.Crud.Hooks.PostSave,
		Relationships: nil,
	})
}
