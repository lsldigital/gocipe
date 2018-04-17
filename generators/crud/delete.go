package crud

import (
	"bytes"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
)

var tmplDelete, _ = template.New("GenerateDelete").Parse(`
// Delete deletes a {{.Name}} record from database by id primary key
func Delete(id int64, tx *sql.Tx) (*sql.Tx, error) {
	var (
		err      error
		txWasNil bool
	)

	if tx == nil {
		txWasNil = true
		tx, err = db.Begin()
		if err != nil {
			return nil, err
		}
	}

	stmt, err := tx.Prepare("DELETE FROM {{.TableName}} WHERE id = $1")
	if err != nil {
		return nil, err
	}
	{{if .PreExecHook}}
	if err := crudPreDelete(id, tx); err != nil {
		fmt.Printf("Error executing crudPreDelete() in Delete(%d) for entity '{{.Name}}': %s", id, err.Error())
		tx.Rollback()
		return nil, err
	}
	{{end}}
	_, err = stmt.Exec(id)
	{{if .PostExecHook}}
	if err := crudPostDelete(id, tx); err != nil {
		fmt.Printf("Error executing crudPostDelete() in Delete(%d) for entity '{{.Name}}': %s", id, err.Error())
		tx.Rollback()
		return nil, err
	}
	{{end}}
	if txWasNil {
		tx.Commit()
	}

	return tx, err
}

// Delete deletes a {{.Name}} record from database and sets id to nil
func (entity *{{.Name}}) Delete(tx *sql.Tx) (*sql.Tx, error) {
	var (
		err      error
		txWasNil bool
	)

	id := *entity.ID

	if tx == nil {
		txWasNil = true
		tx, err = db.Begin()
		if err != nil {
			return nil, err
		}
	}

	stmt, err := tx.Prepare("DELETE FROM {{.TableName}} WHERE id = $1")
	if err != nil {
		return nil, err
	}
	{{if .PreExecHook}}
	if err := crudPreDelete(id, tx); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("Error executing crudPreDelete() in User.Delete() for ID = %d : %s", id, err.Error())
	}
	{{end}}
	_, err = stmt.Exec(id)
	if err != nil {
		entity.ID = nil
	}
	{{if .PostExecHook}}
	if err = crudPostDelete(id, tx); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("Error executing crudPostDelete() in User.Delete() for ID = %d : %s", id, err.Error())
	}
	{{end}}
	if txWasNil {
		tx.Commit()
	}

	return tx, err
}
`)

var tmplDeleteHook, _ = template.New("GenerateDeleteHook").Parse(`
{{if .PreExecHook }}
func crudPreDelete(id int64, tx *sql.Tx) error {
	return nil
}
{{end}}
{{if .PostExecHook }}
func crudPostDelete(id int64, tx *sql.Tx) error {
	return nil
}
{{end}}
`)

//GenerateDelete will generate a function to delete entity from database
func GenerateDelete(structInfo generators.StructureInfo, PreExecHook bool, PostExecHook bool) (string, error) {
	var output bytes.Buffer
	data := new(struct {
		Name         string
		TableName    string
		PreExecHook  bool
		PostExecHook bool
	})

	data.Name = structInfo.Name
	data.TableName = structInfo.TableName
	data.PreExecHook = PreExecHook
	data.PostExecHook = PostExecHook

	err := tmplDelete.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

// GenerateDeleteHook will generate 2 functions: crudPreDelete() and crudPostDelete()
func GenerateDeleteHook(PreExecHook bool, PostExecHook bool) (string, error) {
	var output bytes.Buffer

	data := new(struct {
		PreExecHook  bool
		PostExecHook bool
	})

	data.PreExecHook = PreExecHook
	data.PostExecHook = PostExecHook

	err := tmplDeleteHook.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
