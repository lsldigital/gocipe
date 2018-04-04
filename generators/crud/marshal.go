package crud

import (
	"bytes"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
)

var tmplMarshal, _ = template.New("GenerateMarshal").Parse(`
//MarshalJSON returns JSON representation of the entity
func (entity *{{.Name}}) MarshalJSON() ([]byte, error) {
	var (
		fields []string
		output bytes.Buffer
		err    error
	)
	{{range .Fields}}
	{{.}}{{end}}

	output.WriteString("{")
	output.WriteString(strings.Join(fields, ","))
	output.WriteString("}")

	return output.Bytes(), err
}
`)

var tmplMarshalBool, _ = template.New("GenerateMarshalBool").Parse(`
	if err != nil && entity.{{.Property}} != nil {
		fields = append(fields, ` + "`" + `"{{.Name}}" = ` + "` + strconv.FormatBool(*entity.{{.Property}})" + `)
	}
`)

var tmplMarshalInt, _ = template.New("GenerateMarshalInt").Parse(`
	if err != nil && entity.{{.Property}} != nil {
		fields = append(fields, ` + "`" + `"{{.Name}}" = ` + "` + strconv.FormatInt(*entity.{{.Property}}, 10)" + `)
	}
`)

var tmplMarshalFloat, _ = template.New("GenerateMarshalFloat").Parse(`
	if err != nil && entity.{{.Property}} != nil {
		fields = append(fields, ` + "`" + `"{{.Name}}" = ` + "` + strconv.FormatFloat(*entity.{{.Property}}, 'f', -1, {{.Size}})" + `)
	}
`)

var tmplMarshalTime, _ = template.New("GenerateMarshalTime").Parse(`
	if err != nil && entity.{{.Property}} != nil {
		fields = append(fields, ` + "`" + `"{{.Name}}" = ` + "` + *entity.{{.Property}}.Format(time.RFC3339)" + `)
	}
`)

var tmplMarshalChars, _ = template.New("GenerateMarshalChars").Parse(`
	if err != nil && entity.{{.Property}} != nil {
		fields = append(fields, ` + "`" + `"{{.Name}}" = ` + "` + string(*entity.{{.Property}})" + `)
	}
`)

var tmplMarshalString, _ = template.New("GenerateMarshalString").Parse(`
	if err != nil && entity.{{.Property}} != nil {
		fields = append(fields, ` + "`" + `"{{.Name}}" = ` + "` + *entity.{{.Property}}" + `)
	}
`)

type fieldMarshalTpl struct {
	Property string
	Name     string
	Size     string
}

//GenerateMarshal generates json marshalling code
func GenerateMarshal(structInfo generators.StructureInfo) (string, error) {
	var (
		output bytes.Buffer
		data   struct {
			Name   string
			Fields []string
		}
	)

	data.Name = structInfo.Name

	for _, field := range structInfo.Fields {
		var (
			def  bytes.Buffer
			name string
			ok   bool
		)

		if name, ok = field.Tags.Lookup("json"); !ok || name == "-" {
			continue
		}

		switch field.Type {
		case "bool":
			tmplMarshalBool.Execute(&def, fieldMarshalTpl{field.Property, field.Name, ""})
		case "string":
			tmplMarshalString.Execute(&def, fieldMarshalTpl{field.Property, field.Name, ""})
		case "byte", "rune":
			tmplMarshalChars.Execute(&def, fieldMarshalTpl{field.Property, field.Name, ""})
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
			tmplMarshalInt.Execute(&def, fieldMarshalTpl{field.Property, field.Name, ""})
		case "float32":
			tmplMarshalFloat.Execute(&def, fieldMarshalTpl{field.Property, field.Name, "32"})
		case "float64":
			tmplMarshalFloat.Execute(&def, fieldMarshalTpl{field.Property, field.Name, "64"})
		case "time.Time":
			tmplMarshalTime.Execute(&def, fieldMarshalTpl{field.Property, field.Name, ""})
		default:
			continue
		}

		data.Fields = append(data.Fields, def.String())
	}

	err := tmplMarshal.Execute(&output, data)

	if err != nil {
		return "", err
	}

	return output.String(), nil
}
