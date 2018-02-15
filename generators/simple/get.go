package simple

import (
	"bytes"
	"strings"
	"text/template"
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
func GenerateGet(name string, fields []string) (string, error) {
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

	err := tmplGet.Execute(&output, data)

	if err != nil {
		return "", err
	}

	return output.String(), nil
}
