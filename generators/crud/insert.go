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

	result, err := db.Exec("INSERT INTO {{.TableName}} ({{.SQLFields}}) VALUES ({{.SQLPlaceholders}})", {{.StructFields}})

	if err == nil {
		if id, err = result.LastInsertId(); err == nil {
			entity.ID = &id
		}
	}

	return err
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
	data.TableName = structInfo.TableName
	data.SQLFields = ""
	data.SQLPlaceholders = ""
	data.StructFields = ""

	for i, field := range structInfo.Fields {
		if field.Name == "id" {
			continue
		}

		data.SQLFields += strings.ToLower(field.Name) + ", "
		data.SQLPlaceholders += "$" + strconv.Itoa(i+1) + ", "
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
