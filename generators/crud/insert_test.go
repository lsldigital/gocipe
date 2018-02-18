package crud

import (
	"strings"
	"testing"

	"github.com/fluxynet/gocipe/generators"
)

func TestGenerateInsert(t *testing.T) {
	structInfo := generators.StructureInfo{
		Name: "Persons",
		Fields: []generators.FieldInfo{
			{Name: "id", Type: "int64", Comments: ""},
			{Name: "name", Type: "string", Comments: ""},
			{Name: "email", Type: "string", Comments: ""},
			{Name: "gender", Type: "string", Comments: ""},
		},
	}

	output, err := GenerateInsert(structInfo)
	expected := `
//Insert Will execute an SQLInsert Statement in the database. Prefer using Save instead of Insert directly.
func (entity *Person) Insert(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO ` + "`persons`" + ` (name, email, gender) VALUES (?, ?, ?)")

	if err != nil {
		return err
	}

	result, err := stmt.Exec(entity.name, entity.email, entity.gender)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	entity.id = id
	return nil
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
