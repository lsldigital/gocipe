package crud

import (
	"bytes"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
)

var tmplDelete, _ = template.New("GenerateDelete").Parse(`
// Delete deletes a {{.Name}} record from database by id primary key
func Delete(id int64) error {
    {{if .PreExecHook}}
    if e := deletePreExecHook(id); e != nil {
		fmt.Printf("Error executing deletePreExecHook() in Delete(" + strconv.FormatInt(id, 10) + ") for entity '{{.Name}}': %s", e.Error())
        return e
    }
    {{end}}
	_, err := db.Exec("DELETE FROM {{.TableName}} WHERE id = $1", id)
    {{if .PostExecHook}}
    if e := deletePostExecHook(id); e != nil {
		fmt.Printf("Error executing deletePostExecHook() in Delete(" + strconv.FormatInt(id, 10) + ") for entity '{{.Name}}': %s", e.Error())
        return e
    }
    {{end}}
	return err
}

// Delete deletes a {{.Name}} record from database and sets id to nil
func (entity *{{.Name}}) Delete() error {
    {{if .PreExecHook}}
    if e := deletePreExecHook(*entity.ID); e != nil {
		fmt.Printf("Error executing deletePreExecHook() in {{.Name}}.Delete() for ID = " + strconv.FormatInt(*entity.ID, 10) + ": %s", e.Error())
        return e
    }
    {{end}}
	_, err := db.Exec("DELETE FROM {{.TableName}} WHERE id = $1", entity.ID)
	{{if .PostExecHook}}id := *entity.ID{{end}}
	if err != nil {
		entity.ID = nil
	}
    {{if .PostExecHook}}
    if e := deletePostExecHook(id); e != nil {
		fmt.Printf("Error executing deletePostExecHook() in {{.Name}}.Delete() for ID = " + strconv.FormatInt(*entity.ID, 10) + ": %s", e.Error())
        return e
    }
    {{end}}
	return err
}
`)

var tmplDeleteHook, _ = template.New("GenerateDeleteHook").Parse(`
{{if .PreExecHook }}
func deletePreExecHook(id int64) error {
	return nil
}
{{end}}
{{if .PostExecHook }}
func deletePostExecHook(id int64) error {
	return nil
}
{{end}}
`)

//GenerateDelete will generate a function to delete entity from database
func GenerateDelete(structInfo generators.StructureInfo, PreExecHook bool, PostExecHook bool) (string, error) {
	var output bytes.Buffer
	data := new(struct {
		Name         string
		TableName    string
		PreExecHook  bool
		PostExecHook bool
	})

	data.Name = structInfo.Name
	data.TableName = structInfo.TableName
	data.PreExecHook = PreExecHook
	data.PostExecHook = PostExecHook

	err := tmplDelete.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

// GenerateDeleteHook will generate 2 functions: deletePreExecHook() and deletePostExecHook()
func GenerateDeleteHook(PreExecHook bool, PostExecHook bool) (string, error) {
	var output bytes.Buffer

	data := new(struct {
		PreExecHook  bool
		PostExecHook bool
	})

	data.PreExecHook = PreExecHook
	data.PostExecHook = PostExecHook

	err := tmplDeleteHook.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
