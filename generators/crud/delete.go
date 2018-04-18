package crud

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
)

var tmplDelete, _ = template.New("GenerateDelete").Parse(`
// Delete deletes a {{.Name}} record from database by id primary key
func Delete(id int64, tx *sql.Tx, autocommit bool) error {
	var (
		err error
		{{.ManyManyVars}}
	)

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
	{{.ManyManyFunc}}
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
	var (
		err error
		{{.ManyManyVars}}
	)

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
	{{.ManyManyMethod}}
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

var tmplManyManyDelete, _ = template.New("ManyManyDelete").Parse(`
	stmtMmany, err = tx.Prepare("DELETE FROM {{.PivotTable}} WHERE {{.ThisID}} = $1")

	if err != nil {
		return fmt.Errorf("error preparing transaction statement in ManyManyDelete(%d) {{.ThisID}}-{{.ThatID}} for table '{{.PivotTable}}': %s", {{.ID}}, err)
	}

	_, err = stmtMmany.Exec({{.ID}})
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing transaction statement in ManyManyDelete(%d) {{.ThisID}}-{{.ThatID}} for table '{{.PivotTable}}': %s", {{.ID}}, err)
	}
`)

//GenerateDelete will generate a function to delete entity from database
func GenerateDelete(structInfo generators.StructureInfo, preExecHook bool, postExecHook bool) (string, error) {
	var (
		output         bytes.Buffer
		manyManyFunc   []string
		manyManyMethod []string
	)
	data := new(struct {
		Name           string
		TableName      string
		ManyManyFunc   string
		ManyManyMethod string
		ManyManyVars   string
		PreExecHook    bool
		PostExecHook   bool
	})

	data.Name = structInfo.Name
	data.TableName = structInfo.TableName
	data.PreExecHook = preExecHook
	data.PostExecHook = postExecHook

	for _, field := range structInfo.Fields {
		if field.ManyMany != nil {
			manyManyFunc = append(manyManyFunc, deleteManyMany("id", field))
			manyManyMethod = append(manyManyMethod, deleteManyMany("*entity.ID", field))
		}
	}

	if len(manyManyFunc) != 0 {
		data.ManyManyFunc = strings.Join(manyManyFunc, "\n") + "\n"
		data.ManyManyVars = "stmtMmany *sql.Stmt"
	}

	if len(manyManyMethod) != 0 {
		data.ManyManyMethod = strings.Join(manyManyMethod, "\n") + "\n"
	}

	err := tmplDelete.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

// GenerateDeleteHook will generate 2 functions: crudPreDelete() and crudPostDelete()
func GenerateDeleteHook(preExecHook bool, postExecHook bool) (string, error) {
	var output bytes.Buffer

	data := new(struct {
		PreExecHook  bool
		PostExecHook bool
	})

	data.PreExecHook = preExecHook
	data.PostExecHook = postExecHook

	err := tmplDeleteHook.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

func deleteManyMany(idfield string, field generators.FieldInfo) string {
	var output bytes.Buffer
	data := new(struct {
		ID         string
		Name       string
		PivotTable string
		ThatID     string
		ThisID     string
	})

	data.ID = idfield
	data.Name = field.Name
	data.PivotTable = field.ManyMany.PivotTable
	data.ThatID = field.ManyMany.ThatID
	data.ThisID = field.ManyMany.ThisID

	err := tmplManyManyDelete.Execute(&output, data)
	if err != nil {
		return ""
	}

	return output.String()
}
