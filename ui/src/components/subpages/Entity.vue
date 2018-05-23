<template>
  <v-layout row justify-center>
    <v-dialog v-model="dialog" persistent max-width="900px">
      <v-btn slot="activator" color="primary" dark>ADD DETAILS</v-btn>
      <v-card>
        <v-card-title>
          <span class="headline">Entity details</span>
          <v-divider></v-divider>
        </v-card-title>
        <v-card-text>
          <v-container grid-list-md>
            <v-layout wrap>

              <v-flex xs10>
                <v-subheader> Fields</v-subheader>
              </v-flex>

              <v-layout v-for="(field, index) in fields" :key="'fields'+ index">
                <v-flex xs8>
                  <v-text-field :label="'Field '+ ++index" :id="'const'+index" v-model="fields[index]"></v-text-field>
                </v-flex>

                <v-flex xs1 class="text-xs-center">
                  <v-btn color="primary" dark persistent @click.stop="fielddialogue = !fielddialogue">
                    <i class="material-icons"> add </i>
                    <big> details </big>
                  </v-btn>
                </v-flex>

              </v-layout>

              <v-dialog v-model="fielddialogue" max-width="700px">
                <v-card>
                  <v-card-title>
                    <span class="headline">Field details</span>
                  </v-card-title>
                  <v-card-text>
                    <v-container grid-list-md>
                      <v-layout wrap>
                        <v-flex xs12 sm6 md4>
                          <v-text-field label="Label" hint="Label is the label of the field"></v-text-field>
                        </v-flex>
                        <v-flex xs12 sm6 md4>
                          <v-text-field label="Serialized" hint="Serialized is the name of the field for serialization (e.g. json)"></v-text-field>
                        </v-flex>
                        <v-flex xs12 sm6 md4>
                          <v-list-tile avatar @click="toggle('filterable')" class="pa-2">
                            <v-list-tile-action>
                              <v-checkbox v-model="filterable"></v-checkbox>
                            </v-list-tile-action>
                            <v-list-tile-content>
                              <v-list-tile-title>Filterable</v-list-tile-title>
                              <v-list-tile-sub-title>Field filterable</v-list-tile-sub-title>
                            </v-list-tile-content>
                          </v-list-tile>
                        </v-flex>
                        <v-flex xs12>
                          <v-subheader> Property</v-subheader>
                        </v-flex>
                        <v-flex xs3 class="text-xs-center">
                          <v-btn color="primary" dark persistent @click.stop="propertydialogue = !propertydialogue">
                            <i class="material-icons"> add </i>
                            <big> property </big>
                          </v-btn>
                        </v-flex>

                        <v-flex xs12>
                          <v-subheader> Schema</v-subheader>
                        </v-flex>
                        <v-flex xs3 class="text-xs-center">
                          <v-btn color="primary" dark persistent @click.stop="schemadialogue = !schemadialogue">
                            <i class="material-icons"> add </i>
                            <big> schema </big>
                          </v-btn>
                        </v-flex>

                        <v-flex xs12>
                          <v-subheader> Relationship</v-subheader>
                        </v-flex>
                        <v-flex xs3 class="text-xs-center">
                          <v-btn color="primary" dark persistent @click.stop="relationdialogue = !relationdialogue">
                            <i class="material-icons"> add </i>
                            <big> relation </big>
                          </v-btn>
                        </v-flex>
                      </v-layout>
                    </v-container>

                  </v-card-text>
                  <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn color="blue darken-1" flat @click.native="fielddialogue = false">Close</v-btn>
                    <v-btn color="blue darken-1" flat @click.native="fielddialogue = false">Save</v-btn>
                  </v-card-actions>
                </v-card>
              </v-dialog>

              <v-dialog v-model="propertydialogue" max-width="700px">
                <v-card>
                  <v-card-title>
                    <span class="headline">Add new property</span>
                  </v-card-title>
                  <v-card-text>
                    <v-container grid-list-md>
                      <v-layout wrap>
                        <v-flex xs12 sm6 md6>
                          <v-text-field label="Field" hint="Name of the property"></v-text-field>
                        </v-flex>
                        <v-flex xs12 sm6 md6>
                          <v-text-field label="Type" hint="Type is the data type of the property"></v-text-field>
                        </v-flex>
                      </v-layout>
                    </v-container>

                  </v-card-text>
                  <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn color="blue darken-1" flat @click.native="propertydialogue = false">Close</v-btn>
                    <v-btn color="blue darken-1" flat @click.native="propertydialogue = false">Save</v-btn>
                  </v-card-actions>
                </v-card>
              </v-dialog>

              <v-dialog v-model="schemadialogue" max-width="700px">
                <v-card>
                  <v-card-title>
                    <span class="headline">Add Schema Details</span>
                  </v-card-title>
                  <v-card-text>
                    <v-container grid-list-md>
                      <v-layout wrap>
                        <v-flex xs12>
                          <v-text-field label="Field" hint="Field is the name of the field in database"></v-text-field>
                        </v-flex>
                        <v-flex xs12>
                          <v-text-field label="Type" hint="Type is the data type for the field in database"></v-text-field>
                        </v-flex>
                        <v-flex xs12>
                          <v-list-tile avatar @click="toggle('nullable')" class="pa-2">
                            <v-list-tile-action>
                              <v-checkbox v-model="nullable"></v-checkbox>
                            </v-list-tile-action>
                            <v-list-tile-content>
                              <v-list-tile-title>Nullable</v-list-tile-title>
                            </v-list-tile-content>
                          </v-list-tile>

                          <v-text-field label="Default" hint="Default provides the default value for this field in database"></v-text-field>
                        </v-flex>

                      </v-layout>
                    </v-container>

                  </v-card-text>
                  <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn color="blue darken-1" flat @click.native="schemadialogue = false">Close</v-btn>
                    <v-btn color="blue darken-1" flat @click.native="schemadialogue = false">Save</v-btn>
                  </v-card-actions>
                </v-card>
              </v-dialog>

              <v-dialog v-model="relationdialogue" max-width="700px">
                <v-card>
                  <v-card-title>
                    <span class="headline">Add Relationship Details</span>
                  </v-card-title>
                  <v-card-text>
                    <v-container grid-list-md>
                      <v-layout wrap>
                        <v-flex xs12>
                          <v-text-field label="Type" hint="Type is the data type for the field in database"></v-text-field>
                        </v-flex>

                        <v-flex xs3 class="text-xs-center">
                          <v-btn color="primary" dark persistent @click.stop="relationtarget = !relationtarget">
                            <i class="material-icons"> add </i>
                            <big> Target relationship </big>
                          </v-btn>
                        </v-flex>

                      </v-layout>
                    </v-container>

                  </v-card-text>
                  <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn color="blue darken-1" flat @click.native="relationdialogue = false">Close</v-btn>
                    <v-btn color="blue darken-1" flat @click.native="relationdialogue = false">Save</v-btn>
                  </v-card-actions>
                </v-card>
              </v-dialog>

              <v-dialog v-model="relationtarget" max-width="700px">
                <v-card>
                  <v-card-title>
                    <span class="headline">Add Field Relationship</span>
                  </v-card-title>
                  <v-card-text>
                    <v-container grid-list-md>
                      <v-layout wrap>
                        <v-flex xs12>
                          <v-text-field label="Entity" hint="Entity represents the other entity in the relationship"></v-text-field>
                        </v-flex>

                        <v-flex xs12>
                          <v-text-field label="Endpoint" hint="Endpoint represents the endpoint where query can be made"></v-text-field>
                        </v-flex>

                        <v-flex xs12>
                          <v-text-field label="Query" hint="Query represents the field to use in the query string"></v-text-field>
                        </v-flex>

                        <v-flex xs12>
                          <v-text-field label="Table" hint="Table represents the other table in the relationship"></v-text-field>
                        </v-flex>

                        <v-flex xs12>
                          <v-text-field label="This ID" hint="This ID represents the field in this entity used for the relationship"></v-text-field>
                        </v-flex>

                        <v-flex xs12>
                          <v-text-field label="That ID" hint="That ID represents the field in the other entity used for the relationship"></v-text-field>
                        </v-flex>

                        <v-flex xs12>
                          <v-text-field label="That Field Type" hint="That Field Type represents the field type of the other entity"></v-text-field>
                        </v-flex>

                      </v-layout>
                    </v-container>

                  </v-card-text>
                  <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn color="blue darken-1" flat @click.native="relationtarget = false">Close</v-btn>
                    <v-btn color="blue darken-1" flat @click.native="relationtarget = false">Save</v-btn>
                  </v-card-actions>
                </v-card>
              </v-dialog>

              <v-flex xs12>
                <v-layout column align-center class="pa-3 ">
                  <div class="text-xs-center ">
                    <div>
                      <v-btn color="primary " v-on:click="addfields" dark large>Add Fields</v-btn>
                    </div>
                  </div>
                </v-layout>
              </v-flex>
            </v-layout>
          </v-container>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1 " flat @click.native="dialog=false ">Close</v-btn>
          <v-btn color="blue darken-1 " flat @click.native="dialog=false ">Save</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-layout>
</template>

<script>
export default {
  props: {
    entity: {
      name: "",
      primary_key: "",
      table: "",
      table_constraints: [],
      description: "",
      fields: [],
      schema: {
        create: true,
        drop: true,
        aggregate: true,
        path: ""
      },
      schema: {},
      crud: {},
      rest: {},
      vuetify: {}
    }
  },
  data() {
    return {
      dialog: false,
      fielddialogue: false,
      propertydialogue: false,
      schemadialogue: false,
      relationdialogue: false,
      relationtarget: false,
      filterable: false,
      nullable: false,
      fields: [],
      field: {
        label: "",
        serialized: "",
        property: {
          name: "",
          type: ""
        },
        schema: {
          field: "",
          type: ""
        },
        relationship: {
          type: "",
          target: {
            entity: "",
            endpoint: "",
            query: "",
            table: "",
            thisid: "",
            thatid: "",
            thatfield_type: ""
          }
        },
        widget: {
          type: "",
          options: [
            {
              value: "",
              label: ""
            }
          ],
          target: {
            endpoint: "",
            label: ""
          }
        },
        filterable: true
      }
    };
  },
  methods: {
    addfields() {
      this.fields.push("");
    },

    toggle(name) {
      if (name == "filterable") {
        this.filterable = !this.filterable;
      }
    }
  }
};
</script>