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
			{Name: "id", Type: "int64"},
			{Name: "name", Type: "string"},
			{Name: "email", Type: "string"},
			{Name: "gender", Type: "string"},
		},
	}

	output, err := GenerateInsert(structInfo)
	expected := `
// Insert performs an SQL insert for Persons record and update instance with inserted id.
// Prefer using Save rather than Insert directly.
func (entity *Persons) Insert() error {
	var (
		id  int64
		err error
	)

	err = db.QueryRow("INSERT INTO  (name, email, gender) VALUES ($1, $2, $3) RETURNING id", *entity., *entity., *entity.).Scan(&id)

	if err == nil {
		entity.ID = &id
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
