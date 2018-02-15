package simple

import (
	"bytes"
	"strings"
	"text/template"
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
func GenerateList(name string, fields []string) (string, error) {
	var output bytes.Buffer
	data := new(struct {
		Name         string
		TableName    string
		SQLFields    string
		StructFields string
	})

	data.Name = name
	data.TableName = "`" + strings.ToLower(name) + "s`"
	data.SQLFields = strings.Join(fields, ", ")
	data.StructFields = "entity." + strings.Join(fields, ", entity.")

	err := tmplList.Execute(&output, data)

	if err != nil {
		return "", err
	}

	return output.String(), nil
}
