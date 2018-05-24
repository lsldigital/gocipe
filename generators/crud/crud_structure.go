package crud

import (
	"fmt"
	"strings"

	"github.com/fluxynet/gocipe/util"
)

func generateStructure(entities map[string]util.Entity, entity util.Entity) (string, error) {
	var fields []struct {
		Name       string
		Type       string
		Serialized string
	}

	var initialization []struct {
		Name string
		Type string
	}

	for _, field := range entity.Fields {
		fields = append(fields, struct {
			Name       string
			Type       string
			Serialized string
		}{field.Property.Name, "*" + field.Property.Type, field.Serialized})

		initialization = append(initialization, struct {
			Name string
			Type string
		}{field.Property.Name, fmt.Sprintf("new(%s)", field.Property.Type)})
	}

	relation := func(rel util.Relationship, many, full bool) {
		var t string
		// if full {
		t = fmt.Sprintf("%s", rel.Entity)
		// } else {
		// 	t, _ = util.GetPrimaryKeyDataType(entities[rel.Entity].PrimaryKey)
		// }

		if many {
			t = "[]" + t
		} else {
			t = "*" + t
		}

		fields = append(fields, struct {
			Name       string
			Type       string
			Serialized string
		}{rel.Name, t, rel.Serialized})

	}

	for _, rel := range entity.Relationships {
		switch rel.Type {
		case util.RelationshipTypeOneOne:
			fallthrough
		case util.RelationshipTypeManyOne:
			relation(rel, false, rel.Full)
		case util.RelationshipTypeOneMany:
			fallthrough
		case util.RelationshipTypeManyMany:
			relation(rel, true, rel.Full)
		}
	}

	return util.ExecuteTemplate("crud/partials/structure.go.tmpl", struct {
		Package     string
		Name        string
		Description string
		PrimaryKey  string
		Fields      []struct {
			Name       string
			Type       string
			Serialized string
		}
		Initialization []struct {
			Name string
			Type string
		}
	}{
		Package:        strings.ToLower(entity.Name),
		Name:           entity.Name,
		Description:    entity.Description,
		PrimaryKey:     entity.PrimaryKey,
		Fields:         fields,
		Initialization: initialization,
	})
}
