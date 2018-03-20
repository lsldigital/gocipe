package crud

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
)

var tmplList, _ = template.New("GenerateList").Parse(`
// List returns a slice containing {{.Name}} records
func List() ([]*{{.Name}}, error) {
	var list []*{{.Name}}

	rows, err := db.Query("SELECT {{.SQLFields}} FROM {{.TableName}} ORDER BY id ASC")

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		entity := new({{.Name}})
		err := rows.Scan({{.StructFields}})
		if err != nil {
			return nil, err
		}

		list = append(list, entity)
	}

	return list, nil
}
`)

//GenerateList returns code to return a list of entities from database
func GenerateList(structInfo generators.StructureInfo) (string, error) {
	var output bytes.Buffer
	data := new(struct {
		Name         string
		TableName    string
		SQLFields    string
		StructFields string
	})

	data.Name = structInfo.Name
	data.TableName = "`" + structInfo.TableName + "`"
	data.SQLFields = ""
	data.StructFields = ""

	for _, field := range structInfo.Fields {
		data.SQLFields += strings.ToLower(field.Name) + ", "
		data.StructFields += "entity." + field.Name + ", "
	}

	data.SQLFields = strings.TrimSuffix(data.SQLFields, ", ")
	data.StructFields = strings.TrimSuffix(data.StructFields, ", ")

	err := tmplList.Execute(&output, data)

	if err != nil {
		return "", err
	}

	return output.String(), nil
}
