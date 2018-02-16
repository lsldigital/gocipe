package crud

import (
	"bytes"
	"strings"
	"text/template"
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
func GenerateInsert(name string, fields []string) (string, error) {
	var output bytes.Buffer
	data := new(struct {
		Name            string
		TableName       string
		SQLFields       string
		SQLPlaceholders string
		StructFields    string
	})

	data.Name = name
	data.TableName = "`" + strings.ToLower(name) + "s`"
	data.SQLFields = ""
	data.SQLPlaceholders = ""
	data.StructFields = ""

	for _, field := range fields {
		if strings.Compare(field, "id") == 0 {
			continue
		}

		data.SQLFields += field + ", "
		data.SQLPlaceholders += "?, "
		data.StructFields += "entity." + field + ", "
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
