package crud

import (
	"bytes"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
)

var tmplUnmarshal, _ = template.New("GenerateUnmarshal").Parse(`
//UnmarshalJSON returns an instance of the entity from JSON representation
func (entity *{{.Name}}) UnmarshalJSON(in []byte) error {
	var (
		err error
		raw map[string]string
	)

	err = json.Unmarshal(in, &raw)

	if err != nil {
		return err
	}
	{{range .Fields}}
	{{.}}{{end}}

	return err
}
`)

var tmplUnmarshalBool, _ = template.New("GenerateUnmarshalBool").Parse(`
	if value, ok := raw["{{.Name}}"]; err != nil && ok && value != "null" {
		*entity.{{.Property}}, err = strconv.ParseBool(value)
	}
`)

var tmplUnmarshalInt, _ = template.New("GenerateUnmarshalInt").Parse(`
	if value, ok := raw["{{.Name}}"]; err != nil && ok && value != "null" {
		var i int64
		i, err = strconv.ParseInt(value, 10, {{.Size}})

		if err != nil {
			*entity.{{.Property}} = {{.Type}}(i)
		}
	}
`)

var tmplUnmarshalRune, _ = template.New("GenerateUnmarshalRune").Parse(`
	if value, ok := raw["{{.Name}}"]; err != nil && ok && value != "null" {
		if len(value) > 0 {
			*entity.{{.Property}} = string(value[0])
		}
	}
`)

var tmplUnmarshalFloat, _ = template.New("GenerateUnmarshalFloat").Parse(`
	if value, ok := raw["{{.Name}}"]; err != nil && ok && value != "null" {
		*entity.{{.Property}} = strconv.ParseFloat(value, 10, {{.Size}})
	}
`)

var tmplUnmarshalTime, _ = template.New("GenerateUnmarshalTime").Parse(`
	if value, ok := raw["{{.Name}}"]; err != nil && ok && value != "null" {
		*entity.{{.Property}}, err = time.Parse(time.RFC3339, value)
	}
`)

var tmplUnmarshalString, _ = template.New("GenerateUnmarshalString").Parse(`
	if value, ok := raw["{{.Name}}"]; err != nil && ok && value != "null" {
		*entity.{{.Property}} = value
	}
`)

//GenerateUnmarshal generates json unmarshalling code
func GenerateUnmarshal(structInfo generators.StructureInfo) (string, error) {
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
			tmplUnmarshalBool.Execute(&def, fieldMarshalTpl{field.Property, field.Name, "", field.Type})
		case "string", "byte":
			tmplUnmarshalString.Execute(&def, fieldMarshalTpl{field.Property, field.Name, "", field.Type})
		case "rune":
			tmplUnmarshalRune.Execute(&def, fieldMarshalTpl{field.Property, field.Name, "", field.Type})
		case "int", "uint":
			tmplUnmarshalInt.Execute(&def, fieldMarshalTpl{field.Property, field.Name, "0", field.Type})
		case "int8", "uint8":
			tmplUnmarshalInt.Execute(&def, fieldMarshalTpl{field.Property, field.Name, "8", field.Type})
		case "int16", "uint16":
			tmplUnmarshalInt.Execute(&def, fieldMarshalTpl{field.Property, field.Name, "16", field.Type})
		case "int32", "uint32":
			tmplUnmarshalInt.Execute(&def, fieldMarshalTpl{field.Property, field.Name, "32", field.Type})
		case "int64", "uint64":
			tmplUnmarshalInt.Execute(&def, fieldMarshalTpl{field.Property, field.Name, "64", field.Type})
		case "float32":
			tmplUnmarshalFloat.Execute(&def, fieldMarshalTpl{field.Property, field.Name, "32", field.Type})
		case "float64":
			tmplUnmarshalFloat.Execute(&def, fieldMarshalTpl{field.Property, field.Name, "64", field.Type})
		case "time.Time":
			tmplUnmarshalTime.Execute(&def, fieldMarshalTpl{field.Property, field.Name, "", field.Type})
		default:
			continue
		}

		data.Fields = append(data.Fields, def.String())
	}

	err := tmplUnmarshal.Execute(&output, data)

	if err != nil {
		return "", err
	}

	return output.String(), nil
}
