<template>
    <v-card class="mb-5" color="dark" height="100%">
        <v-list-tile @click="toggle('generate')" avatar="avatar" class="pa-2">
            <v-list-tile-action>
                <v-checkbox v-model="vuetify.generate"></v-checkbox>
            </v-list-tile-action>
            <v-list-tile-content>
                <v-list-tile-title>Generate</v-list-tile-title>
                <v-list-tile-sub-title>Choose to generate back-office in VueJs</v-list-tile-sub-title>
            </v-list-tile-content>

            <v-flex xs8="xs8">
                <v-text-field v-model="vuetify.module" id="module" label="Module" name="module"></v-text-field>
            </v-flex>
        </v-list-tile>

      <v-divider></v-divider> <v-divider></v-divider> <v-divider></v-divider>
      <v-flex xs12 class="text-xs-center">
    		<v-btn @click="pushvuetify" icon color="primary">
      		<i class="material-icons"> check_circle </i>
      	</v-btn>
			</v-flex>
    </v-card>
</template>

<script>
import { mapActions } from 'vuex';
    export default {
        data() {
            return {
                vuetify: {
                    generate: true,
                    module: ""
                }
            }
        },
        methods: {
          ...mapActions(['addvuetify']),
            toggle(name) {
               if (name == "generate") {
                 this.vuetify.generate = !this.vuetify.generate
               }
            },

          pushvuetify() {
    	      this.addvuetify(this.vuetify);
    	    }
        }
    }
</script>