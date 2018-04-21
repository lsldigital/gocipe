<template>
  <div>
    <v-btn block @click.native.stop="dialog = true">Add {{property}}</v-btn>
    <v-dialog v-model="dialog" max-width="290">
      <v-card>
        <v-card-title class="headline">Add new {{property}}</v-card-title>
        <v-card-text>
          <p>{{ types[property].tooltip }}</p>
          <!-- {{model}} -->
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
      let model = Object.assign(this.model);
      this.$emit("interface", model); // handle data and give it back to parent by interface
    }
  },
  created() {
    const definitions = [...this.types[this.property].definition];
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
        } else if (def.type === "array") {
          return {
            [def.json]: []
          };
        }
        return {
          [def.json]: ""
        };
      })
      .reduce((acc, obj) => Object.assign(acc, obj), {});
    // merge an array of object to a single object

    this.model = newDefinitions;
  }
};
</script>
