package vuetify

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/fluxynet/gocipe/generators"
	"github.com/jinzhu/inflection"
)

var tmplList, _ = template.New("GenerateList").Parse(`
<template>
    <div class="container">
        <v-container>
            <v-toolbar color="transparent" flat>
                <v-toolbar-title class="grey--text text--darken-4"><h2>{{.Name}} listing</h2></v-toolbar-title>
                <v-spacer></v-spacer>
                <v-btn ml-4 right small fab dark color="info" :to="{name: '{{.Endpoint}}'}">
                    <v-icon dark>add</v-icon>
                </v-btn>
            </v-toolbar>

            <v-text-field mb-4 append-icon="search" label="Search" single-line hide-details v-model="search"></v-text-field>            
            
            <v-alert :type="message.type" :value="true" v-for="(message, index) in messages" :key="index">
                ᚜ message.text ᚛
            </v-alert>

            <v-data-table :headers="headers" :items="entities" class="elevation-1" :search="search">
                <template slot="items" slot-scope="props">
                    {{.ColumnData}}
                    <td class="justify-center layout px-0">
                        <v-btn icon class="mx-0" :to="{name: '{{.Endpoint}}', params: {'id': props.item.id}  }">
                            <v-icon color="teal">edit</v-icon>
                        </v-btn>
                    </td>
                </template>

                <template slot="no-data">
                    <v-flex ma-4>
                        <v-alert slot="no-results" :value="true" color="info" outline icon="info" v-if="search.length > 0">
                        Your search for "᚜ search ᚛" found no results.
                        </v-alert>
                        <v-alert slot="no-results" :value="true" color="info" outline icon="info" v-else>
                            No {{.Name}} found. Would you like to create one?
                            <v-btn :to="{name: '{{.Endpoint}}'}" color="info">create</v-btn>
                        </v-alert>
                    </v-flex>
                </template>
            </v-data-table>
        </v-container>
    </div>
</template>

<script>
import axios from "axios"
export default {
  props: ["id"],
  data() {
    return {
      messages: [],
      search: "",
      headers: [
        {{.ColumnNames}},
        {'text': 'Action', 'value': null}
      ],
      entities: []
    };
  },
  created() {
      axios.get("/api/{{.Endpoint}}").then(response => {
          this.entities = response.data.entities
      })
  }
};
</script>
`)

//GenerateList produces code for a vuetify list component
func GenerateList(structInfo generators.StructureInfo) (string, error) {
	var (
		output bytes.Buffer
		err    error
		data   struct {
			Name        string
			ColumnData  string
			ColumnNames string
			Endpoint    string
		}
		columnData  []string
		columnNames []string
	)

	data.Name = inflection.Plural(structInfo.Name)
	data.Endpoint = structInfo.TableName

	for _, field := range structInfo.Fields {
		if field.Widget == nil {
			continue
		}
		columnData = append(columnData, `                    <td>᚜ props.item.`+field.Name+` ᚛</td>`)
		columnNames = append(columnNames, `{text: "`+field.Widget.Label+`", value: "`+field.Name+`"}`)
	}

	data.ColumnData = strings.Join(columnData, "\n")
	data.ColumnNames = strings.Join(columnNames, ",\n")
	err = tmplList.Execute(&output, data)
	if err == nil {
		o := output.String()
		o = strings.Replace(o, "᚜", "{{", -1)
		o = strings.Replace(o, "᚛", "}}", -1)
		return o, nil
	}

	return "", err
}
