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
	)

	response.Entities, err = List()

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

//GenerateList will generate a REST handler function for List
func GenerateList(structInfo generators.StructureInfo) (string, error) {
	var output bytes.Buffer
	data := struct {
		Endpoint string
	}{strings.ToLower(structInfo.Name)}

	err := tmplList.Execute(&output, data)

	if err != nil {
		return "", err
	}

	return output.String(), nil
}
