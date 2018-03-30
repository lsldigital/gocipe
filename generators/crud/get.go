package crud

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
)

var tmplGet, _ = template.New("GenerateGet").Parse(`
//Get returns a single {{.Name}} from database by primary key
func Get(id int64) (*{{.Name}}, error) {
	var entity *{{.Name}}

	rows, err := db.Query("SELECT {{.SQLFields}} FROM {{.TableName}} WHERE id = $1 ORDER BY id ASC", id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	if rows.Next() {
		entity = new({{.Name}})
		err := rows.Scan({{.StructFields}})
		if err != nil {
			return nil, err
		}
		return entity, nil
	}
	
	return nil, nil
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
	data.TableName = "`" + structInfo.TableName + "`"
	data.SQLFields = ""
	data.StructFields = ""

	for _, field := range structInfo.Fields {
		data.SQLFields += strings.ToLower(field.Name) + ", "
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
