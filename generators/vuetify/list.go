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
            <v-flex mb-4>
                <h1>{{.Name}} listing</h1>
                <v-text-field mb-4 append-icon="search" label="Search" single-line hide-details v-model="search"></v-text-field>
            </v-flex>
            <v-alert :type="message.type" :value="true" v-for="(message, index) in messages" :key="index">
                ᚜ message.text ᚛
            </v-alert>

            <v-data-table :headers="headers" :items="entities" class="elevation-1" :search="search">
                <template slot="items" slot-scope="props">
                    {{.ColumnData}}
                    <td class="">᚜ props.item.name ᚛</td>
                    <td class="justify-center layout px-0">
                        <v-btn icon class="mx-0" :to="{name: '{{.Endpoint}}', params: {'id': props.item.id}  }">
                            <v-icon color="teal">edit</v-icon>
                        </v-btn>
                    </td>
                </template>

                <template slot="no-data">
                    <v-alert slot="no-results" :value="true" color="error" icon="warning">
                    Your search for "᚜ search ᚛" found no results.
                    </v-alert>
                </template>
            </v-data-table>
        </v-container>
    </div>
</template>

<script>
export default {
  props: ["id"],
  data() {
    return {
      messages: [],
      search: "",
      headers: [
        {{.ColumnNames}}
        {text: "name", value: "name"}
      ],
      entities: []
    };
  },
  created() {
      this.axios.get("/api/{{.Endpoint}}").then(response => {
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
		columnData = append(columnData, `                    <td>᚜ props.item.`+field.Name+` ᚛</td>`)
		columnNames = append(columnNames, `{text: "`+field.Name+`", value: "`+field.Name+`"}`)
	}

	err = tmplList.Execute(&output, data)
	if err == nil {
		o := output.String()
		o = strings.Replace(o, "᚜", "{{", -1)
		o = strings.Replace(o, "᚛", "}}", -1)
		return o, nil
	}

	data.ColumnData = strings.Join(columnData, "\n")
	data.ColumnNames = strings.Join(columnNames, ",\n")

	return "", err
}
