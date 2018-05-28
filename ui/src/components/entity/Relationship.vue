<template>
  <div>
    <v-layout row wrap class="borderwrapper" v-for="(field, index) in fields" :key="index">

      <fieldcomponent v-model="value[index]" />

      <v-btn class=" mx-auto" dark v-on:click="remove(index)">
        <v-icon dark left>remove_circle</v-icon>Delete Field
      </v-btn>

      <v-flex xs12>
        <v-layout column wrap>
          <v-divider></v-divider>
        </v-layout>
      </v-flex>

    </v-layout>

    <v-btn color="success" v-on:click="newObj()">Add Fields</v-btn>
    <!-- <v-btn color="primary" v-on:click="updatefield()">Append To entity</v-btn> -->

  </div>

</template>

<script>
import fieldcomponent from "@/components/entity/Field.vue";
export default {
  components: {
    fieldcomponent
  },
  props: ["value"],
  data() {
    return {
      fields: []
    };
  },
  mounted() {
    if (this.value.length === 0) {
      this.fields.push({});
    }
    this.$emit("input", this.fields);
  },
  methods: {
    newObj() {
      this.fields.push({});
    },

    // updatefield() {
    //   // this.$emit("input", this.fields);
    // },

    remove(index) {
      this.fields.splice(index, 1);
    }
  }
};
</script>

<style>
</style>
