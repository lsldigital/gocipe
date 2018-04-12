package crud

import (
	"bytes"
	"strconv"
	"strings"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
)

var tmplUpdate, _ = template.New("GenerateUpdate").Parse(`
//Update Will execute an SQLUpdate Statement for {{.Name}} in the database. Prefer using Save instead of Update directly.
func (entity *{{.Name}}) Update() error {
	{{if .PreExecHook }}
    if e := crudSavePreExecHook(entity); e != nil {
        fmt.Printf("Error executing crudSavePreExecHook() in {{.Name}}.Update(): %s", e.Error())
        return e
	}
    {{end}}
	_, err := db.Exec("UPDATE {{.TableName}} SET {{.SQLFields}} WHERE id = $1", {{.StructFields}})
	{{if .PostExecHook }}
	if e := crudSavePostExecHook(entity); e != nil {
		fmt.Printf("Error executing crudSavePostExecHook() in {{.Name}}.Update(): %s", e.Error())
		return e
	}
	{{end}}
	return err
}
`)

//GenerateUpdate returns code to update an entity in database
func GenerateUpdate(structInfo generators.StructureInfo, PreExecHook bool, PostExecHook bool) (string, error) {
	var (
		output bytes.Buffer
		index  = 2
	)
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
	data.StructFields = "entity.ID, "
	data.PreExecHook = PreExecHook
	data.PostExecHook = PostExecHook

	for _, field := range structInfo.Fields {
		if field.Name == "ID" {
			continue
		}

		data.SQLFields += field.Name + " = $" + strconv.Itoa(index) + ", "
		data.StructFields += "*entity." + field.Property + ", "
		index++
	}
	data.SQLFields = strings.TrimSuffix(data.SQLFields, ", ")
	data.StructFields = strings.TrimSuffix(data.StructFields, ", ")

	err := tmplUpdate.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
