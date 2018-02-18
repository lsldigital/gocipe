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
			{Name: "id", Type: "int64", Comments: ""},
			{Name: "name", Type: "string", Comments: ""},
			{Name: "email", Type: "string", Comments: ""},
			{Name: "gender", Type: "string", Comments: ""},
		},
	}

	output, err := GenerateUpdate(structInfo)
	expected := `
//Update Will execute an SQLUpdate Statement in the database. Prefer using Save instead of Update directly.
func (entity *Person) Update(db *sql.DB) error {
	stmt, err := db.Prepare("UPDATE ` + "`persons`" + ` SET name = ?, email = ?, gender = ? WHERE id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(entity.name, entity.email, entity.gender, entity.id)
	if err != nil {
		return err
	}
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
