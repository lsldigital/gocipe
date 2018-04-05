package rest

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
)

var tmplList, _ = template.New("GenerateList").Parse(`
//RestList is a REST endpoint for GET /{{.Endpoint}}
func RestList(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		response responseList
		filters  []models.ListFilter
		query    url.Values
	)
	query = r.URL.Query()
	
	{{.Filters}}

	response.Entities, err = List(filters)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "E", "message": "An error occurred"}]}` + "`" + `)
		return
	}

	response.Status = true

	output, err := json.Marshal(response)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "E", "message": "JSON encoding failed"}]}` + "`" + `)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(output))
}
`)

var tmplListFilterBool, _ = template.New("GenerateListFilterBool").Parse(`
	if val := query.Get("{{.Name}}"); val != "" {
		if t, e := strconv.ParseBool(val); e == nil {
			filters = append(filters, models.ListFilter{"{{.Name}}", "=", t})
		}
	}
`)

var tmplListFilterString, _ = template.New("GenerateListFilterString").Parse(`
	if val := query.Get("{{.Name}}"); val != "" {
		filters = append(filters, models.ListFilter{"{{.Name}}", "=", val})
	}

	if val := query.Get("{{.Name}}Lk"); val != "" {
		filters = append(filters, models.ListFilter{"{{.Name}}", "LIKE", "%" + val + "%"})
	}
`)

var tmplListFilterDate, _ = template.New("GenerateListFilterDate").Parse(`
	if val := query.Get("{{.Name}}"); len(val) == 16 {
		if t, e := time.Parse("2006-01-02-15-04", val); e == nil {
			filters = append(filters, models.ListFilter{"{{.Name}}", "=", t})
		}
	}

	if val := query.Get("{{.Name}}Gt"); len(val) == 16 {
		if t, e := time.Parse("2006-01-02-15-04", val); e == nil {
			filters = append(filters, models.ListFilter{"{{.Name}}", ">", t})
		}
	}

	if val := query.Get("{{.Name}}Ge"); len(val) == 16 {
		if t, e := time.Parse("2006-01-02-15-04", val); e == nil {
			filters = append(filters, models.ListFilter{"{{.Name}}", ">=", t})
		}
	}

	if val := query.Get("{{.Name}}Lt"); len(val) == 16 {
		if t, e := time.Parse("2006-01-02-15-04", val); e == nil {
			filters = append(filters, models.ListFilter{"{{.Name}}", "<", t})
		}
	}

	if val := query.Get("{{.Name}}Le"); len(val) == 16 {
		if t, e := time.Parse("2006-01-02-15-04", val); e == nil {
			filters = append(filters, models.ListFilter{"{{.Name}}", "<=", t})
		}
	}
`)

//GenerateList will generate a REST handler function for List
func GenerateList(structInfo generators.StructureInfo) (string, error) {
	var (
		output  bytes.Buffer
		filters []string
		err     error
		data    struct {
			Endpoint string
			Filters  string
		}
	)

	data.Endpoint = strings.ToLower(structInfo.Name)

	for _, field := range structInfo.Fields {
		var segment bytes.Buffer

		if !field.Filterable {
			continue
		}

		switch field.Type {
		case "bool":
			err = tmplListFilterBool.Execute(&segment, struct{ Name string }{field.Name})
		case "string":
			err = tmplListFilterString.Execute(&segment, struct{ Name string }{field.Name})
		case "time.Time":
			err = tmplListFilterDate.Execute(&segment, struct{ Name string }{field.Name})
		default:
			continue
		}

		if err != nil {
			return "", err
		}

		filters = append(filters, segment.String())
	}

	data.Filters = strings.Join(filters, "\n")
	err = tmplList.Execute(&output, data)

	if err != nil {
		return "", err
	}

	return output.String(), nil
}
