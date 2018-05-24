<template>
    <div>
        <v-layout class="borderwrapper" row wrap>

            <v-layout row wrap>
                <v-flex xs10>
                    <v-subheader class="entityheader">Field information</v-subheader>
                </v-flex>
                <v-flex xs3>
                    <v-subheader>Name</v-subheader>
                </v-flex>

                <v-flex xs7>
                    <v-text-field v-model="field.name" label="Name" hint="Label is the label for the field"></v-text-field>
                </v-flex>

                <v-flex xs3>
                    <v-subheader>Serialized</v-subheader>
                </v-flex>

                <v-flex xs7>
                    <v-text-field v-model="field.serialized" label="Serialized" hint="Serialized is the name of the field for serialization (e.g. json)"></v-text-field>
                </v-flex>

                <v-flex xs3>
                    <v-subheader>Filterable</v-subheader>
                </v-flex>

                <v-flex xs7>
                    <v-checkbox v-model="field.filterable"></v-checkbox>
                </v-flex>

                <v-flex xs12>
                    <v-layout column wrap>
                        <v-divider></v-divider>
                    </v-layout>
                </v-flex>

                <v-flex xs10>
                    <v-subheader class="entityheader">Property</v-subheader>
                </v-flex>

                <v-flex xs3>
                    <v-subheader>Name</v-subheader>
                </v-flex>
                <v-flex xs7>
                    <v-text-field v-model="field.property.name" label="Name" hint="Name is the name of the property"></v-text-field>
                </v-flex>
                <v-flex xs3>
                    <v-subheader>Type</v-subheader>
                </v-flex>
                <v-flex xs7>
                    <v-text-field v-model="field.property.type" label="Type" hint="Type is the data type of the property"></v-text-field>
                </v-flex>

                <v-flex xs12>
                    <v-layout column wrap>
                        <v-divider></v-divider>
                    </v-layout>
                </v-flex>

                <v-flex xs10>
                    <v-subheader class="entityheader">Schema</v-subheader>
                </v-flex>

                <v-flex xs3>
                    <v-subheader>Field</v-subheader>
                </v-flex>

                <v-flex xs7>
                    <v-text-field v-model="field.schema.name" label="Field" hint="Field is the name of the field in database"></v-text-field>
                </v-flex>

                <v-flex xs3>
                    <v-subheader>Type</v-subheader>
                </v-flex>

                <v-flex xs7>
                    <v-text-field v-model="field.schema.type" label="Type" hint="Type is the data type for the field in database"></v-text-field>
                </v-flex>
                <v-flex xs3>
                    <v-subheader>Nullable</v-subheader>
                </v-flex>
                <v-flex xs7>
                    <v-checkbox v-model="field.schema.nullable"></v-checkbox>
                </v-flex>

                <v-flex xs3>
                    <v-subheader>Default</v-subheader>
                </v-flex>

                <v-flex xs7>
                    <v-text-field v-model="field.schema.default" label="Default" hint="Default provides the default value for this field in database"></v-text-field>
                </v-flex>

                <v-flex xs12>
                    <v-layout column wrap>
                        <v-divider></v-divider>
                    </v-layout>
                </v-flex>

                <v-flex xs10>
                    <v-subheader class="entityheader">Relationship</v-subheader>
                </v-flex>

                <v-flex xs3>
                    <v-subheader>Type</v-subheader>
                </v-flex>

                <v-flex xs7>
                    <v-text-field v-model="field.relationship.type" label="Type" hint="Type represents the type of relationship"></v-text-field>
                </v-flex>

                <v-flex xs3></v-flex>
                <v-flex xs7>
                    <v-layout column wrap>
                        <v-subheader class="entityheader">Target

                        </v-subheader>
                        <i class="p_italic">Target represents the target of the relationship</i>

                        <v-flex xs7>
                            <v-text-field v-model="field.relationship.target.entity" label="Entity" hint="Entity represents the other entity in the relationship"></v-text-field>
                        </v-flex>

                        <v-flex xs7>
                            <v-text-field v-model="field.relationship.target.endpoint" label="Endpoint" hint="Endpoint represents the endpoint where query can be made"></v-text-field>
                        </v-flex>

                        <v-flex xs7>
                            <v-text-field v-model="field.relationship.target.query" label="Query" hint="Query represents the field to use in the query string"></v-text-field>
                        </v-flex>

                        <v-flex xs7>
                            <v-text-field v-model="field.relationship.target.table" label="Table" hint="">Table represents the other table in the relationship</v-text-field>
                        </v-flex>

                        <v-flex xs7>
                            <v-text-field v-model="field.relationship.target.thisid" label="This ID" hint="ThisID represents the field in this entity used for the relationship"></v-text-field>
                        </v-flex>

                        <v-flex xs7>
                            <v-text-field v-model="field.relationship.target.thatid" label="That Id" hint="ThatID represents the field in the other entity used for the relationship"></v-text-field>
                        </v-flex>

                        <v-flex xs7>
                            <v-text-field v-model="field.relationship.target.thatfield_type" label="That field type" hint="That Field Type represents the field type of the other entity"></v-text-field>
                        </v-flex>
                    </v-layout>
                </v-flex>

                <v-flex xs12>
                    <v-layout column wrap>
                        <v-divider></v-divider>
                    </v-layout>
                </v-flex>

                <v-flex xs10>
                    <v-subheader class="entityheader">Widget</v-subheader>
                </v-flex>
                <v-flex xs3>
                    <v-subheader>Type</v-subheader>
                </v-flex>
                <v-flex xs7>
                    <v-text-field v-model="field.widget.type" label="Type" hint="Type indicates which widget type is represented"></v-text-field>
                </v-flex>

                <v-flex xs3>
                    <v-subheader>Multiple</v-subheader>
                </v-flex>
                <v-flex xs7>
                    <v-checkbox v-model="field.widget.multiple"></v-checkbox>
                </v-flex>

                <v-flex xs3></v-flex>
                <v-flex xs7>
                    <v-layout column wrap>
                        <v-subheader class="entityheader">Widget Target </v-subheader>
                        <i class="p_italic">Widget Target represents a target endpoint to pull data for this widget</i>

                        <v-flex xs7>
                            <v-text-field v-model="field.widget.target.endpoint" label="Endpoint" hint="Endpoint represents an endpoint to pull data from"></v-text-field>
                        </v-flex>

                        <v-flex xs7>
                            <v-text-field v-model="field.widget.target.label" label="Label" hint="Label which field to use for label on data endpoint"></v-text-field>
                        </v-flex>
                    </v-layout>
                </v-flex>

                <v-flex xs12 class="marginbottom">
                    <v-flex xs6 class=" mx-auto mb-5 mt-3">
                        <v-layout column wrap>
                            <div>
                                <v-btn color="secondary" v-on:click="addOptions()" large>Add Widget Options</v-btn>
                            </div>
                            <div v-for="(option, index) in field.widget.options" :key="index">
                                <v-list>
                                    <v-list-tile>
                                        <v-text-field v-model="option.value" label="Value"></v-text-field>
                                    </v-list-tile>
                                    <v-list-tile>
                                        <v-text-field v-model="option.label" label="Label"></v-text-field>
                                    </v-list-tile>
                                    <v-btn icon color="primary" v-on:click="remove(index)" dark>
                                        <i class="material-icons">clear</i>
                                    </v-btn>

                                </v-list>
                            </div>
                        </v-layout>
                    </v-flex>
                </v-flex>
            </v-layout>
        </v-layout>
    </div>
</template>

<script>
export default {
  data() {
    return {
      field: {
        label: "",
        serialized: "",
        property: {
          name: "",
          type: ""
        },
        schema: {
          field: "",
          type: "",
          nullable: false,
          default: ""
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
  mounted() {
    if (typeof this.value === "undefined") {
      this.field.label = "";
      this.field.serialized = "";
      this.field.filterable = false;

      this.field.property.name = "";
      this.field.property.type = "";

      this.field.schema.field = "";
      this.field.schema.type = "";
      this.field.schema.nullable = false;
      this.field.schema.default = "";

      this.field.relationship.type = "";
      this.field.relationship.target.entity = "";
      this.field.relationship.target.endpoint = "";
      this.field.relationship.target.query = "";
      this.field.relationship.target.table = "";
      this.field.relationship.target.thisid = "";
      this.field.relationship.target.thatid = "";
      this.field.relationship.target.thatfield_type = "";

      this.field.widget.type = "";
      this.field.widget.options = [];
      this.field.widget.target.endpoint = "";
      this.field.widget.target.label = "";
    }
  },
  methods: {
    addOptions() {
      const opt = {
        value: "",
        label: ""
      };
      this.field.widget.options.push(opt);
    },
    addfield() {
      this.$emit("input", this.field);
    },
    remove(index) {
      this.field.widget.options.splice(index, 1);
    }
  }
};
</script>

<style>
</style>
