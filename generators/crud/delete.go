package crud

import (
	"bytes"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
)

var tmplDelete, _ = template.New("GenerateDelete").Parse(`
// Delete deletes a {{.Name}} record from database by id primary key
func Delete(id int64) error {
	_, err := db.Exec("DELETE FROM {{.TableName}} WHERE id = $1", id)
	return err
}

// Delete deletes a {{.Name}} record from database and sets id to nil
func (entity *{{.Name}}) Delete() error {
	_, err := db.Exec("DELETE FROM {{.TableName}} WHERE id = $1", entity.ID)
	if err != nil {
		entity.ID = nil
	}
	return err
}
`)

//GenerateDelete will generate a function to delete entity from database
func GenerateDelete(structInfo generators.StructureInfo) (string, error) {
	var output bytes.Buffer
	data := new(struct {
		Name      string
		TableName string
	})

	data.Name = structInfo.Name
	data.TableName = structInfo.TableName

	err := tmplDelete.Execute(&output, data)

	if err != nil {
		return "", err
	}

	return output.String(), nil
}
