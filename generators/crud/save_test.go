package crud

import (
	"strings"
	"testing"
)

func TestGenerateSave(t *testing.T) {
	name := "Person"
	fields := []string{
		"id", "name", "email", "gender",
	}

	output, err := GenerateSave(name, fields)
	expected := `
//Save will persist Person entity to the database
func (entity *Person) Save(db *sql.DB) error {
	if entity.id == 0 {
		error := entity.Insert(db)
	} else {
		error := entity.Update(db)
	}

	if error != nil {
		return error
	}

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
