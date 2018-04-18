package crud

import (
	"bytes"
	"fmt"
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
    if filters, err = crudPreList(filters); err != nil {
		return nil, fmt.Errorf("error executing crudPreList() in List(filters) for entity '{{.Name}}': %s", err)
	}
    {{end}}
	for i, filter := range filters {
		segments = append(segments, filter.Field+" "+filter.Operation+" $"+strconv.Itoa(i+1))
		values = append(values, filter.Value)
	}

	if len(segments) != 0 {
		query += " WHERE " + strings.Join(segments, " AND ")
	}

	rows, err := db.Query(query+" ORDER BY id ASC", values...) {{.ManyIndexDecl}}
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

		list = append(list, entity) {{.ManyIndexAssign}}
	}
	{{.ManyMany}}
	{{if .PostExecHook }}
	if list, err = crudPostList(list); err != nil {
		return nil, fmt.Errorf("error executing crudPostList() in List(filters) for entity '{{.Name}}': %s", err)
	}
	{{end}}
	return list, nil
}
{{.ManyManyLoadRelated}}
`)

var tmplListHook, _ = template.New("GenerateListHook").Parse(`
{{if .PreExecHook }}
func crudPreList(filters []models.ListFilter) ([]models.ListFilter, error) {
	return filters, nil
}
{{end}}
{{if .PostExecHook }}
func crudPostList(list []*{{.Name}}) ([]*{{.Name}}, error) {
	return list, nil
}
{{end}}
`)

var tmplListMany, _ = template.New("GenerateListMany").Parse(`
	if related, e := loadRelated(indexID, "{{.ThisID}}", "{{.ThatID}}", "{{.PivotTable}}"); e == nil {
		for i, v := range related {
			indexID[i].{{.Property}} = append(indexID[i].{{.Property}}, v)
		}
	} else {
		return nil, err
	}
`)

var tmplListManyLoadRelated, _ = template.New("GenerateListManyLoadRelated").Parse(`
func loadRelated(indexID map[int64]*{{.Name}}, thisid string, thatid string, pivot string) (map[int64]int64, error) {
	var (
		placeholder string
		values  []interface{}
		idthis  int64
		idthat  int64
	)

	related := make(map[int64]int64)

	c := 1
	for i := range indexID {
		placeholder += "$" + strconv.Itoa(c) + ","
		values = append(values, i)
		c++
	}
	placeholder = strings.TrimRight(placeholder, ",")

	rows, err := db.Query("SELECT "+thisid+", "+thatid+" FROM "+pivot+" WHERE "+thisid+" IN ("+placeholder+")", values...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&idthis, &idthat)
		if err != nil {
			return nil, err
		}
		related[idthis] = idthat
	}

	return related, nil
}
`)

//GenerateList returns code to return a list of entities from database
func GenerateList(structInfo generators.StructureInfo, preExecHook bool, postExecHook bool) (string, error) {
	var (
		output   bytes.Buffer
		manyMany []string
	)

	data := new(struct {
		Name                string
		TableName           string
		SQLFields           string
		StructFields        string
		ManyIndexDecl       string
		ManyIndexAssign     string
		ManyMany            string
		ManyManyLoadRelated string
		PreExecHook         bool
		PostExecHook        bool
	})

	data.Name = structInfo.Name
	data.TableName = structInfo.TableName
	data.SQLFields = ""
	data.StructFields = ""
	data.PreExecHook = preExecHook
	data.PostExecHook = postExecHook

	for _, field := range structInfo.Fields {
		if field.ManyMany == nil {
			data.SQLFields += field.Name + ", "
			data.StructFields += "entity." + field.Property + ", "
		} else {
			manyMany = append(manyMany, manyGenerateList(&field))
		}
	}

	data.SQLFields = strings.TrimSuffix(data.SQLFields, ", ")
	data.StructFields = strings.TrimSuffix(data.StructFields, ", ")

	if len(manyMany) > 0 {
		data.ManyIndexDecl = fmt.Sprintf("\n\tindexID := make(map[int64]*%s)", structInfo.Name)
		data.ManyIndexAssign = "\n\t\tindexID[*entity.ID] = entity"
		data.ManyMany = "\n" + strings.Join(manyMany, "\n")
		data.ManyManyLoadRelated = manyGenerateListLoadRelated(structInfo.Name)
	}

	err := tmplList.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

// GenerateListHook will generate 2 functions: crudPreList() and crudPostList()
func GenerateListHook(structInfo generators.StructureInfo, preExecHook bool, postExecHook bool) (string, error) {
	var output bytes.Buffer

	data := new(struct {
		Name         string
		PreExecHook  bool
		PostExecHook bool
	})

	data.Name = structInfo.Name
	data.PreExecHook = preExecHook
	data.PostExecHook = postExecHook

	err := tmplListHook.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

func manyGenerateListLoadRelated(entityName string) string {
	var output bytes.Buffer

	data := struct {
		Name string
	}{entityName}

	err := tmplListManyLoadRelated.Execute(&output, data)
	if err != nil {
		return ""
	}

	return output.String()
}

func manyGenerateList(field *generators.FieldInfo) string {
	var output bytes.Buffer

	data := struct {
		Property   string
		ThisID     string
		ThatID     string
		PivotTable string
	}{field.Property, field.ManyMany.ThisID, field.ManyMany.ThatID, field.ManyMany.PivotTable}

	err := tmplListMany.Execute(&output, data)
	if err != nil {
		return ""
	}

	return output.String()
}
