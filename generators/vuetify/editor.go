package vuetify

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
)

var tmplEditor, x = template.New("GenerateEditor").Parse(`
<template>
    <div class="container">
        <h1>{{.Name}}</h1>
		<v-alert :type="message.type" :value="true" v-for="(message, index) in messages" :key="index">
		᚜ message.text ᚛
		</v-alert>
  
        {{.FieldsMarkup}}

        <v-btn color="primary" @click="save()">Save</v-btn>
	</div>
</template>
  
<script>
import axios from "axios";

export default {
    props: ["id"],
    created() {
        if (!this.id) {
            return
        }

        this.axios.get("{{.ServiceURL}}/{{.Endpoint}}/" + this.id).then(response => {
            this.id = response.data.entity.id
            {{.Assignment}}
        })
    },
    data() {
        return {
            select: {
				{{.SelectData}}
            },
            messages: [],
            entity: {
                {{.EntityDecl}}
            }
        }
    },
    watch: {
		{{.SelectWatch}}
    },
    methods: {
        querySelections(fieldname, endpoint, filter, val) {
            this.select[fieldname].loading = true;
            axios.get("{{.ServiceURL}}/" + endpoint + "?" + filter + "Lk=" + encodeURIComponent(val)).then(response => {
                this.select[fieldname].loading = false;
                this.select[fieldname].items = response.data.entities.map(function(e) {
                    return { text: e[filter], value: e.id };
                });
            });
        },
        save() {
            if (this.id) {
                axios.put("{{.ServiceURL}}/{{.Endpoint}}/" + this.id, this.entity).then(this.saved)
            } else {
                axios.post("{{.ServiceURL}}/{{.Endpoint}}", this.entity).then(this.saved)
            }
		},
		saved(response) {
			this.id = response.data.entity.id
			{{.Assignment}}

			this.messages.push({
				type: "success",
				text: "{{.Name}} saved successfully"
			})
		}
    }
}
</script>  
`)

var tmplEditorTextfield, _ = template.New("EditorTextfield").Parse(`
<v-text-field v-model="entity.{{.Field}}" label="{{.Label}}" />
`)

var tmplEditorCheckbox, _ = template.New("EditorCheckbox").Parse(`
<v-checkbox label="{{.Label}}" v-model="entity.{{.Field}}"></v-checkbox>
`)

var tmplEditorDate, _ = template.New("EditorDate").Parse(`
<v-date-picker v-model="entity.{{.Field}}" label="{{.Label}}" />
`)

var tmplEditorSelectRel, _ = template.New("EditorSelectRel").Parse(`
<v-select
    autocomplete
    cache-items
    required
    label="{{.Label}}"
    :loading="select.{{.Field}}.isloading"
    :items="select.{{.Field}}.items"
    :search-input.sync="select.{{.Field}}.search"
	v-model="entity.{{.Field}}"
></v-select>
`)

//GenerateEditor generate editor code
func GenerateEditor(structInfo generators.StructureInfo) (string, error) {
	var (
		output           bytes.Buffer
		err              error
		fieldsAssignment []string
		fieldsEntityDecl []string
		fieldsMarkup     []string
		selectData       []string
		selectWatch      []string
		data             struct {
			Assignment   string
			EntityDecl   string
			Endpoint     string
			Name         string
			FieldsMarkup string
			SelectData   string
			SelectWatch  string
			ServiceURL   string
		}
	)

	for _, field := range structInfo.Fields {
		var (
			markupSegment bytes.Buffer
		)

		if field.Name != "id" {
			fieldsAssignment = append(fieldsAssignment, "            this.entity."+field.Name+" = response.data."+field.Name)
			fieldsEntityDecl = append(fieldsEntityDecl, field.Name+": null")
		}

		if field.Widget == nil {
			continue
		}

		switch field.Widget.Type {
		default:
			continue
		case "textfield":
			tmplEditorTextfield.Execute(&markupSegment, struct {
				Label string
				Field string
			}{field.Widget.Label, field.Name})
		case "checkbox":
			tmplEditorCheckbox.Execute(&markupSegment, struct {
				Label string
				Field string
			}{field.Widget.Label, field.Name})
		case "date":
			tmplEditorDate.Execute(&markupSegment, struct {
				Label string
				Field string
			}{field.Widget.Label, field.Name})
		case "select-rel":
			// if strings.HasPrefix(field.Widget.Options[0], "")
			tmplEditorSelectRel.Execute(&markupSegment, struct {
				Label string
				Field string
			}{field.Widget.Label, field.Name})

			if len(field.Widget.Options) != 2 {
				return "", fmt.Errorf("invalid options for field: %s", field.Property)
			}

			endpoint := field.Widget.Options[0]
			filtername := field.Widget.Options[1]

			selectData = append(selectData, field.Name+": {search: null, isloading: false, items: []}")
			selectWatch = append(selectWatch, `
			"select.`+field.Name+`.search": function(val) {
				val && this.querySelections("`+field.Name+`", "`+endpoint+`", "`+filtername+`", val)
			}`)
		}
		fieldsMarkup = append(fieldsMarkup, markupSegment.String())
	}

	data.Assignment = strings.Join(fieldsAssignment, "\n")
	data.EntityDecl = strings.Join(fieldsEntityDecl, ",\n")
	data.Endpoint = structInfo.TableName
	data.Name = structInfo.Name
	data.FieldsMarkup = strings.Join(fieldsMarkup, "\n")
	data.SelectData = strings.Join(selectData, ",\n")
	data.SelectWatch = strings.Join(selectWatch, ",\n")
	data.ServiceURL = "http://localhost:8888"

	err = tmplEditor.Execute(&output, data)
	if err == nil {
		o := output.String()
		o = strings.Replace(o, "᚜", "{{", -1)
		o = strings.Replace(o, "᚛", "}}", -1)
		return o, nil
	}

	return "", err
}
