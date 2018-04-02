package crud

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
)

var tmplModel, _ = template.New("GenerateModel").Parse(`
import "database/sql"

var db *sql.DB

// Inject allows injection of services into the package
func Inject(database *sql.DB) {
	db = database
}

//New return a new {{.Name}} instance
func New() *{{.Name}} {
	entity := new({{.Name}})
	{{.Fields}}

	return entity
}
`)

//GenerateModel generates code for basic entity
func GenerateModel(structInfo generators.StructureInfo) (string, error) {
	var (
		output bytes.Buffer
		fields []string
	)

	data := new(struct {
		Name   string
		Fields string
	})

	data.Name = structInfo.Name

	for _, field := range structInfo.Fields {
		item := "entity." + field.Property + " = new(" + strings.TrimLeft(field.Type, "*") + ")"
		fields = append(fields, item)
	}

	data.Fields = strings.Join(fields, "\n    ")

	err := tmplModel.Execute(&output, data)

	if err != nil {
		return "", err
	}

	return output.String(), nil
}
