package rest

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
	"github.com/jinzhu/inflection"
)

var tmplRoutes, _ = template.New("GenerateRoutes").Parse(`
//RegisterRoutes registers routes with a mux Router
func RegisterRoutes(router *mux.Router) {
{{.Routes}}
}
`)

//GenerateRoutes create function RegisterRoutes to register route endpoints
func GenerateRoutes(structInfo generators.StructureInfo, g generator) (string, error) {
	var (
		output bytes.Buffer
		routes []string
	)

	endpoint := inflection.Plural(strings.ToLower(structInfo.Name))
	data := new(struct {
		Routes string
	})

	if g.GenerateGet {
		routes = append(routes, "\t"+`router.HandleFunc("/`+endpoint+`/{id}", RestGet).Methods("GET")`)
	}

	if g.GenerateList {
		routes = append(routes, "\t"+`router.HandleFunc("/`+endpoint+`", RestList).Methods("GET")`)
	}

	if g.GenerateDelete {
		routes = append(routes, "\t"+`router.HandleFunc("/`+endpoint+`/{id}", RestDelete).Methods("DELETE")`)
	}

	if g.GenerateCreate {
		routes = append(routes, "\t"+`router.HandleFunc("/`+endpoint+`", RestCreate).Methods("POST")`)
	}

	if g.GenerateUpdate {
		routes = append(routes, "\t"+`router.HandleFunc("/`+endpoint+`/{id}", RestUpdate).Methods("PUT")`)
	}

	data.Routes = strings.Join(routes, "\n")
	err := tmplRoutes.Execute(&output, data)

	if err != nil {
		return "", err
	}

	return output.String(), nil
}
