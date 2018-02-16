package crud

import (
	"strings"
	"testing"
)

func TestGenerateList(t *testing.T) {
	name := "Person"
	fields := []string{
		"id", "name", "email", "gender",
	}

	output, err := GenerateList(name, fields)
	expected := `
//List returns all Person entities stored in database
func List(db *sql.DB) ([]Person, error) {
	var (
		entity   *Person
		entities []Person
	)

	query, err := db.Query("SELECT id, name, email, gender FROM ` + "`persons`" + `")
	if err != nil {
		return entities, err
	}
	defer query.Close()

	entities = make([]Person, 10)

	for query.Next() {
		entity = new(Person)
		query.Scan(entity.id, entity.name, entity.email, entity.gender)
		entities = append(entities, *entity)
	}

	return entities, nil
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
