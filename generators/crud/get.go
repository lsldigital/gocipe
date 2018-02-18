package crud

import (
	"bytes"
	"strings"
	"text/template"
	"projects/gocipe/generators"
)

var tmplGet, _ = template.New("GenerateGet").Parse(`
//Get returns a single {{.Name}} from database
func Get(db *sql.DB, id int) (*{{.Name}}, error) {
	var entity = new({{.Name}})

	query := db.QueryRow("SELECT {{.SQLFields}} FROM {{.TableName}} WHERE id = ? LIMIT 1", id)
	err := query.Scan({{.StructFields}})

	return entity, err
}
`)

//GenerateGet generates code to get an entity from database
func GenerateGet(structInfo generators.StructureInfo) (string, error) {
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
		data.StructFields += "entity." + field.Name + ", "
	}

	data.SQLFields = strings.TrimSuffix(data.SQLFields, ", ")
	data.StructFields = strings.TrimSuffix(data.StructFields, ", ")

	err := tmplGet.Execute(&output, data)

	if err != nil {
		return "", err
	}

	return output.String(), nil
}
