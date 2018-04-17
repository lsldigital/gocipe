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
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "E", "message": "Failed to read body"}]}` + "`" + `)
		return
	}

	response.Entity = New()
	err = json.Unmarshal(rawbody, response.Entity)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "E", "message": "Failed to decode body"}]}` + "`" + `)
		return
	}
	response.Entity.ID = nil

	tx, err = db.Begin()
	if err != nil {
		return
	}

	{{if .PreExecHook}}
	if tx, err = restCreatePreExecHook(w, r, response.Entity, tx); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "E", "message": "restCreatePreExecHook failed for '{{.Endpoint}}'"}]}` + "`" + `)
		_ = tx.Rollback()
		return
	}
    {{end}}

	tx, err = response.Entity.Save(tx)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "E", "message": "Save failed"}]}` + "`" + `)
		_ = tx.Rollback()
		return
	}

	{{if .PostExecHook}}
	if tx, err = restCreatePostExecHook(w, r, response.Entity, tx); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "E", "message": "restCreatePreExecHook failed for '{{.Endpoint}}'"}]}` + "`" + `)
		_ = tx.Rollback()
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
		fmt.Fprint(w, ` + "`" + `{"status": false, "messages": [{"type": "E", "message": "JSON encoding failed"}]}` + "`" + `)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(output))
}
`)

var tmplCreateHook, _ = template.New("GenerateCreateHook").Parse(`
{{if .PreExecHook }}
func restCreatePreExecHook(w http.ResponseWriter, r *http.Request, entity *{{.Name}}, tx *sql.Tx) (*sql.Tx, error) {
	return tx, nil
}
{{end}}
{{if .PostExecHook }}
func restCreatePostExecHook(w http.ResponseWriter, r *http.Request, entity *{{.Name}}, tx *sql.Tx) (*sql.Tx, error) {
	return tx, nil
}
{{end}}
`)

//GenerateCreate will generate a REST handler function for Create
func GenerateCreate(structInfo generators.StructureInfo, PreExecHook bool, PostExecHook bool) (string, error) {
	var output bytes.Buffer
	data := struct {
		Endpoint     string
		PreExecHook  bool
		PostExecHook bool
	}{strings.ToLower(structInfo.Name), PreExecHook, PostExecHook}

	err := tmplCreate.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

// GenerateCreateHook will generate 2 functions: restCreatePreExecHook() and restCreatePostExecHook()
func GenerateCreateHook(structInfo generators.StructureInfo, PreExecHook bool, PostExecHook bool) (string, error) {
	var output bytes.Buffer

	data := new(struct {
		Name         string
		PreExecHook  bool
		PostExecHook bool
	})

	data.Name = structInfo.Name
	data.PreExecHook = PreExecHook
	data.PostExecHook = PostExecHook

	err := tmplCreateHook.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
