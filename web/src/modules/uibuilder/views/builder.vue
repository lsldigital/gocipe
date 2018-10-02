<template>
  <v-layout>
    <v-toolbar app fixed clipped-left dense dark class="blue-grey darken-4 top-toolbar">
      <v-toolbar-side-icon @click.stop="drawer = !drawer"></v-toolbar-side-icon>
      <v-toolbar-title>
        JSON Builder
      </v-toolbar-title>

      <v-spacer/>

      <v-btn outline small :href="`data:${download_link}`" download="data.json">
        <v-icon>file_download</v-icon> DOWNLOAD
      </v-btn>

      <v-spacer/>

      <v-btn outline small @click="'restore'">
        <v-icon>restore</v-icon> RESTORE
      </v-btn>
      <v-btn outline small @click="'save'">
        <v-icon>done</v-icon> SAVE
      </v-btn>

    </v-toolbar>
    <Navigation :drawer="drawer" />

    <v-flex sm12>
      <component :is="getType(gocipe)" :value="gocipe"> </component>
    </v-flex>

  </v-layout>
</template>

<script>
import Toolbar from "../../../views/components/ui-elements/Toolbar.vue";
import Navigation from "../../../views/components/ui-elements/Navigation.vue";
import draggable from "vuedraggable";
// import gocipe from "../gocipe.json";
import gocipe from "../simple.gocipe.json";
// import gocipe from "../types.json";
import booleanType from "./components/booleanType.vue";
import stringType from "./components/stringType.vue";

export default {
  data() {
    return {
      gocipe,
      drawer: true
    };
  },
  computed: {
    download_link() {
      return (
        "text/json;charset=utf-8," +
        encodeURIComponent(JSON.stringify(this.gocipe))
      );
    }
  },
  methods: {
    getType: value => {
      let TYPEOF = typeof value;
      return TYPEOF + "Type";
    },
    save() {
      localStorage.setItem("gocipe", this.gocipe);
    },
    restore() {
      this.gocipe = localStorage.getItem("gocipe");
    }
  },
  components: {
    Toolbar,
    Navigation,
    draggable,
    booleanType,
    stringType
  }
};
</script>


<style lang="scss" scoped>
.code {
  background: red;
}
</style>




