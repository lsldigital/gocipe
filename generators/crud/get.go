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
    if entity, e := getPreExecHook(id); e != nil {
		fmt.Printf("Error executing getPreExecHook() in Get(" + strconv.FormatInt(id, 10) + ") for entity '{{.Name}}': %s", e.Error())
        return entity, e
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
        if entity, e := getPostExecHook(entity); e != nil {
			fmt.Printf("Error executing getPostExecHook() in Get(" + strconv.FormatInt(id, 10) + ") for entity '{{.Name}}': %s", e.Error())
			return entity, e
        }
        {{end}}
		return entity, nil
	}

	return nil, nil
}
`)

var tmplGetHook, _ = template.New("GenerateGetHook").Parse(`
{{if .PreExecHook }}
func getPreExecHook(id int64) (*User, error) {
	return New(), nil
}
{{end}}
{{if .PostExecHook }}
func getPostExecHook(entity *User) (*User, error) {
	return entity, nil
}
{{end}}
`)

//GenerateGet generates code to get an entity from database
func GenerateGet(structInfo generators.StructureInfo, PreExecHook bool, PostExecHook bool) (string, error) {
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
	data.PreExecHook = PreExecHook
	data.PostExecHook = PostExecHook

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

// GenerateGetHook will generate 2 functions: getPreExecHook() and getPostExecHook()
func GenerateGetHook(PreExecHook bool, PostExecHook bool) (string, error) {
	var output bytes.Buffer

	data := new(struct {
		PreExecHook  bool
		PostExecHook bool
	})

	data.PreExecHook = PreExecHook
	data.PostExecHook = PostExecHook

	err := tmplGetHook.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
