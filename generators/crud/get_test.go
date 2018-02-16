package crud

import (
	"strings"
	"testing"
)

func TestGenerateGet(t *testing.T) {
	name := "Person"
	fields := []string{
		"id", "name", "email", "gender",
	}

	output, err := GenerateGet(name, fields)
	expected := `
//Get returns a single Person from database
func Get(db *sql.DB, id int) (*Person, error) {
	var entity = new(Person)

	query := db.QueryRow("SELECT id, name, email, gender FROM ` + "`persons`" + ` WHERE id = ? LIMIT 1", id)
	err := query.Scan(entity.id, entity.name, entity.email, entity.gender)

	return entity, err
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
