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
			{Name: "id", Type: "int64", Comments: ""},
			{Name: "name", Type: "string", Comments: ""},
			{Name: "email", Type: "string", Comments: ""},
			{Name: "gender", Type: "string", Comments: ""},
		},
	}

	output, err := GenerateDelete(structInfo)
	expected := `
//Delete delete single Person entity from database
func Delete(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM ` + "`persons`" + ` WHERE id = ?", id)
	
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
