package crud

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
)

var tmplGet, _ = template.New("GenerateGet").Parse(`
//Get returns a single {{.Name}} from database by primary key
func Get(id int64) (*{{.Name}}, error) {
	var entity = New()
	{{if .PreExecHook }}
    if err := crudPreGet(id); err != nil {
		return nil, fmt.Errorf("error executing crudPreGet() in Get(%d) for entity '{{.Name}}': %s", id, err)
	}
    {{end}}
	rows, err := db.Query("SELECT {{.SQLFields}} FROM {{.TableName}} t {{.Joins}}WHERE id = $1 ORDER BY t.id ASC", id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() { {{.JoinVarsDecl}}
		err := rows.Scan({{.StructFields}})
		if err != nil {
			return nil, err
		} {{.JoinVarsAssgn}}
	}
	{{if .PostExecHook }}
	if err = crudPostGet(entity); err != nil {
		return nil, fmt.Errorf("error executing crudPostGet() in Get(%d) for entity '{{.Name}}': %s", id, err)
	}
	{{end}}

	return entity, nil
}
`)

var tmplGetHook, _ = template.New("GenerateGetHook").Parse(`
{{if .PreExecHook }}
func crudPreGet(id int64) error {
	return nil
}
{{end}}
{{if .PostExecHook }}
func crudPostGet(entity *{{.Name}}) error {
	return nil
}
{{end}}
`)

//GenerateGet generates code to get an entity from database
func GenerateGet(structInfo generators.StructureInfo, preExecHook bool, postExecHook bool) (string, error) {
	var (
		output        bytes.Buffer
		joinCount     int
		joins         []string
		joinVarsDecl  []string
		joinVarsAssgn []string
	)
	data := new(struct {
		Name          string
		TableName     string
		SQLFields     string
		StructFields  string
		Joins         string
		JoinVarsDecl  string
		JoinVarsAssgn string
		PreExecHook   bool
		PostExecHook  bool
	})

	data.Name = structInfo.Name
	data.TableName = structInfo.TableName
	data.SQLFields = ""
	data.StructFields = ""
	data.PreExecHook = preExecHook
	data.PostExecHook = postExecHook

	for _, field := range structInfo.Fields {
		if field.ManyMany == nil {
			data.SQLFields += "t." + field.Name + ", "
			data.StructFields += "entity." + field.Property + ", "
		} else {
			data.SQLFields += fmt.Sprintf("jt%d.%s, ", joinCount, field.ManyMany.ThatID)
			joins = append(joins, fmt.Sprintf("%s jt%d ON (t.%s = jt%d.%s)", field.ManyMany.PivotTable, joinCount, field.ManyMany.ThisID, joinCount, field.ManyMany.ThatID))
			joinVarsDecl = append(joinVarsDecl, fmt.Sprintf("j%d int64", joinCount))
			data.StructFields += fmt.Sprintf("&j%d, ", joinCount)
			joinVarsAssgn = append(joinVarsAssgn, fmt.Sprintf("entity.%s = append(entity.%s, j%d)", field.Property, field.Property, joinCount))
			joinCount++
		}
	}

	data.SQLFields = strings.TrimSuffix(data.SQLFields, ", ")
	data.StructFields = strings.TrimSuffix(data.StructFields, ", ")

	if joinCount > 0 {
		data.Joins = "INNER JOIN " + strings.Join(joins, " INNER JOIN ") + " "
		data.JoinVarsDecl = "\nvar (\n\t\t" + strings.Join(joinVarsDecl, "\n\t\t") + "\t\t)\n"
		data.JoinVarsAssgn = "\n\t\t" + strings.Join(joinVarsAssgn, "\n\t\t")
	}

	err := tmplGet.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

// GenerateGetHook will generate 2 functions: crudPreGet() and crudPostGet()
func GenerateGetHook(structInfo generators.StructureInfo, preExecHook bool, postExecHook bool) (string, error) {
	var output bytes.Buffer

	data := new(struct {
		Name         string
		PreExecHook  bool
		PostExecHook bool
	})

	data.Name = structInfo.Name
	data.PreExecHook = preExecHook
	data.PostExecHook = postExecHook

	err := tmplGetHook.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
