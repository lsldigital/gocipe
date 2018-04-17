package rest

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
)

var tmplCreate, _ = template.New("GenerateCreate").Parse(`
//RestCreate is a REST endpoint for POST /{{.Endpoint}}
func RestCreate(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		rawbody  []byte
		response responseSingle
		tx       *sql.Tx
	)

	rawbody, err = ioutil.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "error", "text": "Failed to read body"}]}` + "`" + `)
		return
	}

	response.Entity = New()
	err = json.Unmarshal(rawbody, response.Entity)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "error", "text": "Failed to decode body"}]}` + "`" + `)
		return
	}
	response.Entity.ID = nil

	tx, err = db.Begin()
	if err != nil {
		return
	}

	{{if .PreExecHook}}
	if err = restPreCreate(w, r, response.Entity, tx); err != nil {
		tx.Rollback()
		return
	}
    {{end}}

	err = response.Entity.Save(tx, false)
	if err != nil {
		tx.Rollback()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "error", "text": "Save failed"}]}` + "`" + `)
		return
	}

	{{if .PostExecHook}}
	if err = restPostCreate(w, r, response.Entity, tx); err != nil {
		tx.Rollback()
		return
	}
	{{end}}
	
	if err = tx.Commit(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "E", "message": "RestCreate could not commit transaction"}]}` + "`" + `)
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

var tmplCreateHook, _ = template.New("GenerateCreateHook").Parse(`
{{if .PreExecHook }}
func restPreCreate(w http.ResponseWriter, r *http.Request, entity *{{.Name}}, tx *sql.Tx) error {
	return nil
}
{{end}}
{{if .PostExecHook }}
func restPostCreate(w http.ResponseWriter, r *http.Request, entity *{{.Name}}, tx *sql.Tx) error {
	return nil
}
{{end}}
`)

//GenerateCreate will generate a REST handler function for Create
func GenerateCreate(structInfo generators.StructureInfo, preExecHook bool, postExecHook bool) (string, error) {
	var output bytes.Buffer
	data := struct {
		Endpoint     string
		PreExecHook  bool
		PostExecHook bool
	}{strings.ToLower(structInfo.Name), preExecHook, postExecHook}

	err := tmplCreate.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

// GenerateCreateHook will generate 2 functions: restPreCreate() and restPostCreate()
func GenerateCreateHook(structInfo generators.StructureInfo, preExecHook bool, postExecHook bool) (string, error) {
	var output bytes.Buffer

	data := new(struct {
		Name         string
		PreExecHook  bool
		PostExecHook bool
	})

	data.Name = structInfo.Name
	data.PreExecHook = preExecHook
	data.PostExecHook = postExecHook

	err := tmplCreateHook.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
