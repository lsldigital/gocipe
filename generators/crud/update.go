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

	{{.SnippetsBefore}}

	stmt, err := tx.Prepare("UPDATE {{.TableName}} SET {{.SQLFields}} WHERE id = $1")
	if err != nil {
		return err
	}

	{{if .PreExecHook }}
    if err := crudPreSave(entity, tx); err != nil {
		tx.Rollback()
        return fmt.Errorf("error executing crudPreSave() in {{.Name}}.Update(): %s", err)
	}
    {{end}}
	_, err = stmt.Exec({{.StructFields}})
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing transaction statement in {{.Name}}.Update(): %s", err)
	}
	{{.ManyMany}}
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
func GenerateUpdate(structInfo generators.StructureInfo, preExecHook bool, postExecHook bool) (string, error) {
	var (
		output         bytes.Buffer
		index          = 2
		snippetsBefore []string
		manyMany       []string
	)

	data := new(struct {
		Name           string
		TableName      string
		SQLFields      string
		StructFields   string
		SnippetsBefore string
		ManyMany       string
		ManyManyVars   string
		PreExecHook    bool
		PostExecHook   bool
	})

	data.Name = structInfo.Name
	data.TableName = structInfo.TableName
	data.SQLFields = ""
	data.StructFields = "*entity.ID, "
	data.PreExecHook = preExecHook
	data.PostExecHook = postExecHook

	for _, field := range structInfo.Fields {
		if field.Property == "ID" {
			continue
		} else if field.Property == "UpdatedAt" {
			snippetsBefore = append(snippetsBefore, "*entity.UpdatedAt = time.Now()")
		}

		if field.ManyMany == nil {
			data.SQLFields += field.Name + " = $" + strconv.Itoa(index) + ", "
			data.StructFields += "*entity." + field.Property + ", "
			index++
		} else {
			manyMany = append(manyMany, deleteManyMany("*entity.ID", field))
			manyMany = append(manyMany, insertManyMany(field))
		}
	}

	data.SQLFields = strings.TrimSuffix(data.SQLFields, ", ")
	data.StructFields = strings.TrimSuffix(data.StructFields, ", ")

	if len(snippetsBefore) != 0 {
		data.SnippetsBefore = strings.Join(snippetsBefore, "\n")
	}

	if len(manyMany) != 0 {
		data.ManyMany = strings.Join(manyMany, "\n") + "\n"
		data.ManyManyVars = "stmtMmany *sql.Stmt"
	}

	err := tmplUpdate.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
