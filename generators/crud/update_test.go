package crud

import (
	"strings"
	"testing"

	"github.com/fluxynet/gocipe/generators"
)

func TestGenerateUpdate(t *testing.T) {
	structInfo := generators.StructureInfo{
		Name: "Persons",
		Fields: []generators.FieldInfo{
			{Name: "id", Type: "int64"},
			{Name: "name", Type: "string"},
			{Name: "email", Type: "string"},
			{Name: "gender", Type: "string"},
		},
	}

	output, err := GenerateUpdate(structInfo)
	expected := `
//Update Will execute an SQLUpdate Statement for Persons in the database. Prefer using Save instead of Update directly.
func (entity *Persons) Update() error {
	_, err := db.Exec("UPDATE  SET id = $2, name = $3, email = $4, gender = $5 WHERE id = $1", entity.ID, *entity., *entity., *entity., *entity.)

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
