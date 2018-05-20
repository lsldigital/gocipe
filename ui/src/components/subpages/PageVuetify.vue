<template>
    <v-card class="mb-5" color="dark" height="100%">
        <v-list-tile @click="toggle('vuetify_generate')" avatar="avatar" class="pa-2">
            <v-list-tile-action>
                <v-checkbox v-model="checkedValues['vuetify_generate']"></v-checkbox>
            </v-list-tile-action>
            <v-list-tile-content>
                <v-list-tile-title>Generate</v-list-tile-title>
                <v-list-tile-sub-title>Choose to generate back-office in VueJs</v-list-tile-sub-title>
            </v-list-tile-content>

            <v-flex xs8="xs8">
                <v-text-field id="module" label="Module" name="module"></v-text-field>
            </v-flex>
        </v-list-tile>
    </v-card>
</template>

<script>

    export default {
        data() {
            return {
                checkedValues: {
                    vuetify_generate: false
                }
            }
        },
        methods: {
            toggle(name) {
                this.checkedValues[name] = !this.checkedValues[name]
            }
        }
    }
</script>