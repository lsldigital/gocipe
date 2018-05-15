<template>
    <v-card color="dark" class="mb-5" height="100%">
        <v-list-tile avatar @click="toggle('bootstrap_generate')" class="pa-2">
            <v-list-tile-action>
                <v-checkbox v-model="checkedValues['bootstrap_generate']" @click.prevent=""></v-checkbox>
            </v-list-tile-action>
            <v-list-tile-content>
                <v-list-tile-title>Generate</v-list-tile-title>
                <v-list-tile-sub-title>Generate bootstrap</v-list-tile-sub-title>
            </v-list-tile-content>
        </v-list-tile>

        <v-divider></v-divider> 

        <v-layout row >
            <v-flex xs3>
                <v-subheader> Environment Name</v-subheader>
            </v-flex>

            <v-flex xs8>
                <v-text-field
                id="testing"
                name="input-1"
                label="Setting (Key)"
                ></v-text-field>
            </v-flex>

            <v-flex xs3>
                <v-subheader>Environment Description</v-subheader>
            </v-flex>
            <v-flex xs8>
                <v-text-field
                id="testing"
                name="input-1"
                label="Setting (Value)"
                ></v-text-field>
            </v-flex>
        </v-layout> 

        <v-divider></v-divider> 

        <v-flex xs12>
            <v-layout column align-center class="pa-3">
                <div class="text-xs-center">
                    <div>
                        <v-btn color="primary" dark large >Add Environment Fields</v-btn>
                    </div>
                </div>
            </v-layout>
        </v-flex>      
    </v-card>
</template>

<script>

export default {
    data () {
      return {
        checkedValues: {
          bootstrap_generate: false,
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