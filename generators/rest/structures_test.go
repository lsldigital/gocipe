package rest

import (
	"strings"
	"testing"

	"github.com/fluxynet/gocipe/generators"
)

func TestGenerateStructures(t *testing.T) {
	structInfo := generators.StructureInfo{
		Name: "Persons",
		Fields: []generators.FieldInfo{
			{Name: "id", Type: "int64"},
			{Name: "name", Type: "string"},
			{Name: "email", Type: "string"},
			{Name: "gender", Type: "string"},
		},
	}

	output, err := GenerateStructures(structInfo)
	expected := `
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type responseSingle struct {
	Status   bool      ` + "`" + `json:"status"` + "`" + `
	Messages []message ` + "`" + `json:"messages"` + "`" + `
	Entity   *Persons ` + "`" + `json:"entity"` + "`" + `
}

type responseList struct {
	Status   bool        ` + "`" + `json:"status"` + "`" + `
	Messages []message   ` + "`" + `json:"messages"` + "`" + `
	Entities []*Persons ` + "`" + `json:"entities"` + "`" + `
}

type message struct {
	Type    rune   ` + "`" + `json:"type"` + "`" + `
	Message string ` + "`" + `json:"message"` + "`" + `
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
