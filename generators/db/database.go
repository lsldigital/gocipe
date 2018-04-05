package db

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/fluxynet/gocipe/generators"
)

//tableField represents an SQL table field
type tableField struct {
	Name        string
	Type        string
	Constraints string
	Separator   string
}

var tmplDatabase, _ = template.New("GenerateDatabase").Parse(`
DROP TABLE IF EXISTS {{.TableName}};

CREATE TABLE {{.TableName}} ({{range .Fields}}
	"{{.Name}}" {{.Type}}{{.Constraints}}{{.Separator}}{{end}}
);
`)

//GenerateDatabase return code for database sql
func GenerateDatabase(structInfo generators.StructureInfo) (string, error) {
	var (
		output    bytes.Buffer
		numFields = len(structInfo.Fields)
	)
	data := new(struct {
		TableName string
		Fields    []tableField
	})

	data.TableName = structInfo.TableName

	for i, sfield := range structInfo.Fields {
		var (
			field       tableField
			constraints []string
		)

		field.Name = sfield.Name
		field.Type = sfield.DBType

		if field.Name == "id" {
			constraints = append(constraints, "PRIMARY KEY")
		}

		if !sfield.Nullable {
			constraints = append(constraints, "NOT NULL")
		}

		if sfield.Default != "" {
			constraints = append(constraints, "DEFAULT " + sfield.Default)
		}

		if i < numFields-1 {
			field.Separator = ","
		}

		if len(constraints) > 0 {
			field.Constraints = " " + strings.Join(constraints, " ")
		}

		data.Fields = append(data.Fields, field)
	}

	err := tmplDatabase.Execute(&output, data)

	if err != nil {
		return "", err
	}

	return output.String(), nil
}
