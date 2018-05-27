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

    <v-flex xs10>
      <v-layout class=" mx-auto mb-5 mt-3" column wrap>
        <div column class="text-xs-center">
          <h2>Upload your gocipe file</h2>
          <input type="file" @change="onFileChange">

        </div>
      </v-layout>
      <v-layout class=" mx-auto mb-5 mt-3" column wrap>
        <div class="text-xs-center">
          <v-btn color="success" @click="generateJson">generate json</v-btn>

          <router-link :to="{ name: 'home' }">Apply</router-link>

          <!-- <li><router-link :to="{ name: 'home' }">home</router-link></li> -->
        </div>
      </v-layout>
    </v-flex>
  </v-flex>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
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
    ...mapActions([
      "addrest",
      "addcrud",
      "addschema",
      "addvuetify",
      "addbootstrap",
      "addhttp",
      "addentities"
    ]),

    onFileChange: function(event) {
      var input = event.target;
      if (input.files && input.files[0]) {
        var reader = new FileReader();

        reader.onload = (e => {
          var filedata = e.target.result;
          var obj = JSON.parse(filedata);
          if (obj.bootstrap !== undefined) {
            this.addbootstrap(obj.bootstrap);
          }

          if (obj.http !== undefined) {
            this.addhttp(obj.http);
          }

          if (obj.schema !== undefined) {
            this.addschema(obj.schema);
          }

          if (obj.crud !== undefined) {
            this.addcrud(obj.crud);
          }

          if (obj.rest !== undefined) {
            this.addrest(obj.rest);
          }

          if (obj.vuetify !== undefined) {
            this.addvuetify(obj.vuetify);
          }
          var ent = [];

          console.log();
          if (obj.entities !== undefined) {
            console.log(obj.entities);
            this.addentities(obj.entities);
          }
        }).bind(this);
        // Start the reader job - read file as a data url (base64 format)
        reader.readAsText(input.files[0]);

        var recipe = this.$store.state["gocipe"];

        const data = JSON.stringify(recipe);

        this.$;
      }
    },
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
    generateJson: function(gocipe) {
      var recipe = this.$store.state["gocipe"];
      const data = JSON.stringify(recipe);
      const blob = new Blob([data], { type: "text/plain" });
      const e = document.createEvent("MouseEvents"),
        a = document.createElement("a");
      a.download = "gocipe.json";
      a.href = window.URL.createObjectURL(blob);
      a.dataset.downloadurl = ["text/json", a.download, a.href].join(":");
      e.initEvent(
        "click",
        true,
        false,
        window,
        0,
        0,
        0,
        0,
        0,
        false,
        false,
        false,
        false,
        0,
        null
      );
      a.dispatchEvent(e);
    }
  }
};
</script>

<style>
.link-component {
  color: white;
}

.stepper__step.stepper__step--active {
  border-bottom: 1px solid #004a8f;
}
</style>



