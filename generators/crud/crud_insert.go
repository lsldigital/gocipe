package crud

import (
	"fmt"
	"strings"

	"github.com/fluxynet/gocipe/util"
)

// generateInsert produces code for database insertion of entity (INSERT INTO)
func generateInsert(entities map[string]util.Entity, entity util.Entity) (string, error) {
	var before, after, sqlPlaceholders, sqlfields, structFields []string
	count := 1

	if entity.PrimaryKey != util.PrimaryKeySerial {
		sqlPlaceholders = append(sqlPlaceholders, fmt.Sprintf("$%d", count))
		sqlfields = append(sqlfields, "id")
		structFields = append(structFields, "entity.ID")

		count++
	}

	for _, field := range entity.Fields {
		if field.Property.Name == "CreatedAt" {
			before = append(before, "entity.CreatedAt = timestamp.TimestampNow()")
		} else if field.Property.Name == "UpdatedAt" {
			before = append(before, "entity.UpdatedAt = timestamp.TimestampNow()")
		}

		sqlPlaceholders = append(sqlPlaceholders, fmt.Sprintf("$%d", count))
		sqlfields = append(sqlfields, fmt.Sprintf("%s", field.Schema.Field))

		if field.Property.Type == "time" {
			prop := strings.ToLower(field.Property.Name)
			before = append(before, fmt.Sprintf("%s, _ := ptypes.Timestamp(entity.%s)", prop, field.Property.Name))
			structFields = append(structFields, fmt.Sprintf("%s", prop))
		} else {
			structFields = append(structFields, fmt.Sprintf("entity.%s", field.Property.Name))
		}

		count++
	}

	return util.ExecuteTemplate("crud/partials/insert.go.tmpl", struct {
		Before          []string
		After           []string
		EntityName      string
		PrimaryKey      string
		SQLFields       string
		SQLPlaceholders string
		StructFields    string
		Table           string
		HasPreHook      bool
		HasPostHook     bool
	}{
		Before:          before,
		After:           after,
		EntityName:      entity.Name,
		PrimaryKey:      entity.PrimaryKey,
		SQLFields:       strings.Join(sqlfields, ", "),
		SQLPlaceholders: strings.Join(sqlPlaceholders, ", "),
		StructFields:    strings.Join(structFields, ", "),
		Table:           entity.Table,
		HasPostHook:     entity.Crud.Hooks.PreSave,
		HasPreHook:      entity.Crud.Hooks.PostSave,
	})
}
