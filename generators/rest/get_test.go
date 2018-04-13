package rest

import (
	"strings"
	"testing"

	"github.com/fluxynet/gocipe/generators"
)

func TestGenerateGet(t *testing.T) {
	structInfo := generators.StructureInfo{
		Name: "Persons",
		Fields: []generators.FieldInfo{
			{Name: "id", Type: "int64"},
			{Name: "name", Type: "string"},
			{Name: "email", Type: "string"},
			{Name: "gender", Type: "string"},
		},
	}

	output, err := GenerateGet(structInfo)
	expected := `
//RestGet is a REST endpoint for GET /persons/{id}
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
`

	if err != nil {
		t.Errorf("Got error: %s", err)
	} else if strings.Compare(output, expected) != 0 {
		t.Errorf(`Output not as expected. Output length = %d Expected length = %d
--- # Output ---
%s
--- # Expected ---
%s
------------------
`, len(output), len(expected), output, expected)
	}
}
