package crud

import (
	"github.com/fluxynet/gocipe/util"
)

func generateProtobuf(entities map[string]util.Entity) (string, error) {
	type protoField struct {
		Name  string
		Type  string
		Index int
	}

	type protoEntity struct {
		Name        string
		Description string
		Fields      []protoField
	}
	var (
		ents    []protoEntity
		hasTime bool
		imports []string
	)

	for _, entity := range entities {
		var ent = protoEntity{Name: entity.Name, Description: entity.Description}
		count := 1

		ent.Fields = append(ent.Fields, protoField{Index: count, Name: "ID", Type: "string"})
		count++

		for _, field := range entity.Fields {
			t := field.Property.Type
			if field.Property.Type == "time" {
				hasTime = true
				t = "google.protobuf.Timestamp"
			}
			ent.Fields = append(ent.Fields, protoField{Index: count, Name: field.Property.Name, Type: t})
			count++
		}

		for _, rel := range entity.Relationships {
			var t string
			related := entities[rel.Entity]

			t = related.Name

			switch rel.Type {
			case util.RelationshipTypeManyMany, util.RelationshipTypeManyManyOwner, util.RelationshipTypeManyManyInverse:
				ent.Fields = append(ent.Fields, protoField{Index: count, Name: rel.Name, Type: "repeated " + t})
			case util.RelationshipTypeOneOne:
				ent.Fields = append(ent.Fields, protoField{Index: count, Name: rel.Name, Type: t})
			case util.RelationshipTypeOneMany:
				ent.Fields = append(ent.Fields, protoField{Index: count, Name: rel.Name, Type: "repeated " + t})
			case util.RelationshipTypeManyOne:
				ent.Fields = append(ent.Fields, protoField{Index: count, Name: rel.Name + "ID", Type: f})
				count++

				ent.Fields = append(ent.Fields, protoField{Index: count, Name: rel.Name, Type: t})
			}

			count++
		}

		for _, ref := range entity.References {
			// IDField
			ent.Fields = append(ent.Fields, protoField{Index: count, Name: ref.IDField.Property.Name, Type: ref.IDField.Property.Type})
			count++

			// TypeField
			ent.Fields = append(ent.Fields, protoField{Index: count, Name: ref.TypeField.Property.Name, Type: ref.TypeField.Property.Type})
			count++
		}

		ents = append(ents, ent)
	}

	if hasTime {
		imports = append(imports, `import "google/protobuf/timestamp.proto";`)
	}

	return util.ExecuteTemplate("crud/protobuf.proto.tmpl", struct {
		Entities      []protoEntity
		Imports       []string
		AppImportPath string
	}{ents, imports, util.AppImportPath})
}
