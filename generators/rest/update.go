package rest

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
)

var tmplUpdate, _ = template.New("GenerateUpdate").Parse(`
//RestUpdate is a REST endpoint for PUT /{{.Endpoint}}/{id}
func RestUpdate(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		rawbody  []byte
		id       int64
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
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "error", "message": "Invalid ID"}]}` + "`" + `)
		return
	}

	response.Entity, err = Get(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "error", "message": "An error occurred"}]}` + "`" + `)
		return
	}

	if response.Entity == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "error", "message": "Entity not found"}]}` + "`" + `)
		return
	}

	rawbody, err = ioutil.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "error", "message": "Failed to read body"}]}` + "`" + `)
		return
	}

	err = json.Unmarshal(rawbody, response.Entity)
	if err != nil {
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "error", "message": "Failed to decode body"}]}` + "`" + `)
			return
		}
	}
	response.Entity.ID = &id

	err = response.Entity.Save()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "error", "message": "Save failed"}]}` + "`" + `)
		return
	}

	output, err := json.Marshal(response)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "error", "message": "JSON encoding failed"}]}` + "`" + `)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(output))
}
`)

//GenerateUpdate will generate a REST handler function for Update
func GenerateUpdate(structInfo generators.StructureInfo) (string, error) {
	var output bytes.Buffer
	data := struct {
		Endpoint string
	}{strings.ToLower(structInfo.Name)}

	err := tmplUpdate.Execute(&output, data)

	if err != nil {
		return "", err
	}

	return output.String(), nil
}
