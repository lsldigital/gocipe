package crud

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
)

var tmplList, _ = template.New("GenerateList").Parse(`
// List returns a slice containing {{.Name}} records
func List(filters []models.ListFilter) ([]*{{.Name}}, error) {
	var (
		list     []*{{.Name}}
		segments []string
		values   []interface{}
		err      error
	)

	query := "SELECT {{.SQLFields}} FROM {{.TableName}}"
	{{if .PreExecHook }}
    if filters, err = crudListPreExecHook(filters); err != nil {
		fmt.Printf("Error executing crudListPreExecHook() in List(filters) for entity '{{.Name}}': %s", err.Error())
        return nil, err
	}
    {{end}}
	for i, filter := range filters {
		segments = append(segments, filter.Field+" "+filter.Operation+" $"+strconv.Itoa(i+1))
		values = append(values, filter.Value)
	}

	if len(segments) != 0 {
		query += " WHERE " + strings.Join(segments, " AND ")
	}

	rows, err := db.Query(query+" ORDER BY id ASC", values...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		entity := New()
		err := rows.Scan({{.StructFields}})
		if err != nil {
			return nil, err
		}

		list = append(list, entity)
	}
	{{if .PostExecHook }}
	if list, err = crudListPostExecHook(list); err != nil {
		fmt.Printf("Error executing crudListPostExecHook() in List(filters) for entity '{{.Name}}': %s", err.Error())
		return nil, err
	}
	{{end}}
	return list, nil
}
`)

var tmplListHook, _ = template.New("GenerateListHook").Parse(`
{{if .PreExecHook }}
func crudListPreExecHook(filters []models.ListFilter) ([]models.ListFilter, error) {
	return filters, nil
}
{{end}}
{{if .PostExecHook }}
func crudListPostExecHook(list []*{{.Name}}) ([]*{{.Name}}, error) {
	return list, nil
}
{{end}}
`)

//GenerateList returns code to return a list of entities from database
func GenerateList(structInfo generators.StructureInfo, PreExecHook bool, PostExecHook bool) (string, error) {
	var output bytes.Buffer
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
	data.StructFields = ""
	data.PreExecHook = PreExecHook
	data.PostExecHook = PostExecHook

	for _, field := range structInfo.Fields {
		data.SQLFields += field.Name + ", "
		data.StructFields += "entity." + field.Property + ", "
	}

	data.SQLFields = strings.TrimSuffix(data.SQLFields, ", ")
	data.StructFields = strings.TrimSuffix(data.StructFields, ", ")

	err := tmplList.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

// GenerateListHook will generate 2 functions: crudListPreExecHook() and crudListPostExecHook()
func GenerateListHook(structInfo generators.StructureInfo, PreExecHook bool, PostExecHook bool) (string, error) {
	var output bytes.Buffer

	data := new(struct {
		Name         string
		PreExecHook  bool
		PostExecHook bool
	})

	data.Name = structInfo.Name
	data.PreExecHook = PreExecHook
	data.PostExecHook = PostExecHook

	err := tmplListHook.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
