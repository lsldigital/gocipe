<template>
  <v-flex xs12 pa-3>
    <v-stepper v-model="e1">
      <v-stepper-header>
        <v-stepper-step :complete="e1 > 1" step="1">BOOTSTRAP</v-stepper-step>
        <v-divider></v-divider>
        <v-stepper-step :complete="e1 > 2" step="2">HTTP</v-stepper-step>
        <v-divider></v-divider>
        <v-stepper-step :complete="e1 > 3" step="3">SCHEMA</v-stepper-step>
        <v-divider></v-divider>
        <v-stepper-step :complete="e1 > 4" step="4">CRUD</v-stepper-step>
        <v-divider></v-divider>
        <v-stepper-step :complete="e1 > 5" step="5">REST API</v-stepper-step>
        <v-divider></v-divider>
        <v-stepper-step :complete="e1 > 6" step="6">VUETIFY</v-stepper-step>
        <v-divider></v-divider>
        <v-stepper-step step="7">ENTITIES</v-stepper-step>
        </v-stepper-header>	 	
      <v-stepper-items>
      <v-stepper-content step="1" >
        <page-bootstrap />
        <v-btn color="primary" @click.native="e1 = 2">Continue</v-btn>
        <v-btn flat>Cancel</v-btn>
      </v-stepper-content>
      <v-stepper-content step="2">
        <page-http />
        <v-btn color="primary" @click.native="e1 = 3">Continue</v-btn>
        <v-btn flat>Cancel</v-btn>
      </v-stepper-content>
      <v-stepper-content step="3">
        <page-schema />
        <v-btn color="primary" @click.native="e1 = 4">Continue</v-btn>
        <v-btn flat>Cancel</v-btn>
      </v-stepper-content>
      <v-stepper-content step="4">
        <page-crud />
        <v-btn color="primary" @click.native="e1 = 5">Continue</v-btn>
        <v-btn flat>Cancel</v-btn>
      </v-stepper-content>
      <v-stepper-content step="5">
        <page-rest />
        <v-btn color="primary" @click.native="e1 = 6">Continue</v-btn>
        <v-btn flat>Cancel</v-btn>
      </v-stepper-content>
        <v-stepper-content step="6">
          <page-vuetify />
          <v-btn color="primary" @click.native="e1 = 7">Continue</v-btn>
          <v-btn flat>Cancel</v-btn>
        </v-stepper-content>
        <v-stepper-content step="7">
          <v-card color="dark" class="mb-5" height="100%"></v-card>
          <v-btn color="primary" @click.native="e1 = 1">Continue</v-btn>
          <v-btn flat>Cancel</v-btn>
        </v-stepper-content>
      </v-stepper-items>
    </v-stepper>
  </v-flex>
</template>

<script>
import PageBootstrap from '@/components/subpages/PageBootstrap.vue'
import PageSchema from '@/components/subpages/PageSchema.vue'
import PageHttp from '@/components/subpages/PageHttp.vue'
import PageCrud from '@/components/subpages/PageCrud.vue'
import PageRest from '@/components/subpages/PageRest.vue'
import PageVuetify from '@/components/subpages/PageVuetify.vue'
export default {
    components: {
      PageBootstrap,
      PageSchema,
      PageHttp,
      PageCrud,
      PageRest,
      PageVuetify
    },
    data () {
      return {
        e1: 4
      }
    },
    methods: {
      toggle (name) {
        this.checkedValues[name] = !this.checkedValues[name]
      }
    }
  }
</script>