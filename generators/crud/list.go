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
		list []*{{.Name}}
		segments []string
		values []interface{}
	)

	query := "SELECT {{.SQLFields}} FROM {{.TableName}}"

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

	return list, nil
}
`)

//GenerateList returns code to return a list of entities from database
func GenerateList(structInfo generators.StructureInfo) (string, error) {
	var output bytes.Buffer
	data := new(struct {
		Name         string
		TableName    string
		SQLFields    string
		StructFields string
	})

	data.Name = structInfo.Name
	data.TableName = structInfo.TableName
	data.SQLFields = ""
	data.StructFields = ""

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
