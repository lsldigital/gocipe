package simple

import (
	"strings"
	"testing"
)

func TestGenerateUpdate(t *testing.T) {
	name := "Person"
	fields := []string{
		"id", "name", "email", "gender",
	}

	output, err := GenerateUpdate(name, fields)
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
