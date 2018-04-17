package crud

import (
	"bytes"
	"strconv"
	"strings"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
)

var tmplUpdate, _ = template.New("GenerateUpdate").Parse(`
//Update Will execute an SQLUpdate Statement for {{.Name}} in the database. Prefer using Save instead of Update directly.
func (entity *{{.Name}}) Update(tx *sql.Tx, autocommit bool) error {
	var err error

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return err
		}
	}

	stmt, err := tx.Prepare("UPDATE users SET id = $2, auth_code = $3, alias = $4, name = $5, callback = $6, status = $7 WHERE id = $1")
	if err != nil {
		return err
	}

	{{if .PreExecHook }}
    if err := crudPreSave(entity, tx); err != nil {
		tx.Rollback()
        return fmt.Errorf("error executing crudPreSave() in {{.Name}}.Update(): %s", err)
	}
    {{end}}
	_, err = stmt.Exec("UPDATE {{.TableName}} SET {{.SQLFields}} WHERE id = $1", {{.StructFields}})
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing transaction statement in {{.Name}}.Update(): %s", err)
	}

	{{if .PostExecHook }}
	if err := crudPostSave(entity, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing crudPostSave() in {{.Name}}.Update(): %s", err)
	}
	{{end}}
	if autocommit {
		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("error committing transaction in {{.Name}}.Update(): %s", err)
		}
	}

	return nil
}
`)

//GenerateUpdate returns code to update an entity in database
func GenerateUpdate(structInfo generators.StructureInfo, PreExecHook bool, PostExecHook bool) (string, error) {
	var (
		output bytes.Buffer
		index  = 2
	)
	data := new(struct {
		Name         string
		TableName    string
		SQLFields    string
		StructFields string
		PreExecHook  bool
		PostExecHook bool
	})

	data.Name = structInfo.Name
	data.TableName = structInfo.TableName
	data.SQLFields = ""
	data.StructFields = "entity.ID, "
	data.PreExecHook = PreExecHook
	data.PostExecHook = PostExecHook

	for _, field := range structInfo.Fields {
		if field.Name == "ID" {
			continue
		}

		data.SQLFields += field.Name + " = $" + strconv.Itoa(index) + ", "
		data.StructFields += "*entity." + field.Property + ", "
		index++
	}
	data.SQLFields = strings.TrimSuffix(data.SQLFields, ", ")
	data.StructFields = strings.TrimSuffix(data.StructFields, ", ")

	err := tmplUpdate.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
