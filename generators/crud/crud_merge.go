package crud

import (
	"fmt"
	"strings"

	"github.com/fluxynet/gocipe/util"
)

// generateMerge produces code for database merge of entity (INSERT/ON CONFLICT UPDATE)
func generateMerge(entities map[string]util.Entity, entity util.Entity) (string, error) {
	var (
		before, related, sqlfieldsInsert, sqlfieldsUpdate, sqlPlaceholders, structFields []string
		count                                                                            = 1
	)

	sqlfieldsInsert = append(sqlfieldsInsert, `"id"`)
	structFields = append(structFields, "entity.ID")
	sqlPlaceholders = append(sqlPlaceholders, fmt.Sprintf("$%d", count))
	count++

	for _, field := range entity.Fields {
		if field.Property.Name == "CreatedAt" {
			before = append(before, "entity.CreatedAt = ptypes.TimestampNow()")
		} else if field.Property.Name == "UpdatedAt" {
			before = append(before, "entity.UpdatedAt = ptypes.TimestampNow()")
		}

		sqlPlaceholders = append(sqlPlaceholders, fmt.Sprintf("$%d", count))
		sqlfieldsUpdate = append(sqlfieldsUpdate, fmt.Sprintf(`"%s" = $%d`, field.Schema.Field, count))
		sqlfieldsInsert = append(sqlfieldsInsert, fmt.Sprintf(`"%s"`, field.Schema.Field))

		if field.Property.Type == "time" {
			prop := strings.ToLower(field.Property.Name)
			before = append(before, fmt.Sprintf("%s, _ := ptypes.Timestamp(entity.%s)", prop, field.Property.Name))
			structFields = append(structFields, fmt.Sprintf("%s", prop))
		} else {
			structFields = append(structFields, fmt.Sprintf("entity.%s", field.Property.Name))
		}

		count++
	}

	for _, rel := range entity.Relationships {
		if rel.Type == util.RelationshipTypeManyMany ||
			rel.Type == util.RelationshipTypeManyManyOwner ||
			rel.Type == util.RelationshipTypeManyManyInverse ||
			rel.Type == util.RelationshipTypeOneMany {
			related = append(related, fmt.Sprintf("repo.Save%s(ctx, tx, false, entity.ID, entity.%s...)", util.RelFuncName(rel), rel.Name))
		} else if rel.Type == util.RelationshipTypeManyOne {
			sqlPlaceholders = append(sqlPlaceholders, fmt.Sprintf("$%d", count))
			sqlfieldsUpdate = append(sqlfieldsUpdate, fmt.Sprintf(`"%s" = $%d`, rel.ThisID, count))
			sqlfieldsInsert = append(sqlfieldsInsert, fmt.Sprintf(`"%s"`, rel.ThisID))
			structFields = append(structFields, fmt.Sprintf("entity.%sID", rel.Name))
			count++
		}
	}

	return util.ExecuteTemplate("crud/partials/merge.go.tmpl", struct {
		EntityName      string
		PrimaryKey      string
		Before          []string
		Related         []string
		Table           string
		SQLFieldsInsert string
		SQLPlaceholders string
		SQLFieldsUpdate string
		StructFields    string
		HasPreHook      bool
		HasPostHook     bool
	}{
		EntityName:      entity.Name,
		PrimaryKey:      entity.PrimaryKey,
		Before:          before,
		Related:         related,
		Table:           entity.Table,
		SQLFieldsInsert: strings.Join(sqlfieldsInsert, ", "),
		SQLPlaceholders: strings.Join(sqlPlaceholders, ", "),
		SQLFieldsUpdate: strings.Join(sqlfieldsUpdate, ", "),
		StructFields:    strings.Join(structFields, ", "),
		HasPreHook:      entity.CrudHooks.PreSave,
		HasPostHook:     entity.CrudHooks.PostSave,
	})
}
