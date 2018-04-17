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
    <v-container>
        <v-toolbar color="transparent" flat>
            <v-toolbar-title class="grey--text text--darken-4 ml-0"><h2>{{.Name}}</h2></v-toolbar-title>
            <v-spacer></v-spacer>
            <v-btn mr-0 color="primary" :to="{name: '{{.Endpoint}}Edit', params:{id: 0}}">
                <v-icon dark>add</v-icon> Add
            </v-btn>
        </v-toolbar>
        
        <v-alert :type="message.type === 'E' ? 'error' : message.type" :value="true" v-for="(message, index) in messages" :key="index">
            ᚜ message.text ᚛
        </v-alert>

        <v-alert type="info" value="true"  color="primary" outline icon="info" v-if="entities.length === 0">
            No {{.Name}} exist. Would you like to create one now?
            <v-btn :to="{name: '{{.Endpoint}}Edit', params:{id: 0}}" color="primary">create new</v-btn>
        </v-alert>
        <template v-else>
            <v-text-field mb-4 append-icon="search" label="Search" single-line hide-details v-model="search"></v-text-field>            
            <v-data-table :headers="headers" :items="entities" class="elevation-1" :search="search">
                <template slot="items" slot-scope="props">
                    {{.ColumnData}}
                    <td class="justify-center layout px-0">
                        <v-btn icon class="mx-0" :to="{name: '{{.Endpoint}}Edit', params: {'id': props.item.id}  }">
                            <v-icon color="teal">edit</v-icon>
                        </v-btn>
                    </td>
                </template>

                <template slot="no-data">
                    <v-flex ma-4>
                        <v-alert slot="no-results" :value="true" color="primary" outline icon="info" v-if="search.length > 0">
                        Your search for "᚜ search ᚛" found no results.
                        </v-alert>
                        <v-alert slot="no-results" :value="true" color="primary" outline icon="info" v-else>
                            No {{.Name}} found.
                        </v-alert>
                    </v-flex>
                </template>
            </v-data-table>
        </template>
    </v-container>
</template>

<script>
import axios from "axios"
export default {
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
    axios
      .get("/api/{{.Endpoint}}")
      .then(response => {
        this.entities = response.data.entities;
      })
      .catch(error => {
        this.messages = [...this.messages, ...error.response.data.messages];
      });
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
