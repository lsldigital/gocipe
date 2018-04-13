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

	rawbody, err = ioutil.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "E", "message": "Failed to read body"}]}` + "`" + `)
		return
	}

	err = json.Unmarshal(rawbody, response.Entity)
	if err != nil {
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "E", "message": "Failed to decode body"}]}` + "`" + `)
			return
		}
	}
	response.Entity.ID = &id

	{{if .PreExecHook}}
    if err = restUpdatePreExecHook(w, r, response.Entity); err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "E", "message": err.Error()}]}` + "`" + `)
        return
    }
    {{end}}

	err = response.Entity.Save()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "E", "message": "Save failed"}]}` + "`" + `)
		return
	}

	{{if .PostExecHook}}
    if err = restUpdatePostExecHook(w, r, response.Entity); err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "E", "message": err.Error()}]}` + "`" + `)
        return
    }
    {{end}}

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

var tmplUpdateHook, _ = template.New("GenerateUpdateHook").Parse(`
{{if .PreExecHook }}
func restUpdatePreExecHook(w http.ResponseWriter, r *http.Request, entity *{{.Name}}) error {
	return nil
}
{{end}}
{{if .PostExecHook }}
func restUpdatePostExecHook(w http.ResponseWriter, r *http.Request, entity *{{.Name}}) error {
	return nil
}
{{end}}
`)

//GenerateUpdate will generate a REST handler function for Update
func GenerateUpdate(structInfo generators.StructureInfo, PreExecHook bool, PostExecHook bool) (string, error) {
	var output bytes.Buffer
	data := struct {
		Endpoint     string
		PreExecHook  bool
		PostExecHook bool
	}{strings.ToLower(structInfo.Name), PreExecHook, PostExecHook}

	err := tmplUpdate.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

// GenerateUpdateHook will generate 2 functions: restUpdatePreExecHook() and restUpdatePostExecHook()
func GenerateUpdateHook(structInfo generators.StructureInfo, PreExecHook bool, PostExecHook bool) (string, error) {
	var output bytes.Buffer

	data := new(struct {
		Name         string
		PreExecHook  bool
		PostExecHook bool
	})

	data.Name = structInfo.Name
	data.PreExecHook = PreExecHook
	data.PostExecHook = PostExecHook

	err := tmplUpdateHook.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
