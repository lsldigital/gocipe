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

CREATE TABLE {{.Table}} (
	{{range .Fields}}
	"{{.Name}}" {{.Type}} {{.Constraints}}{{.Separator}}
	{{end}}
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

	data.TableName = "`" + structInfo.TableName + "`"

	for i, sfield := range structInfo.Fields {
		var (
			field       tableField
			constraints []string
		)

		field.Name = sfield.Tags.Get("json")
		field.Type = sfield.Tags.Get("dbtype")

		if field.Name == "id" {
			constraints = append(constraints, "PRIMARY KEY")
		}

		if sfield.Tags.Get("nullable") != "true" {
			constraints = append(constraints, "NOT NULL")
		}

		if def := sfield.Tags.Get("default"); def != "" {
			constraints = append(constraints, "DEFAULT "+def)
		}

		if i < numFields-1 {
			field.Separator = ","
		}

		field.Constraints = strings.Join(constraints, " ")
		data.Fields = append(data.Fields, field)
	}

	err := tmplDatabase.Execute(&output, data)

	if err != nil {
		return "", err
	}

	return output.String(), nil
}
