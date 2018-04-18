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
		{{.ManyManyVars}}
	)

	if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return err
		}
	}

	{{.SnippetsBefore}}
	
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
	{{.ManyMany}}
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

var tmplManyManyInsert, _ = template.New("ManyManyInsert").Parse(`
	stmtMmany, err = tx.Prepare("INSERT INTO {{.PivotTable}} ({{.ThisID}}, {{.ThatID}}) VALUES ($1, $2)")
	
	if err != nil {
		return fmt.Errorf("error preparing transaction statement in ManyManyInsert(%d) {{.ThisID}}-{{.ThatID}} for table '{{.PivotTable}}': %s", *entity.ID, err)
	}

	for _, relatedID := range entity.{{.RelatedProperty}} {
		_, err = stmtMmany.Exec(entity.ID, relatedID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error executing transaction statement in ManyManyInsert(%d) {{.ThisID}}-{{.ThatID}} for table '{{.PivotTable}}': %s", *entity.ID, err)
		}
	}
`)

//GenerateInsert generate function to insert an entity in database
func GenerateInsert(structInfo generators.StructureInfo, PreExecHook bool, PostExecHook bool) (string, error) {
	var (
		output         bytes.Buffer
		snippetsBefore []string
		manyMany       []string
	)

	data := new(struct {
		Name            string
		TableName       string
		SQLFields       string
		SQLPlaceholders string
		StructFields    string
		SnippetsBefore  string
		ManyMany        string
		ManyManyVars    string
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
		if field.Property == "ID" {
			continue
		} else if field.Property == "CreatedAt" {
			snippetsBefore = append(snippetsBefore, "*entity.CreatedAt = time.Now()")
		} else if field.Property == "UpdatedAt" {
			snippetsBefore = append(snippetsBefore, "*entity.UpdatedAt = time.Now()")
		}

		if field.ManyMany == nil {
			data.SQLFields += field.Name + ", "
			data.SQLPlaceholders += "$" + strconv.Itoa(i) + ", "
			data.StructFields += "*entity." + field.Property + ", "
		} else {
			manyMany = append(manyMany, insertManyMany(field))
		}
	}

	data.SQLFields = strings.TrimSuffix(data.SQLFields, ", ")
	data.SQLPlaceholders = strings.TrimSuffix(data.SQLPlaceholders, ", ")
	data.StructFields = strings.TrimSuffix(data.StructFields, ", ")

	if len(snippetsBefore) != 0 {
		data.SnippetsBefore = strings.Join(snippetsBefore, "\n")
	}

	if len(manyMany) != 0 {
		data.ManyMany = strings.Join(manyMany, "\n") + "\n"
		data.ManyManyVars = "stmtMmany *sql.Stmt"
	}

	err := tmplInsert.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

func insertManyMany(field generators.FieldInfo) string {
	var output bytes.Buffer

	data := new(struct {
		PivotTable      string
		RelatedProperty string
		ThatID          string
		ThisID          string
	})

	data.PivotTable = field.ManyMany.PivotTable
	data.RelatedProperty = field.Property
	data.ThisID = field.ManyMany.ThisID
	data.ThatID = field.ManyMany.ThatID

	err := tmplManyManyInsert.Execute(&output, data)
	if err != nil {
		return ""
	}

	return output.String()
}
