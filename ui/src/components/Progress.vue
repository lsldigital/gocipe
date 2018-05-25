<template>
  <v-flex xs12 pa-3>
    <v-stepper v-model="e1">
      <v-stepper-header>
        <div v-for="(step, key, index) in steps" :key="key">
          <template v-if="index < 6">
            <v-stepper-step :complete="e1 > index + 1" :step="index + 1">
              <a class="link-component" v-on:click="loadpage(key)"> {{step}} </a>
            </v-stepper-step>

          </template>
          <template v-else>
            <v-stepper-step :step="index + 1">
              <a class="link-component" v-on:click="loadpage(key)"> {{step}} </a>
            </v-stepper-step>
          </template>
        </div>
      </v-stepper-header>

      <v-stepper-items>

        <v-stepper-content v-for="(step, key, index) in steps" :key="key" :step="key">
          <component :is="'page-'+steps[key]"></component>
          <div v-if="index < 6">
            <v-btn color="primary" v-on:click="setback" @click.native="e1 = index+2">Continue</v-btn>
          </div>

          <div v-else>
            <v-btn color="primary" @click.native="e1 = 1">Continue</v-btn>
          </div>

          <div v-if="index < 1">
            <v-btn @click.native="b1" v-on:click="setlastpage" flat>Back</v-btn>
          </div>

          <div v-else>
            <v-btn @click.native="b1" v-on:click="setpage" flat>Back</v-btn>
          </div>

        </v-stepper-content>
      </v-stepper-items>
    </v-stepper>

    <v-btn @click="generateJson" color="success">generate json</v-btn>
  </v-flex>
</template>

<script>
import { mapGetters } from "vuex";
import PageBootstrap from "@/components/subpages/PageBootstrap.vue";
import PageSchema from "@/components/subpages/PageSchema.vue";
import PageHttp from "@/components/subpages/PageHttp.vue";
import PageCrud from "@/components/subpages/PageCrud.vue";
import PageRest from "@/components/subpages/PageRest.vue";
import PageVuetify from "@/components/subpages/PageVuetify.vue";
import PageEntities from "@/components/subpages/PageEntities.vue";

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
  data() {
    return {
      isActive: false,
      e1: 0,
      b1: 0,
      steps: {
        1: "bootstrap",
        2: "http",
        3: "schema",
        4: "crud",
        5: "rest",
        6: "vuetify",
        7: "entities"
      }
    };
  },
  methods: {
    ...mapGetters(["getgocipe"]),
    setback() {
      this.b1 = this.e1 - 1;
    },
    setpage() {
      this.e1 = this.e1 - 1;
    },
    setlastpage() {
      this.e1 = 7;
      this.b1 = this.e1 - 1;
    },
    loadpage(key) {
      this.e1 = key;
      this.b1 = this.e1 - 1;
    },
    generateJson() {
      var some = this.$store.getters.getgocipe;
      console.log(some);
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
};
</script>

<style>
.link-component {
  color: white;
}

/* .link-component.active {
  text-decoration: underline;
} */

.stepper__step.stepper__step--active {
  border-bottom: 1px solid #004a8f;
}
</style>



