package crud

import (
	"strings"
	"testing"

	"github.com/fluxynet/gocipe/generators"
)

func TestGenerateList(t *testing.T) {
	structInfo := generators.StructureInfo{
		Name: "Persons",
		Fields: []generators.FieldInfo{
			{Name: "id", Type: "int64", Comments: ""},
			{Name: "name", Type: "string", Comments: ""},
			{Name: "email", Type: "string", Comments: ""},
			{Name: "gender", Type: "string", Comments: ""},
		},
	}

	output, err := GenerateList(structInfo)
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

	entities = []Person

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
