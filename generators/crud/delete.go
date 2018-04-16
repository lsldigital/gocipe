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
	if tx, err := crudDeletePreExecHook(id, tx); err != nil {
		fmt.Printf("Error executing deletePreExecHook() in Delete(%d) for entity '{{.Name}}': %s", id, err.Error())
		_ = tx.Rollback()
		return nil, err
	}
	{{end}}
	_, err = stmt.Exec(id)
	{{if .PostExecHook}}
	if tx, err := crudDeletePostExecHook(id, tx); err != nil {
		fmt.Printf("Error executing deletePostExecHook() in Delete(%d) for entity '{{.Name}}': %s", id, err.Error())
		_ = tx.Rollback()
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
	if tx, err := crudDeletePreExecHook(id, tx); err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("Error executing deletePreExecHook() in User.Delete() for ID = %d : %s", id, err.Error())
	}
	{{end}}
	_, err = stmt.Exec(id)
	if err != nil {
		entity.ID = nil
	}
	{{if .PostExecHook}}
	if tx, err = crudDeletePostExecHook(id, tx); err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("Error executing deletePostExecHook() in User.Delete() for ID = %d : %s", id, err.Error())
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
func crudDeletePreExecHook(id int64, tx *sql.Tx) (*sql.Tx, error) {
	return tx, nil
}
{{end}}
{{if .PostExecHook }}
func crudDeletePostExecHook(id int64, tx *sql.Tx) (*sql.Tx, error) {
	return tx, nil
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

// GenerateDeleteHook will generate 2 functions: deletePreExecHook() and deletePostExecHook()
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
