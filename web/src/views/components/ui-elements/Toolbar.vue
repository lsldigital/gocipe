<template>
  <div>
    <v-toolbar app fixed clipped-left dense dark class="blue-grey darken-4 top-toolbar">
      <v-toolbar-side-icon @click.stop="drawer = !drawer"></v-toolbar-side-icon>
      <v-toolbar-title><router-link class="home-link" :to="{name:'home'}">{{ siteName }}</router-link></v-toolbar-title>

      <v-spacer/>


      <!-- <v-btn light @click="$store.dispatch('save')">
        <v-icon>done</v-icon> SAVE
      </v-btn>
      <v-btn light @click="$store.dispatch('restore')">
        <v-icon>restore</v-icon> RESTORE
      </v-btn> -->
      
      <v-spacer/>

      <!-- Logged in -->
      <template v-if="getAuth">
        <v-btn icon :to="{name: 'settings'}">
          <v-icon>settings</v-icon>
        </v-btn>
        <v-chip  color="success" text-color="white">
          <v-avatar>
            <v-icon>check_circle</v-icon>
          </v-avatar>
          {{ getName }}
        </v-chip>
        <v-menu bottom left offset-y :open-on-hover="true">
          <v-btn icon slot="activator" dark>
            <v-icon>more_vert</v-icon>
          </v-btn>
          <v-list>
            <v-list-tile :key="1" @click="logout()">
              <v-list-tile-title>Logout</v-list-tile-title>
            </v-list-tile>
          </v-list>
        </v-menu>
      </template>
      <!-- Logged Out -->
      <template v-else>
        <v-btn color="warning" text-color="black" @click="login()">Sign in to continue
          <v-icon>lock</v-icon>
        </v-btn>
      </template>
    </v-toolbar>
    <Navigation :drawer="drawer" />
  </div>
</template>

<script>
import { mapGetters, commit, mapMutations, mapActions } from "vuex";
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
