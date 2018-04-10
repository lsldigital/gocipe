package crud

import (
	"strings"
	"testing"

	"github.com/fluxynet/gocipe/generators"
)

func TestGenerateSave(t *testing.T) {
	structInfo := generators.StructureInfo{
		Name: "Persons",
		Fields: []generators.FieldInfo{
			{Name: "id", Type: "int64"},
			{Name: "name", Type: "string"},
			{Name: "email", Type: "string"},
			{Name: "gender", Type: "string"},
		},
	}

	output, err := GenerateSave(structInfo)
	expected := `
// Save either inserts or updates a Persons record based on whether or not id is nil
func (entity *Persons) Save() error {
	if entity.ID == nil {
		return entity.Insert()
	}
	return entity.Update()
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
