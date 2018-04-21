<template>
  <div>
    <v-toolbar app fixed clipped-left dense dark class="blue-grey darken-4 top-toolbar">
      <v-toolbar-side-icon @click.stop="drawer = !drawer"></v-toolbar-side-icon>
      <v-toolbar-title>
        <router-link class="home-link" :to="{name:'home'}">{{ siteName }}</router-link>
      </v-toolbar-title>

      <v-spacer/>

      <v-btn outline small @click="$store.dispatch('restore')">
        <v-icon>file_download</v-icon> DOWNLOAD
      </v-btn>

      <v-spacer/>

      <v-btn outline small @click="$store.dispatch('restore')">
        <v-icon>restore</v-icon> RESTORE
      </v-btn>
      <v-btn outline small @click="$store.dispatch('save')">
        <v-icon>done</v-icon> SAVE
      </v-btn>

    </v-toolbar>
    <Navigation :drawer="drawer" />
  </div>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
import Navigation from "./Navigation";

export default {
  data() {
    return {
      drawer: true
    };
  },
  computed: {
    ...mapGetters(["settings", "user"]),
    ...mapGetters("user", ["getName", "getAuth"]),
    siteName() {
      if (this.settings.siteName.value !== "") {
        return this.settings.siteName.value;
      } else {
        return "Default";
      }
    }
  },
  methods: {
    ...mapActions("user", ["login", "logout"])
  },
  components: {
    Navigation
  }
};
</script>

<style lang="scss" scoped>
.home-link {
  text-decoration: none;
  color: white;
}
</style>
