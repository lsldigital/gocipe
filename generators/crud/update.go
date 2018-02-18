package crud

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
)

var tmplUpdate, _ = template.New("GenerateUpdate").Parse(`
//Update Will execute an SQLUpdate Statement in the database. Prefer using Save instead of Update directly.
func (entity *{{.Name}}) Update(db *sql.DB) error {
	stmt, err := db.Prepare("UPDATE {{.TableName}} SET {{.SQLFields}} WHERE id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec({{.StructFields}})
	if err != nil {
		return err
	}
}
`)

//GenerateUpdate returns code to update an entity in database
func GenerateUpdate(structInfo generators.StructureInfo) (string, error) {
	var output bytes.Buffer
	data := new(struct {
		Name         string
		TableName    string
		SQLFields    string
		StructFields string
	})

	data.Name = structInfo.Name
	data.TableName = "`" + strings.ToLower(structInfo.Name) + "s`"
	data.SQLFields = ""
	data.StructFields = ""

	for _, field := range structInfo.Fields {
		if field.Name == "id" {
			continue
		}

		data.SQLFields += field.Name + " = ?, "
		data.StructFields += "entity." + field.Name + ", "
	}
	data.SQLFields = strings.TrimSuffix(data.SQLFields, ", ")
	data.StructFields += "entity.id"

	err := tmplUpdate.Execute(&output, data)

	if err != nil {
		return "", err
	}

	return output.String(), nil
}
