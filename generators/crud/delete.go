package crud

import (
	"bytes"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
)

var tmplDelete, _ = template.New("GenerateDelete").Parse(`
// Delete deletes a {{.Name}} record from database by id primary key
func Delete(id int64, tx *sql.Tx, autocommit bool) error {
	var err error

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return err
		}
	}

	stmt, err := tx.Prepare("DELETE FROM {{.TableName}} WHERE id = $1")
	if err != nil {
		return err
	}
	{{if .PreExecHook}}
	if err := crudPreDelete(id, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing crudPreDelete() in Delete(%d) for entity '{{.Name}}': %s", id, err)
	}
	{{end}}
	_, err = stmt.Exec(id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing transaction statement in Delete(%d) for entity '{{.Name}}': %s", id, err)
	}
	{{if .PostExecHook}}
	if err := crudPostDelete(id, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("Error executing crudPostDelete() in Delete(%d) for entity '{{.Name}}': %s", id, err)
	}
	{{end}}
	if autocommit {
		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("error committing transaction in Delete(%d) for '{{.Name}}': %s", id, err)
		}
	}

	return err
}

// Delete deletes a {{.Name}} record from database and sets id to nil
func (entity *{{.Name}}) Delete(tx *sql.Tx, autocommit bool) error {
	var err error

	id := *entity.ID

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return err
		}
	}

	stmt, err := tx.Prepare("DELETE FROM {{.TableName}} WHERE id = $1")
	if err != nil {
		return err
	}
	{{if .PreExecHook}}
	if err := crudPreDelete(id, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing crudPreDelete() in {{.Name}}.Delete() for ID = %d : %s", id, err)
	}
	{{end}}
	_, err = stmt.Exec(id)
	if err == nil {
		entity.ID = nil
	} else {
		tx.Rollback()
		return fmt.Errorf("error executing transaction statement in {{.Name}}.Delete() for ID = %d : %s", id, err)
	}
	{{if .PostExecHook}}
	if err = crudPostDelete(id, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing crudPostDelete() in {{.Name}}.Delete() for ID = %d : %s", id, err)
	}
	{{end}}
	if autocommit {
		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("error committing transaction in {{.Name}}.Delete() for ID = %d : %s", id, err)
		}
	}

	return err
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
