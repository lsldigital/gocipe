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
  props: ["value"],
  mounted() {
    if (this.value !== undefined) {
      this.vuetify = this.value;
    }
    this.$emit("input", this.vuetify);
  },
  watch: {
    value: function(query) {
      this.vuetify = query;
    }
  },
  methods: {}
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
