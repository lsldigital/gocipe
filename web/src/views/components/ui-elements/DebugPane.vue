<template>
    <v-card>
        <v-toolbar dense dark color="cyan darken-2">
        <v-btn icon @click="closePane" dark>
            <v-icon>close</v-icon>
        </v-btn>
        <v-toolbar-title>Developement Logger</v-toolbar-title>
        <v-spacer></v-spacer>
        <v-toolbar-items>
            <!-- <v-btn dark flat @click.native="dialog = false">Save</v-btn> -->
        </v-toolbar-items>
        </v-toolbar>
        <v-list two-line>
          <template v-for="(item, index) in getEvents">
            <v-list-tile
              ripple
              :key="item.message"
            >
              <v-list-tile-content>
                <v-list-tile-title>{{ item.message }}</v-list-tile-title>
                <v-list-tile-sub-title class="text--primary">{{ item.type }}</v-list-tile-sub-title>
                <v-list-tile-sub-title>{{ item.message }}</v-list-tile-sub-title>
              </v-list-tile-content>
              <v-list-tile-action>
                <v-list-tile-action-text>{{ item.message }}</v-list-tile-action-text>
                <v-icon
                  color="yellow darken-2"
                >star</v-icon>
              </v-list-tile-action>
            </v-list-tile>
            <v-divider v-if="index + 1 < getEvents.length" :key="index"></v-divider>
          </template>
        </v-list>
    </v-card>
</template>

<script>
import { mapGetters } from "vuex";

export default {
  computed: {
    ...mapGetters("debugpane", ["getEvents"])
  },
  methods: {
    closePane(event) {
      this.$store.dispatch("debugpane/closePane");
    }
  }
};
</script>
