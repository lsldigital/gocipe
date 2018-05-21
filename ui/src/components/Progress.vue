<template>
  <v-flex xs12 pa-3>
    <v-stepper v-model="e1">
      <v-stepper-header>
       <template  v-for="(step,key, index) in steps" >
         <div v-if="index < 6">
          <v-stepper-step :complete="e1 > index + 1" :step="index + 1">{{step}}</v-stepper-step>
          <v-divider></v-divider>
          </div>
          <div v-else>
          <v-stepper-step :step="index + 1">{{step}}</v-stepper-step>
          </div>
       </template>
      </v-stepper-header>

    <v-stepper-items>
     
      <v-stepper-content v-for="(step, key, index) in steps"  :step="key" >

          <component :is="'page-'+steps[key]"></component>
          <div v-if= "index < 6">
          <v-btn color="primary"  v-on:click="setback" @click.native="e1 = index+2">Continue</v-btn>
          
          </div>

          <div v-else>
          <v-btn color="primary" @click.native="e1 = 1">Continue</v-btn>
          <!-- <v-btn @click.native="b1 = 0" flat>Back</v-btn> -->
          </div>

          <div v-if= "index < 1">
            <v-btn @click.native="b1" v-on:click="setlastpage" flat>Back</v-btn>
          </div>

          <div v-else>
            <v-btn @click.native="b1"  v-on:click="setpage" flat>Back</v-btn>
          </div>

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
import PageEntities from '@/components/subpages/PageEntities.vue'

export default {
    components: {
      PageBootstrap,
      PageSchema,
      PageHttp,
      PageCrud,
      PageRest,
      PageVuetify,
      PageEntities
    },
    data () {
      return {
        e1: 0,
        b1: 0,
        steps: {1: "bootstrap", 2 : "http", 3: "schema", 4: "crud", 5: "rest", 6 : "vuetify", 7 : "entities"}
      }
    },
    methods: {
      setback() {
        this.b1 = this.e1-1
      },
      setpage() {
        this.e1 = this.e1-1
      },
      setlastpage() {
        this.e1 = 7
        this.b1 = this.e1-1
      }
      // ...mapActions(['addbootstrap']),
      // ...mapActions(['addhttp']),
      // addtostore: function (event) {
      //   if (this.e1 == 0) {
      //     this.addbootstrap(this.boostrap);
      //   } else if (this.e1 == 1) {
      //     this.addhttp(this.addhttp)
      //   }
      // }
    }
}
</script>

