package crud

import (
	"strings"
	"testing"

	"github.com/fluxynet/gocipe/generators"
)

func TestGenerateDelete(t *testing.T) {
	structInfo := generators.StructureInfo{
		Name: "Persons",
		Fields: []generators.FieldInfo{
			{Name: "id", Type: "int64"},
			{Name: "name", Type: "string"},
			{Name: "email", Type: "string"},
			{Name: "gender", Type: "string"},
		},
	}

	output, err := GenerateDelete(structInfo)
	expected := `
// Delete deletes a Persons record from database by id primary key
func Delete(id int64) error {
	_, err := db.Exec("DELETE FROM  WHERE id = $1", id)
	return err
}

// Delete deletes a Persons record from database and sets id to nil
func (entity *Persons) Delete() error {
	_, err := db.Exec("DELETE FROM  WHERE id = $1", entity.ID)
	if err != nil {
		entity.ID = nil
	}
	return err
}
`

	if err != nil {
		t.Errorf("Got error: %s", err)
	} else if strings.Compare(output, expected) != 0 {
		t.Errorf(`Output not as expected. Output length = %d Expected length = %d
--- # Output ---
%s
--- # Expected ---
%s
------------------
`, len(output), len(expected), output, expected)
	}
}
