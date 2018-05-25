<template>
    <div>
        <v-layout class="borderwrapper" row wrap>
            <v-layout class="wrapper" row wrap>
                <v-flex xs12>
                    <v-subheader class="entityheader">Vuetify Widget Options</v-subheader>
                </v-flex>
                <v-flex xs3>
                    <v-checkbox v-model="vuetify.generate" label="Generate" hint="Generate represents whether or not to generate vuetify assets"></v-checkbox>
                </v-flex>

                <v-flex xs6>
                    <v-text-field v-model="vuetify.module" id="module" label="Module" hint="Module represents the location where the gocipe module will be generated"></v-text-field>
                </v-flex>

            </v-layout>
            <!-- <v-flex xs12>
                <v-layout column wrap>
                    <v-divider></v-divider>
                    <v-flex class=" mx-auto">
                        <v-btn color="info" v-on:click="addvuetify()" dark large>Add Widget</v-btn>
                    </v-flex>
                    <v-divider></v-divider>
                </v-layout>
            </v-flex> -->

        </v-layout>
    </div>
</template>

<script>
export default {
  data() {
    return {
      vuetify: {
        generate: true,
        module: ""
      }
    };
  },

  mounted() {
    if (typeof this.value === "undefined") {
      this.vuetify.generate = true;
      this.vuetify.module = "";
    }
    this.$emit("input", this.vuetify);
  },
  methods: {
    // addvuetify: function() {
    //   this.$emit("input", this.vuetify);
    // }
  }
};
</script>

<style>
.layout.borderwrapper.row.wrap {
  padding-top: 20px;
}

.flex.marginbottom.xs5 {
  margin-bottom: 25px;
}

button.addcustomcolor.btn.btn {
  color: black;
  background-color: #c6c6c6;
}
</style>
