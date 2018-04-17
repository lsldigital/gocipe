package crud

import (
	"bytes"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
)

var tmplSave, _ = template.New("GenerateSave").Parse(`
// Save either inserts or updates a {{.Name}} record based on whether or not id is nil
func (entity *{{.Name}}) Save(tx *sql.Tx) (*sql.Tx, error) {
	if entity.ID == nil {
		return tx, entity.Insert()
	}
	return tx, entity.Update()
}
`)

var tmplSaveHook, _ = template.New("GenerateSaveHook").Parse(`
{{if .PreExecHook }}
func crudSavePreExecHook(entity *{{.Name}}) error {
	return nil
}
{{end}}
{{if .PostExecHook }}
func crudSavePostExecHook(entity *{{.Name}}) error {
	return nil
}
{{end}}
`)

//GenerateSave return code to save entity in database
func GenerateSave(structInfo generators.StructureInfo) (string, error) {
	var output bytes.Buffer

	err := tmplSave.Execute(&output, struct{ Name string }{structInfo.Name})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

// GenerateSaveHook will generate 2 functions: crudSavePreExecHook() and crudSavePostExecHook()
func GenerateSaveHook(structInfo generators.StructureInfo, PreExecHook bool, PostExecHook bool) (string, error) {
	var output bytes.Buffer

	data := new(struct {
		Name         string
		PreExecHook  bool
		PostExecHook bool
	})

	data.Name = structInfo.Name
	data.PreExecHook = PreExecHook
	data.PostExecHook = PostExecHook

	err := tmplSaveHook.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
