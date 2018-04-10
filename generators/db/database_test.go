package db

import (
	"strings"
	"testing"

	"github.com/fluxynet/gocipe/generators"
)

func TestGenerateDatabase(t *testing.T) {
	structInfo := generators.StructureInfo{
		TableName: "Persons",
		Fields: []generators.FieldInfo{
			{Name: "id", Type: "int64"},
			{Name: "name", Type: "string"},
			{Name: "email", Type: "string"},
			{Name: "gender", Type: "string"},
		},
	}

	output, err := GenerateDatabase(structInfo)
	expected := `
DROP TABLE IF EXISTS Persons;

CREATE TABLE Persons (
	"id"  PRIMARY KEY NOT NULL,
	"name"  NOT NULL,
	"email"  NOT NULL,
	"gender"  NOT NULL
);
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
