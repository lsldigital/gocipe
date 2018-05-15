<template>
    <v-card color="dark" class="mb-5" height="100%">
        <v-list-tile avatar @click="toggle('http_generate')" class="pa-2">
            <v-list-tile-action>
                <v-checkbox v-model="checkedValues['http_generate']" @click.prevent=""></v-checkbox>
            </v-list-tile-action>
            <v-list-tile-content>
                <v-list-tile-title>HTTP Server</v-list-tile-title>
                <v-list-tile-sub-title>To Generate Http service</v-list-tile-sub-title>
            </v-list-tile-content>

            <v-flex xs8>
                <v-text-field
                id="port_number"
                name="port_number"
                label="Port number"
                ></v-text-field>
            </v-flex>
        </v-list-tile>
    </v-card>
</template>

<script>

export default {
    data () {
      return {
        checkedValues: {
          http_generate: false,
        }
      }
    },
    methods: {
      toggle (name) {
        this.checkedValues[name] = !this.checkedValues[name]
      }
    }
  }
</script>