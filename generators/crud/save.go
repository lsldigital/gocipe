package crud

import (
	"bytes"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
)

var tmplSave, _ = template.New("GenerateSave").Parse(`
//Save will persist {{.Name}} entity to the database
func (entity *{{.Name}}) Save(db *sql.DB) error {
	if entity.id == 0 {
		error := entity.Insert(db)
	} else {
		error := entity.Update(db)
	}

	if error != nil {
		return error
	}

	return nil
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
