<template>
  <v-card color="dark" class="mb-5" height="100%">
	<v-toolbar color="#212121">
      <i class="material-icons">settings</i>
      <v-toolbar-title>Entities specification</v-toolbar-title>
      <v-spacer></v-spacer>
    </v-toolbar>

            <!-- <v-layout class="borderwrapper" row wrap v-for="(row, index) in bootstrap.settings" v-bind:key="row.id" :v-model="nextid" > -->
         
            <!-- <v-flex xs3>
              <v-subheader>Name</v-subheader>		
            </v-flex>

            <v-flex xs7>
                <v-text-field v-model="row.name" :id="'name'+index" :name="row.name"  label="Name" ></v-text-field>
            </v-flex>

			<v-flex xs3>
              <v-subheader> Type</v-subheader>
            </v-flex>

            <v-flex xs7>
                <v-text-field v-model="row.type" :id="'type'+index" :name="row.type"  label="Type" ></v-text-field>
            </v-flex>
	
            <v-flex xs3>
                <v-subheader> Description</v-subheader>
            </v-flex>

            <v-flex xs7>
                <v-text-field v-model="row.description" :id="'desc'+index" :name="row.description"  label="Description" ></v-text-field>
            </v-flex> -->

        	<!-- </v-layout> -->
  </v-card>
</template>

<script>

    export default {
        data() {
            return {
                entities: [
                    {
                        "name": "",
                        "primary_key": "",
                        "table": "",
                        "table_constraints": [],
                        "description": "",
                        "fields": [
                            {
                                "label": "",
                                "serialized": "",
                                "property": {
                                    "name": "",
                                    "type": ""
                                },
                                "schema": {
                                    "field": "",
                                    "type": ""
                                },
                                "relationship": {
                                    "type": "",
                                    "target": {
                                        "entity":"",
                                        "endpoint":"",
                                        "query":"",
                                        "table":"",
                                        "thisid":"",
                                        "thatid":"",
                                        "thatfield_type":""
                                    }
                                },
                                "widget": {
                                    "type": "",
                                    "options": [
                                        {
                                            "value": "",
                                            "label": ""
                                        }
                                    ],
                                    "target": {
                                        "endpoint": "",
                                        "label": ""
                                    },
                                },
                                "filterable": true
                            }
                        ],
                        "schema": {
                            "create": true,
                            "drop": true,
                            "aggregate": true,
                            "path": ""
                        },
                        "schema": {},
                        "crud": {},
                        "rest": {},
                        "vuetify": {}
                    }
                ],
                entity : []
            }
        },
        methods: {
            toggle(name) {
                if (name == "generate") {
                    this.http.generate = !this.http.generate
                }
            }
        }
    }
</script>