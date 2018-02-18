package crud

import (
	"bytes"
	"strings"
	"text/template"
	"projects/gocipe/generators"
)

var tmplList, _ = template.New("GenerateList").Parse(`
//List returns all {{.Name}} entities stored in database
func List(db *sql.DB) ([]{{.Name}}, error) {
	var (
		entity   *{{.Name}}
		entities []{{.Name}}
	)

	query, err := db.Query("SELECT {{.SQLFields}} FROM {{.TableName}}")
	if err != nil {
		return entities, err
	}
	defer query.Close()

	entities = make([]{{.Name}}, 10)

	for query.Next() {
		entity = new({{.Name}})
		query.Scan({{.StructFields}})
		entities = append(entities, *entity)
	}

	return entities, nil
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
	data.TableName = "`" + strings.ToLower(structInfo.Name) + "s`"
	data.SQLFields = ""
	data.StructFields = ""

	for _, field := range structInfo.Fields {
		data.SQLFields += field.Name + ", "
		data.StructFields += "entity" + field.Name + ", "
	}

	data.SQLFields = strings.TrimSuffix(data.SQLFields, ", ")
	data.StructFields = strings.TrimSuffix(data.StructFields, ", ")

	err := tmplList.Execute(&output, data)

	if err != nil {
		return "", err
	}

	return output.String(), nil
}
