package rest

import (
	"bytes"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
)

var tmplStructures, _ = template.New("GenerateStructures").Parse(`
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
	Entity   *{{.Name}} ` + "`" + `json:"entity"` + "`" + `
}

type responseList struct {
	Status   bool        ` + "`" + `json:"status"` + "`" + `
	Messages []message   ` + "`" + `json:"messages"` + "`" + `
	Entities []*{{.Name}} ` + "`" + `json:"entities"` + "`" + `
}

type message struct {
	Type    rune   ` + "`" + `json:"type"` + "`" + `
	Message string ` + "`" + `json:"message"` + "`" + `
}
`)

//GenerateStructures will generate structures used for the REST endpoints
func GenerateStructures(structInfo generators.StructureInfo) (string, error) {
	var output bytes.Buffer
	data := struct {
		Name string
	}{structInfo.Name}

	err := tmplStructures.Execute(&output, data)

	if err != nil {
		return "", err
	}

	return output.String(), nil
}
