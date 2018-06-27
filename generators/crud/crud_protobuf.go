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
		pkey, err := util.GetPrimaryKeyDataType(entity.PrimaryKey)
		count := 1

		if err != nil {
			return "", err
		}

		ent.Fields = append(ent.Fields, protoField{Index: count, Name: "ID", Type: pkey})
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
			case util.RelationshipTypeManyMany:
				ent.Fields = append(ent.Fields, protoField{Index: count, Name: rel.Name, Type: "repeated " + t})
			case util.RelationshipTypeOneOne:
				ent.Fields = append(ent.Fields, protoField{Index: count, Name: rel.Name, Type: t})
			case util.RelationshipTypeOneMany:
				ent.Fields = append(ent.Fields, protoField{Index: count, Name: rel.Name, Type: "repeated " + t})
			case util.RelationshipTypeManyOne:
				f, err := util.GetPrimaryKeyDataType(entities[rel.Entity].PrimaryKey)
				if err != nil {
					return "", err
				}
				ent.Fields = append(ent.Fields, protoField{Index: count, Name: rel.Name + "ID", Type: f})
				count++

				ent.Fields = append(ent.Fields, protoField{Index: count, Name: rel.Name, Type: t})
			}

			count++
		}

		ents = append(ents, ent)
	}

	if hasTime {
		imports = append(imports, `import "google/protobuf/timestamp.proto";`)
	}

	return util.ExecuteTemplate("crud/protobuf.proto.tmpl", struct {
		Entities          []protoEntity
		Imports           []string
		ProjectImportPath string
	}{ents, imports, util.ProjectImportPath})
}
