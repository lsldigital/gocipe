<template>
  <v-card color="dark" class="mb-5" height="100%">
    <v-list two-line subheader>
      <v-list-tile avatar @click="toggle('create')" class="pa-2">
        <v-list-tile-action>
          <v-checkbox v-model="schema.create"></v-checkbox>
        </v-list-tile-action>
        <v-list-tile-content>
          <v-list-tile-title>Create</v-list-tile-title>
          <v-list-tile-sub-title>Allow create function for generating schema</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>
      <v-divider></v-divider>
      <v-list-tile avatar @click="toggle('drop')" class="pa-2">
        <v-list-tile-action>
          <v-checkbox v-model="schema.drop"></v-checkbox>
        </v-list-tile-action>
        <v-list-tile-content>
          <v-list-tile-title>Drop</v-list-tile-title>
          <v-list-tile-sub-title>Allow Drop function for generating schema</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>
      <v-divider></v-divider>
      <v-list-tile avatar @click="toggle('aggregate')" class="pa-2">
        <v-list-tile-action>
          <v-checkbox v-model="schema.aggregate"></v-checkbox>
        </v-list-tile-action>
        <v-list-tile-content>
          <v-list-tile-title>Aggregate</v-list-tile-title>
          <v-list-tile-sub-title>Allow Aggregate function for generating schema</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>
    </v-list>

    <v-divider></v-divider>
    <v-divider></v-divider>
    <v-divider></v-divider>
    <v-flex xs12 class="text-xs-center">
      <v-btn @click="pushschema" icon color="primary">
        <i class="material-icons"> check_circle </i>
      </v-btn>
    </v-flex>
  </v-card>
</template>

<script>
import gocipe from "../../assets/gocipe.json";
import { mapActions } from "vuex";
export default {
  data() {
    return {
      gocipe,
      schema: {
        create: true,
        drop: true,
        aggregate: true
      }
    };
  },
  mounted() {
    if (this.gocipe.schema !== undefined) {
      this.schema.create =
        this.gocipe.schema.create === undefined
          ? true
          : this.gocipe.schema.create;

      this.schema.drop =
        this.gocipe.schema.drop === undefined ? true : this.gocipe.schema.drop;

      this.schema.aggregate =
        this.gocipe.schema.aggregate === undefined
          ? true
          : this.gocipe.schema.aggregate;
    }
  },
  methods: {
    ...mapActions(["addschema"]),
    toggle(name) {
      if (name == "create") {
        this.schema.create = !this.schema.create;
      }

      if (name == "drop") {
        this.schema.drop = !this.schema.drop;
      }

      if (name == "aggregate") {
        this.schema.aggregate = !this.schema.aggregate;
      }
    },

    pushschema() {
      this.addschema(this.schema);
    }
  }
};
</script>