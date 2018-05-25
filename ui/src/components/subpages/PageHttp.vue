<template>
    <v-card class="mb-5" color="dark" height="100%">
        <v-list-tile @click="toggle('generate')" avatar="avatar" class="pa-2">
            <v-list-tile-action>
                <v-checkbox v-model="http.generate"></v-checkbox>
            </v-list-tile-action>
            <v-list-tile-content>
                <v-list-tile-title>HTTP Server</v-list-tile-title>
                <v-list-tile-sub-title>To Generate Http service</v-list-tile-sub-title>
            </v-list-tile-content>

            <v-flex class="mr-5" xs4="xs4">
                <v-text-field id="port_number" label="Port number" name="port_number" v-model="http.port"></v-text-field>
            </v-flex>

            <v-flex xs4="xs4">
                <v-text-field id="prefix" label="Prefix" name="prefix" v-model="http.prefix"></v-text-field>
            </v-flex>

        </v-list-tile>
        <v-divider></v-divider>
        <v-divider></v-divider>
        <v-divider></v-divider>
        <v-flex xs12 class="text-xs-center">
            <v-btn @click="pushhttp" icon color="primary">
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
      http: {
        generate: false,
        prefix: "",
        port: ""
      }
    };
  },
  mounted() {
    if (this.gocipe.http !== undefined) {
      this.http.generate =
        this.gocipe.http.generate === undefined
          ? ""
          : this.gocipe.http.generate;

      this.http.prefix =
        this.gocipe.http.prefix === undefined ? "" : this.gocipe.http.prefix;
      this.http.port =
        this.gocipe.http.port === undefined ? "" : this.gocipe.http.port;
    }
  },
  methods: {
    ...mapActions(["addhttp"]),
    toggle(name) {
      if (name == "generate") {
        this.http.generate = !this.http.generate;
      }
    },
    pushhttp() {
      this.addhttp(this.http);
    }
  }
};
</script>