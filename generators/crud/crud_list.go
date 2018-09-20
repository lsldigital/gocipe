package crud

import (
	"fmt"
	"strings"

	"github.com/fluxynet/gocipe/util"
)

// generateList produces code for database retrieval of list of entities with optional filters
func generateList(entities map[string]util.Entity, entity util.Entity) (string, error) {
	var sqlfields, structfields, before, after, orderfields []string
	related := make(map[string]string)

	sqlfields = append(sqlfields, fmt.Sprintf(`t."%s"`, "id"))
	structfields = append(structfields, fmt.Sprintf("&entity.%s", "ID"))
	orderfields = append(orderfields, "id")

	for _, field := range entity.Fields {
		if field.Property.Type == "time" {
			prop := strings.ToLower(field.Property.Name)
			before = append(before, fmt.Sprintf("var %s time.Time", prop))
			structfields = append(structfields, fmt.Sprintf("&%s", prop))
			after = append(after, fmt.Sprintf("entity.%s, _ = ptypes.TimestampProto(%s)", field.Property.Name, prop))
		} else {
			structfields = append(structfields, fmt.Sprintf("&entity.%s", field.Property.Name))
		}
		orderfields = append(orderfields, field.Schema.Field)
		sqlfields = append(sqlfields, fmt.Sprintf(`t."%s"`, field.Schema.Field))
	}

	for _, rel := range entity.Relationships {
		switch rel.Type {
		case util.RelationshipTypeManyMany,
			util.RelationshipTypeManyManyOwner,
			util.RelationshipTypeManyManyInverse,
			util.RelationshipTypeOneMany,
			util.RelationshipTypeManyOne:
			related[rel.Name] = fmt.Sprintf("err = repo.Load%s(ctx, entities...)", util.RelFuncName(rel))
		}
		if rel.Type == util.RelationshipTypeManyOne {
			sqlfields = append(sqlfields, fmt.Sprintf(`t."%s"`, rel.ThisID))
			structfields = append(structfields, fmt.Sprintf("&entity.%sID", rel.Name))
		}
	}

	return util.ExecuteTemplate("crud/partials/list.go.tmpl", struct {
		EntityName   string
		SQLFields    string
		StructFields string
		OrderFields  string
		Table        string
		Before       []string
		After        []string
		Related      map[string]string
		HasPreHook   bool
		HasPostHook  bool
		DefaultSort  string
	}{
		EntityName:   entity.Name,
		Table:        entity.Table,
		SQLFields:    strings.Join(sqlfields, ", "),
		StructFields: strings.Join(structfields, ", "),
		OrderFields:  `"` + strings.Join(orderfields, `","`) + `"`,
		Before:       before,
		After:        after,
		Related:      related,
		HasPreHook:   entity.CrudHooks.PreList,
		HasPostHook:  entity.CrudHooks.PostList,
		DefaultSort:  entity.DefaultSort,
	})
}
