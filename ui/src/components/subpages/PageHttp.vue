<template>
    <v-card color="dark" class="mb-5" height="100%">
        <v-list-tile avatar @click="toggle('http.generate')" class="pa-2">
            <v-list-tile-action>
                <v-checkbox v-model="http.generate" @click.prevent=""></v-checkbox>
            </v-list-tile-action>
            <v-list-tile-content>
                <v-list-tile-title>HTTP Server</v-list-tile-title>
                <v-list-tile-sub-title>To Generate Http service</v-list-tile-sub-title>
            </v-list-tile-content>

            <v-flex xs4 class="mr-5">
                <v-text-field id="port_number" v-model="http.port" name="port_number" label="Port number" ></v-text-field>
            </v-flex>

            <v-flex xs4>
                <v-text-field id="prefix" v-model="http.prefix" name="prefix" label="Prefix" ></v-text-field>
            </v-flex>

        </v-list-tile>

    </v-card>
</template>

<script>

export default {
    data () {
      return {
        http: {
            "generate": false,
            "prefix": "",
            "port": ""
        }
      }
    },
    methods: {
      toggle (name) {

				if (this.http.generate == true) {
					this.http.generate =  false
				} else {
					this.http.generate =  true
				}
      }
    }
  }
</script>