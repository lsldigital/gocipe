<template>
  <div>
    <v-btn block @click.native.stop="dialog = true">Add {{property}}</v-btn>
    <v-dialog v-model="dialog" max-width="290">
      <v-card>
        <v-card-title class="headline">Add new {{property}}</v-card-title>
        <v-card-text>
          <!-- {{  }} -->
          <template v-if="Object.keys(types).indexOf(property) !== -1">
            <div v-for="(field, index) in types[property].definition" :key="index">
              <template v-if="field.type === 'string'">
                <v-text-field :label="field.label" :hint="field.tooltip" v-model="model[field.json]"></v-text-field>
              </template>
              <template v-else-if="field.type === 'bool'">
                <v-checkbox :label="field.label" v-model="model[field.json]"></v-checkbox>
              </template>
              <template v-else>
                prolly an array of type {{field.type}}
              </template>
            </div>
          </template>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="darken-1" flat="flat" @click.native="dialog = false">Cancel</v-btn>
          <v-btn color="green darken-1" flat="flat" @click.native="handleDataFc">Confirm</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>


<script>
import types from "./../../types.js";

export default {
  data() {
    return {
      dialog: false,
      types,
      model: []
    };
  },
  props: ["property"],
  computed: {},
  methods: {
    handleDataFc() {
      this.dialog = false;
      let model = Object.assign({}, this.model);
      this.$emit("interface", model); // handle data and give it back to parent by interface
    }
  },
  created() {
    let definedTypes = Object.keys(types);
    let definitions = [];
    console.log(this.property);
    if (definedTypes.indexOf(this.property) !== -1) {
      console.info("Found your property", this.property);
      definitions = this.types[this.property].definition;
      console.log(definitions);
      const newDefinitions = definitions
        .map(def => {
          if (def.type === "string") {
            return {
              [def.json]: ""
            };
          } else if (def.type === "bool") {
            return {
              [def.json]: false
            };
          } else if (definedTypes.indexOf(def.json) !== -1) {
            console.log("Found your type", def.type);
            return {
              [def.json]: []
            };
          } else if (definedTypes.indexOf(def.json) === -1) {
            console.error("Didn't find your type ", def.type);
          }
          return {
            [def.json]: ""
          };
        })
        .reduce((acc, obj) => Object.assign(acc, obj), {});
      // merge an array of object to a single object

      this.model = newDefinitions;
    } else {
      console.error("Type", this.property, "is not defined");
      return;
    }
  }
};
</script>
