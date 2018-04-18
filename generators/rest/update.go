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
		tx       *sql.Tx
		{{if .Hooks}}stop     bool{{end}}
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
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "error", "text": "Invalid ID"}]}` + "`" + `)
		return
	}

	response.Entity, err = Get(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "error", "text": "An error occurred"}]}` + "`" + `)
		return
	}

	if response.Entity == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "error", "text": "Entity not found"}]}` + "`" + `)
		return
	}

	rawbody, err = ioutil.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "error", "text": "Failed to read body"}]}` + "`" + `)
		return
	}

	err = json.Unmarshal(rawbody, response.Entity)
	if err != nil {
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "error", "text": "Failed to decode body"}]}` + "`" + `)
			return
		}
	}
	response.Entity.ID = &id

	tx, err = db.Begin()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "error", "text": "Failed to process"}]}` + "`" + `)
		return
	}

	{{if .PreExecHook}}
    if stop, err = restPreUpdate(w, r, response.Entity, tx); err != nil {
		tx.Rollback()
        return
    } else if stop {
		return
	}
    {{end}}

	err = response.Entity.Save(tx, false)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "error", "text": "Save failed"}]}` + "`" + `)
		return
	}

	{{if .PostExecHook}}
    if stop, err = restPostUpdate(w, r, response.Entity, tx); err != nil {
		tx.Rollback()
        return
    } else if stop {
		return
	}
	{{end}}
	
	if err = tx.Commit(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "E", "message": "RestUpdate could not commit transaction"}]}` + "`" + `)
		return
	}

	output, err := json.Marshal(response)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "error", "text": "JSON encoding failed"}]}` + "`" + `)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(output))
}
`)

var tmplUpdateHook, _ = template.New("GenerateUpdateHook").Parse(`
{{if .PreExecHook }}
func restPreUpdate(w http.ResponseWriter, r *http.Request, entity *{{.Name}}, tx *sql.Tx) (bool, error) {
	return false, nil
}
{{end}}
{{if .PostExecHook }}
func restPostUpdate(w http.ResponseWriter, r *http.Request, entity *{{.Name}}, tx *sql.Tx) (bool, error) {
	return false, nil
}
{{end}}
`)

//GenerateUpdate will generate a REST handler function for Update
func GenerateUpdate(structInfo generators.StructureInfo, preExecHook bool, postExecHook bool) (string, error) {
	var output bytes.Buffer
	data := struct {
		Endpoint     string
		PreExecHook  bool
		PostExecHook bool
		Hooks        bool
	}{strings.ToLower(structInfo.Name), preExecHook, postExecHook, preExecHook || postExecHook}

	err := tmplUpdate.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

// GenerateUpdateHook will generate 2 functions: restPreUpdate() and restPostUpdate()
func GenerateUpdateHook(structInfo generators.StructureInfo, preExecHook bool, postExecHook bool) (string, error) {
	var output bytes.Buffer

	data := new(struct {
		Name         string
		PreExecHook  bool
		PostExecHook bool
	})

	data.Name = structInfo.Name
	data.PreExecHook = preExecHook
	data.PostExecHook = postExecHook

	err := tmplUpdateHook.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
