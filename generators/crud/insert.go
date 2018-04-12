package crud

import (
	"bytes"
	"strconv"
	"strings"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
)

var tmplInsert, _ = template.New("GenerateInsert").Parse(`
// Insert performs an SQL insert for {{.Name}} record and update instance with inserted id.
// Prefer using Save rather than Insert directly.
func (entity *{{.Name}}) Insert() error {
	var (
		id  int64
		err error
	)
	{{if .PreExecHook }}
    if e := savePreExecHook(entity); e != nil {
        fmt.Printf("Error executing savePreExecHook() in {{.Name}}.Insert(): %s", e.Error())
        return e
	}
    {{end}}
	err = db.QueryRow("INSERT INTO {{.TableName}} ({{.SQLFields}}) VALUES ({{.SQLPlaceholders}}) RETURNING id", {{.StructFields}}).Scan(&id)

	if err == nil {
		entity.ID = &id
	}
	{{if .PostExecHook }}
	if e := savePostExecHook(entity); e != nil {
		fmt.Printf("Error executing savePostExecHook() in {{.Name}}.Insert(): %s", e.Error())
		return e
	}
	{{end}}
	return err
}
`)

//GenerateInsert generate function to insert an entity in database
func GenerateInsert(structInfo generators.StructureInfo, PreExecHook bool, PostExecHook bool) (string, error) {
	var output bytes.Buffer
	data := new(struct {
		Name            string
		TableName       string
		SQLFields       string
		SQLPlaceholders string
		StructFields    string
		PreExecHook     bool
		PostExecHook    bool
	})

	data.Name = structInfo.Name
	data.TableName = structInfo.TableName
	data.SQLFields = ""
	data.SQLPlaceholders = ""
	data.StructFields = ""
	data.PreExecHook = PreExecHook
	data.PostExecHook = PostExecHook

	for i, field := range structInfo.Fields {
		if field.Name == "id" {
			continue
		}

		data.SQLFields += field.Name + ", "
		data.SQLPlaceholders += "$" + strconv.Itoa(i) + ", "
		data.StructFields += "*entity." + field.Property + ", "
	}

	data.SQLFields = strings.TrimSuffix(data.SQLFields, ", ")
	data.SQLPlaceholders = strings.TrimSuffix(data.SQLPlaceholders, ", ")
	data.StructFields = strings.TrimSuffix(data.StructFields, ", ")

	err := tmplInsert.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
