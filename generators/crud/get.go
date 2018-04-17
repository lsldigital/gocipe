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
	var entity = New()
	{{if .PreExecHook }}
    if err := crudPreGet(id); err != nil {
		return nil, fmt.Errorf("error executing crudPreGet() in Get(%d) for entity '{{.Name}}': %s", id, err)
	}
    {{end}}
	rows, err := db.Query("SELECT {{.SQLFields}} FROM {{.TableName}} WHERE id = $1 ORDER BY id ASC", id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	if rows.Next() {
		err := rows.Scan({{.StructFields}})
		if err != nil {
			return nil, err
		}
        {{if .PostExecHook }}
		if err = crudPostGet(entity); err != nil {
			return nil, fmt.Errorf("error executing crudPostGet() in Get(%d) for entity '{{.Name}}': %s", id, err)
		}
        {{end}}
		return entity, nil
	}

	return nil, nil
}
`)

var tmplGetHook, _ = template.New("GenerateGetHook").Parse(`
{{if .PreExecHook }}
func crudPreGet(id int64) error {
	return nil
}
{{end}}
{{if .PostExecHook }}
func crudPostGet(entity *{{.Name}}) error {
	return nil
}
{{end}}
`)

//GenerateGet generates code to get an entity from database
func GenerateGet(structInfo generators.StructureInfo, preExecHook bool, postExecHook bool) (string, error) {
	var output bytes.Buffer
	data := new(struct {
		Name         string
		TableName    string
		SQLFields    string
		StructFields string
		PreExecHook  bool
		PostExecHook bool
	})

	data.Name = structInfo.Name
	data.TableName = structInfo.TableName
	data.SQLFields = ""
	data.StructFields = ""
	data.PreExecHook = preExecHook
	data.PostExecHook = postExecHook

	for _, field := range structInfo.Fields {
		data.SQLFields += field.Name + ", "
		data.StructFields += "entity." + field.Property + ", "
	}

	data.SQLFields = strings.TrimSuffix(data.SQLFields, ", ")
	data.StructFields = strings.TrimSuffix(data.StructFields, ", ")

	err := tmplGet.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

// GenerateGetHook will generate 2 functions: crudPreGet() and crudPostGet()
func GenerateGetHook(structInfo generators.StructureInfo, preExecHook bool, postExecHook bool) (string, error) {
	var output bytes.Buffer

	data := new(struct {
		Name         string
		PreExecHook  bool
		PostExecHook bool
	})

	data.Name = structInfo.Name
	data.PreExecHook = preExecHook
	data.PostExecHook = postExecHook

	err := tmplGetHook.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
