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

//GenerateSave return code to save entity in database
func GenerateSave(structInfo generators.StructureInfo) (string, error) {
	var output bytes.Buffer
	err := tmplSave.Execute(&output, struct{ Name string }{structInfo.Name})

	if err != nil {
		return "", err
	}

	return output.String(), nil
}
