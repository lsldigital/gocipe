package simple

import (
	"bytes"
	"strings"
	"text/template"
)

var tmplDelete, _ = template.New("GenerateDelete").Parse(`
//Delete delete single {{.Name}} entity from database
func Delete(db *sql.DB, id int) error {
	stmt, err := db.Prepare("DELETE FROM {{.TableName}} WHERE id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
`)

//GenerateDelete will generate a function to delete entity from database
func GenerateDelete(name string) (string, error) {
	var output bytes.Buffer
	data := new(struct {
		Name      string
		TableName string
	})

	data.Name = name
	data.TableName = "`" + strings.ToLower(name) + "s`"

	err := tmplDelete.Execute(&output, data)

	if err != nil {
		return "", err
	}

	return output.String(), nil
}
