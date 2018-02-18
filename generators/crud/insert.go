package crud

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
)

var tmplInsert, _ = template.New("GenerateInsert").Parse(`
//Insert Will execute an SQLInsert Statement in the database. Prefer using Save instead of Insert directly.
func (entity *{{.Name}}) Insert(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO {{.TableName}} ({{.SQLFields}}) VALUES ({{.SQLPlaceholders}})")

	if err != nil {
		return err
	}

	result, err := stmt.Exec({{.StructFields}})
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	entity.id = id
	return nil
}
`)

//GenerateInsert generate function to insert an entity in database
func GenerateInsert(structInfo generators.StructureInfo) (string, error) {
	var output bytes.Buffer
	data := new(struct {
		Name            string
		TableName       string
		SQLFields       string
		SQLPlaceholders string
		StructFields    string
	})

	data.Name = structInfo.Name
	data.TableName = "`" + strings.ToLower(structInfo.Name) + "s`"
	data.SQLFields = ""
	data.SQLPlaceholders = ""
	data.StructFields = ""

	for _, field := range structInfo.Fields {
		if field.Name == "id" {
			continue
		}

		data.SQLFields += field.Name + ", "
		data.SQLPlaceholders += "?, "
		data.StructFields += "entity." + field.Name + ", "
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
