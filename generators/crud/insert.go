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
func (entity *{{.Name}}) Insert(tx *sql.Tx, autocommit bool) error {
	var (
		id  int64
		err error
	)

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return err
		}
	}
	
	stmt, err := tx.Prepare("INSERT INTO {{.TableName}} ({{.SQLFields}}) VALUES ({{.SQLPlaceholders}}) RETURNING id")
	if err != nil {
		return err
	}
	{{if .PreExecHook }}
    if err := crudPreSave(entity, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing crudPreSave() in {{.Name}}.Insert(): %s", err)
	}
    {{end}}
	err = stmt.QueryRow({{.StructFields}}).Scan(&id)
	if err == nil {
		entity.ID = &id
	} else {
		tx.Rollback()
		return fmt.Errorf("error executing transaction statement in {{.Name}}: %s", err)
	}
	{{if .PostExecHook }}
	if err := crudPostSave(entity, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing crudPostSave() in {{.Name}}.Insert(): %s", err)
	}
	{{end}}
	if autocommit {
		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("error committing transaction in {{.Name}}.Insert(): %s", err)
		}
	}

	return nil
}
`)

//GenerateInsert generate function to insert an entity in database
func GenerateInsert(structInfo generators.StructureInfo, PreExecHook bool, PostExecHook bool) (string, error) {
	var output bytes.Buffer
	data := new(struct {
		Name            string
		TableName       string
		SQLFields       string
		SQLPlaceholders string
		StructFields    string
		PreExecHook     bool
		PostExecHook    bool
	})

	data.Name = structInfo.Name
	data.TableName = structInfo.TableName
	data.SQLFields = ""
	data.SQLPlaceholders = ""
	data.StructFields = ""
	data.PreExecHook = PreExecHook
	data.PostExecHook = PostExecHook

	for i, field := range structInfo.Fields {
		if field.Name == "id" {
			continue
		}

		data.SQLFields += field.Name + ", "
		data.SQLPlaceholders += "$" + strconv.Itoa(i) + ", "
		data.StructFields += "*entity." + field.Property + ", "
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
