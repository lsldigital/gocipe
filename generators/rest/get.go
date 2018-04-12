package rest

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
)

var tmplGet, _ = template.New("GenerateGet").Parse(`
//RestGet is a REST endpoint for GET /{{.Endpoint}}/{id}
func RestGet(w http.ResponseWriter, r *http.Request) {
	var (
		id       int64
		err      error
		response responseSingle
	)

	vars := mux.Vars(r)
	valid := false
	if _, ok := vars["id"]; ok {
		id, err = strconv.ParseInt(vars["id"], 10, 64)
		valid = err == nil && id > 0
	}

	if !valid {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "E", "message": "Invalid ID"}]}` + "`" + `)
		return
	}

	response.Entity, err = Get(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "E", "message": "An error occurred"}]}` + "`" + `)
		return
	}

	if response.Entity == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "E", "message": "Entity not found"}]}` + "`" + `)
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

var tmplGetHook, _ = template.New("GenerateGetHook").Parse(`
{{if .PreExecHook }}
func restGetPreExecHook(id int64) (*User, error) {
	return New(), nil
}
{{end}}
{{if .PostExecHook }}
func restGetPostExecHook(entity *User) (*User, error) {
	return entity, nil
}
{{end}}
`)

//GenerateGet will generate a REST handler function for GET
func GenerateGet(structInfo generators.StructureInfo, PreExecHook bool, PostExecHook bool) (string, error) {
	var output bytes.Buffer
	data := struct {
		Endpoint     string
		PreExecHook  bool
		PostExecHook bool
	}{strings.ToLower(structInfo.Name), PreExecHook, PostExecHook}

	err := tmplGet.Execute(&output, data)

	if err != nil {
		return "", err
	}

	return output.String(), nil
}

// GenerateGetHook will generate 2 functions: restGetPreExecHook() and restGetPostExecHook()
func GenerateGetHook(PreExecHook bool, PostExecHook bool) (string, error) {
	var output bytes.Buffer

	data := new(struct {
		PreExecHook  bool
		PostExecHook bool
	})

	data.PreExecHook = PreExecHook
	data.PostExecHook = PostExecHook

	err := tmplGetHook.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
