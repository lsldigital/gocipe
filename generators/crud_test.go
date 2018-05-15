package generators

import (
	"fmt"
	"testing"

	"github.com/fluxynet/gocipe/util"
)

func TestCrub(t *testing.T) {

	pack := "person"

	fieldProperty := util.FieldProperty{
		Name: "Name",
		Type: "string",
	}

	fieldSchema := util.FieldSchema{
		Field:    "name",
		Type:     "VARCHAR(255)",
		Nullable: true,
		Default:  "jeshta",
	}

	fieldRelationshipTarget := util.FieldRelationshipTarget{
		Entity: "Color",
		Table:  "color",
		ThisID: "person_id",
		ThatID: "color_id",
	}

	relationship := util.FieldRelationship{
		Type:   "many-to-many",
		Target: fieldRelationshipTarget,
	}

	field := util.Field{
		Label:        "Person",
		Serialized:   "person",
		Property:     fieldProperty,
		Schema:       fieldSchema,
		Relationship: relationship,
	}
	fields := []util.Field{}

	f := append(fields, field)

	schema := util.SchemaOpts{
		Create:    true,
		Drop:      true,
		Aggregate: true,
		Path:      "some/path/schema",
	}

	crud_hooks := util.CrudHooks{
		PreCreate:  false,
		PostCreate: false,
		PreRead:    false,
		PostRead:   false,
		PreList:    false,
		PostList:   false,
		PreUpdate:  false,
		PostUpdate: false,
		PreDelete:  false,
		PostDelete: false,
	}

	crud := util.CrudOpts{
		Create:   true,
		Read:     true,
		ReadList: true,
		Update:   true,
		Delete:   true,
		Hooks:    crud_hooks,
	}

	hooks := util.RestHooks{
		PreCreate:  false,
		PostCreate: false,
		PreRead:    false,
		PostRead:   false,
		PreList:    false,
		PostList:   false,
		PreUpdate:  false,
		PostUpdate: false,
		PreDelete:  false,
		PostDelete: false,
	}

	restOpts := util.RestOpts{
		Create:   false,
		Read:     false,
		ReadList: false,
		Update:   false,
		Delete:   false,
		Prefix:   "/api",
		Hooks:    hooks,
	}

	entity := util.Entity{}
	entity.Name = "Person"
	entity.Table = "person"
	entity.Description = "this is descr"
	entity.Fields = f
	entity.Schema = &schema
	entity.Crud = &crud
	entity.Rest = &restOpts

	// pkg := ""
	structure, _ := util.ExecuteTemplate("crud_structure.go.tmpl", struct {
		Entity  util.Entity
		Package string
	}{entity, pack})

	fmt.Println(structure)

}
