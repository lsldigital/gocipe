package crud

import (
	"bytes"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
)

var tmplSave, _ = template.New("GenerateSave").Parse(`
// Save either inserts or updates a {{.Name}} record based on whether or not id is nil
func (entity *{{.Name}}) Save() error {
	if entity.ID == nil {
		return entity.Insert()
	}
	return entity.Update()
}
`)

var tmplSaveHook, _ = template.New("GenerateSaveHook").Parse(`
{{if .PreExecHook }}
func savePreExecHook(entity *User) error {
	return nil
}
{{end}}
{{if .PostExecHook }}
func savePostExecHook(entity *User) error {
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

// GenerateSaveHook will generate 2 functions: savePreExecHook() and savePostExecHook()
func GenerateSaveHook(PreExecHook bool, PostExecHook bool) (string, error) {
	var output bytes.Buffer

	data := new(struct {
		PreExecHook  bool
		PostExecHook bool
	})

	data.PreExecHook = PreExecHook
	data.PostExecHook = PostExecHook

	err := tmplSaveHook.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
