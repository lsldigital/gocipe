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

		for i, field := range entity.Fields {
			t := field.Property.Type
			if field.Property.Type == "time" {
				hasTime = true
				t = "google.protobuf.Timestamp"
			}
			ent.Fields = append(ent.Fields, protoField{Index: i + 1, Name: field.Property.Name, Type: t})
		}
		ents = append(ents, ent)
	}

	if hasTime {
		imports = append(imports, `import "google/protobuf/timestamp.proto";`)
	}

	return util.ExecuteTemplate("crud/protobuf.proto.tmpl", struct {
		Entities []protoEntity
		Imports  []string
	}{ents, imports})
}
