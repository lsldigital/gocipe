<template>
    <v-card class="mb-5" color="dark" height="100%">
        <v-list-tile avatar="avatar" class="pa-2">
            <v-flex xs2>
                <v-checkbox v-model="vuetify.generate" label="Generate"></v-checkbox>
            </v-flex>
            <v-flex xs8="xs8">
                <v-text-field v-model="vuetify.module" id="module" label="Module" name="module"></v-text-field>
            </v-flex>
        </v-list-tile>

        <v-divider></v-divider>
        <v-divider></v-divider>
        <v-divider></v-divider>
        <v-flex xs12 class="text-xs-center">
            <v-btn @click="pushvuetify" icon color="primary">
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
      vuetify: {
        generate: true,
        module: ""
      }
    };
  },
  mounted() {
    if (this.gocipe.vuetify !== undefined) {
      this.vuetify.create =
        this.gocipe.vuetify.generate === undefined
          ? false
          : this.gocipe.vuetify.generate;

      this.vuetify.module =
        this.gocipe.vuetify.module === undefined
          ? ""
          : this.gocipe.vuetify.module;
    }
  },
  methods: {
    ...mapActions(["addvuetify"]),

    pushvuetify() {
      this.addvuetify(this.vuetify);
    }
  }
};
</script>