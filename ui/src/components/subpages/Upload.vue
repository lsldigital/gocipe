<template>
    <v-flex xs12 pa-3>

        <v-flex xs10>
            <v-layout class=" mx-auto mb-5 mt-3" column wrap>
                <v-spacer></v-spacer>
                <div column class="text-xs-center">
                    <h2>Upload your gocipe file</h2>
                    <input type="file" @change="onFileChange">

                </div>
            </v-layout>
            <v-layout class=" mx-auto mb-5 mt-3" column wrap>

                <div class="text-xs-center">

                    <router-link class="applyupload" :to="{ name: 'home' }">Apply Changes from File</router-link>

                </div>
            </v-layout>
        </v-flex>
        <v-divider></v-divider>
        <v-flex xs10>
            <v-layout class=" mx-auto mb-5 mt-3" column wrap>
                <div column class="text-xs-center">
                    <h2>Or</h2>
                    <h2>Create a New Gocipe file</h2>

                </div>
            </v-layout>
            <v-layout class=" mx-auto mb-5 mt-3" column wrap>

                <div class="text-xs-center">

                    <router-link class="routeLink" :to="{ name: 'home' }">New Gocipe</router-link>

                </div>
            </v-layout>
        </v-flex>
    </v-flex>
</template>

<script>
import { mapGetters, mapActions } from "vuex";

export default {
  data() {
    return {
      isActive: false,
      e1: 0,
      b1: 0
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

a.routeLink.router-link-active {
  padding: 20px;
  text-decoration: none;
  font-size: 36px;
  color: #ff9800;
}

a.router-link-active.applyupload {
  text-decoration: none;
  font-size: 30px;
  padding: 20px;
}
</style>



