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

    </v-card>
</template>

<script>

    export default {
        data() {
            return {
                http: {
                    "generate": false,
                    "prefix": "",
                    "port": ""
                }
            }
        },
        methods: {
            toggle(name) {
                if (name == "generate") {
                    this.http.generate = !this.http.generate
                }
            }
        }
    }
</script>