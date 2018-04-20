<template>
    <v-navigation-drawer v-bind:value="drawer" clipped fixed app>
        <v-list dense>
            <v-list-tile :to="{ name: 'settings'}">
                <v-list-tile-action>
                    <v-icon>settings</v-icon>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Settings</v-list-tile-title>
                </v-list-tile-content>
            </v-list-tile>
            <v-list-tile :to="{ name: 'gocipe'}">
                <v-list-tile-action>
                    <v-icon>settings</v-icon>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Gocipe</v-list-tile-title>
                </v-list-tile-content>
            </v-list-tile>
        </v-list>

        <v-divider></v-divider>
    </v-navigation-drawer>
</template>

<script>
// import TopicsEdit from "@/views/Entities/TopicsEdit";
import { mapGetters } from "vuex";

export default {
  data() {
    return {
      routes: []
    };
  },
  props: ["drawer"],
  computed: {},
  created() {},
  mounted() {}
};
</script>
