package vuetify

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
)

var tmplEditor, _ = template.New("GenerateEditor").Parse(`
<template>
    <div class="container">
		<v-toolbar color="transparent" flat>
            <v-toolbar-title class="grey--text text--darken-4 ml-0"><h2>{{.Name}}</h2></v-toolbar-title>
            <v-spacer></v-spacer>
            <v-btn ml-0 small color="grey" flat :to="{name: '{{.Endpoint}}List'}">
                <v-icon dark>arrow_back</v-icon> Back
            </v-btn>
        </v-toolbar>
		<v-alert :type="message.type" :value="true" v-for="(message, index) in messages" :key="index">
		᚜ message.text ᚛
		</v-alert>
  
        {{.FieldsMarkup}}

        <v-btn color="primary" @click="save()">Save</v-btn>
        <v-btn color="gray" :to="{name: '{{.Endpoint}}List'}">Cancel</v-btn>
	</div>
</template>
  
<script>
import axios from "axios"

export default {
    props: ["id"],
    created() {
        if (!this.id) {
            return
        }

        axios.get("/api/{{.Endpoint}}/" + this.id).then(response => {
            this.id = response.data.entity.id
			{{.Assignment}}
			{{.CreatedAssgn}}
        })
    },
    data() {
        return {
            select: {
				{{.SelectData}}
			},
			dates: {
				{{.DateData}}
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
            this.select[fieldname].loading = true
            axios.get("/api/" + endpoint + "?" + filter + "Lk=" + encodeURIComponent(val)).then(response => {
                this.select[fieldname].loading = false
                this.select[fieldname].items = response.data.entities.map(function(e) {
                    return { text: e[filter], value: e.id }
                })
            })
        },
        save() {
            if (this.id) {
                axios.put("/api/{{.Endpoint}}/" + this.id, this.entity).then(this.saved)
            } else {
                axios.post("/api/{{.Endpoint}}", this.entity).then(this.saved)
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

var tmplEditorTextarea, _ = template.New("EditorTextarea").Parse(`
<v-text-field v-model="entity.{{.Field}}" label="{{.Label}}" multiline />
`)

var tmplEditorNumber, _ = template.New("EditorRange").Parse(`
<v-text-field v-model="entity.{{.Field}}" label="{{.Label}}" type="number" />
`)

var tmplEditorPassword, _ = template.New("EditorPassword").Parse(`
<v-text-field
	v-model="entity.{{.Field}}"
	:append-icon="e1 ? 'visibility' : 'visibility_off'"
	:append-icon-cb="() => (e1 = !e1)"
	:type="e1 ? 'password' : 'text'"
	counter
  ></v-text-field>
`)

var tmplEditorCheckbox, _ = template.New("EditorCheckbox").Parse(`
<v-checkbox label="{{.Label}}" v-model="entity.{{.Field}}"></v-checkbox>
`)

var tmplEditorDate, _ = template.New("EditorDate").Parse(`
<v-menu
	ref="menu_{{.Field}}"
	lazy
	:close-on-content-click="false"
	v-model="dates.{{.Field}}.menuAppear"
	transition="scale-transition"
	offset-y
	full-width
	:nudge-right="40"
	min-width="290px"
	:return-value.sync="dates.{{.Field}}.value"
	>
	<v-text-field
		slot="activator"
		label="{{.Label}}"
		v-model="dates.{{.Field}}.value"
		prepend-icon="event"
		readonly
		></v-text-field>
		<v-date-picker v-model="dates.{{.Field}}.value" @change="entity.{{.Field}} = dates.{{.Field}}.value + 'T00:00:00Z'" no-title scrollable>
		<v-spacer></v-spacer>
		<v-btn flat color="primary" @click="menu_{{.Field}} = false">Cancel</v-btn>
		<v-btn flat color="primary" @click="$refs.menu_{{.Field}}.save(dates.{{.Field}}.value)">OK</v-btn>
		</v-date-picker>
</v-menu>
`)

var tmplEditorTime, _ = template.New("EditorTime").Parse(`
<div>
    <v-time-picker v-model="entity.{{.Field}}" label="{{.Label}}" :landscape="landscape"></v-time-picker>
  </div>
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
	{{.Options}}
></v-select>
`)

var tmplEditorSelect, _ = template.New("EditorSelect").Parse(`
<v-select
	autocomplete
	cache-items
	required
	label="{{.Label}}"
	:items="select.{{.Field}}.items"
	v-model="entity.{{.Field}}"
	{{.Options}}
></v-select>
`)

var tmplEditorToggle, _ = template.New("EditorToggle").Parse(`
<v-switch
	label="{{.Label}}"
	v-model="entity.{{.Field}}"
></v-switch>
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
		dateData         []string
		createdAssgn     []string
		data             struct {
			Assignment   string
			EntityDecl   string
			Endpoint     string
			Name         string
			FieldsMarkup string
			SelectData   string
			SelectWatch  string
			DateData     string
			CreatedAssgn string
		}
	)

	for _, field := range structInfo.Fields {
		var (
			markupSegment bytes.Buffer
		)

		if field.Name != "id" {
			fieldsAssignment = append(fieldsAssignment, "            this.entity."+field.Name+" = response.data.entity."+field.Name)
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

		case "textarea":
			tmplEditorTextarea.Execute(&markupSegment, struct {
				Label string
				Field string
			}{field.Widget.Label, field.Name})

		case "number":
			tmplEditorTextarea.Execute(&markupSegment, struct {
				Label string
				Field string
			}{field.Widget.Label, field.Name})

		case "password":
			tmplEditorPassword.Execute(&markupSegment, struct {
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
			dateData = append(dateData, field.Name+": {value: null, menuAppear: false}")
			createdAssgn = append(createdAssgn, "this.dates."+field.Name+".value = response.data.entity."+field.Name+".substr(0,10)")

		case "select":
			var options []string

			for _, v := range field.Widget.Options {
				if v == "multiple" {
					options = append(options, "multiple chips")
				}
			}

			tmplEditorSelect.Execute(&markupSegment, struct {
				Label   string
				Field   string
				Options string
			}{field.Widget.Label, field.Name, strings.Join(options, " ")})

			var values []string
			for _, pair := range field.Widget.Data {
				p := strings.SplitN(pair, ":", 2)
				if len(p) == 1 {
					values = append(values, fmt.Sprintf(`{text: "%s", value: "%s"}`, p[0], p[0]))
				} else {
					values = append(values, fmt.Sprintf(`{text: "%s", value: "%s"}`, p[0], p[1]))
				}
			}

			selectData = append(selectData, field.Name+": { items: ["+strings.Join(values, ", ")+"]}")

		case "select-rel":
			var options []string

			for _, v := range field.Widget.Options {
				if v == "multiple" {
					options = append(options, "multiple chips")
				}
			}

			tmplEditorSelect.Execute(&markupSegment, struct {
				Label   string
				Field   string
				Options string
			}{field.Widget.Label, field.Name, strings.Join(options, " ")})

			if len(field.Widget.Data) != 2 {
				return "", fmt.Errorf("invalid options for field: %s", field.Property)
			}

			endpoint := field.Widget.Data[0]
			filtername := field.Widget.Data[1]

			selectData = append(selectData, field.Name+": {search: null, isloading: false, items: []}")
			selectWatch = append(selectWatch, `
			"select.`+field.Name+`.search": function(val) {
				val && this.querySelections("`+field.Name+`", "`+endpoint+`", "`+filtername+`", val)
			}`)

		case "toggle":
			tmplEditorToggle.Execute(&markupSegment, struct {
				Label string
				Field string
			}{field.Widget.Label, field.Name})
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
	data.DateData = strings.Join(dateData, ",\n")
	data.CreatedAssgn = strings.Join(createdAssgn, ",\n")

	err = tmplEditor.Execute(&output, data)
	if err == nil {
		o := output.String()
		o = strings.Replace(o, "᚜", "{{", -1)
		o = strings.Replace(o, "᚛", "}}", -1)
		return o, nil
	}

	return "", err
}
