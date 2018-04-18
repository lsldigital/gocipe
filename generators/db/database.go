package db

import (
	"bytes"
	"strings"
	"text/template"

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

{{.ManyMany}}
`)

var tmplDatabaseManyMany, _ = template.New("GenerateDatabaseManyMany").Parse(`
DROP TABLE IF EXISTS {{.TableName}};

CREATE TABLE {{.TableName}} (
	"{{.ThisID}}" INT NOT NULL,
	"{{.ThatID}}" INT NOT NULL
);
`)

//GenerateDatabase return code for database sql
func GenerateDatabase(structInfo generators.StructureInfo) (string, error) {
	var (
		output    bytes.Buffer
		numFields = len(structInfo.Fields)
		manymany  []string
	)
	data := new(struct {
		TableName string
		Fields    []tableField
		ManyMany  string
	})

	data.TableName = structInfo.TableName

	for i, sfield := range structInfo.Fields {
		var (
			field       tableField
			constraints []string
		)

		if sfield.ManyMany != nil {
			if strings.Compare(sfield.ManyMany.ThisID, sfield.ManyMany.ThatID) == -1 {
				manymany = append(manymany, manymanyDatabase(sfield))
			}
			continue
		}

		field.Name = sfield.Name
		field.Type = sfield.DBType

		if field.Name == "id" {
			constraints = append(constraints, "PRIMARY KEY")
		}

		if !sfield.Nullable {
			constraints = append(constraints, "NOT NULL")
		}

		if sfield.Default != "" {
			constraints = append(constraints, "DEFAULT "+sfield.Default)
		}

		if i < numFields-1 {
			field.Separator = ","
		}

		if len(constraints) > 0 {
			field.Constraints = " " + strings.Join(constraints, " ")
		}

		data.Fields = append(data.Fields, field)
	}

	if len(manymany) > 0 {
		data.ManyMany = strings.Join(manymany, "\n")
	}

	err := tmplDatabase.Execute(&output, data)

	if err != nil {
		return "", err
	}

	return output.String(), nil
}

func manymanyDatabase(field generators.FieldInfo) string {
	var output bytes.Buffer
	data := new(struct {
		TableName string
		ThisID    string
		ThatID    string
		ManyMany  string
	})

	data.TableName = field.ManyMany.PivotTable
	data.ThisID = field.ManyMany.ThisID
	data.ThatID = field.ManyMany.ThatID

	err := tmplDatabaseManyMany.Execute(&output, data)

	if err != nil {
		return ""
	}

	return output.String()
}
