<template>
  <div>

    <v-layout row wrap class="borderwrapper" v-for="(field, index) in fields" :key="index">

      <fieldcomponent v-model="fields[index]" />

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
    if (this.value !== undefined) {
      var i;
      try {
        for (i = 0; this.value.length; i++) {
          this.fields.push(this.value[i]);
        }
      } catch (err) {
        console.log("error");
      }
    }
    this.$emit("input", this.fields);
  },
  watch: {
    value: function(query) {
      this.fields = query;
    }
  },
  methods: {
    newObj() {
      this.fields.push({
        label: "",
        serialized: "",
        property: {
          name: "",
          type: ""
        },
        schema: {
          field: "",
          type: ""
        },
        relationship: {
          type: "",
          target: {
            entity: "",
            endpoint: "",
            query: "",
            table: "",
            thisid: "",
            thatid: "",
            thatfield_type: ""
          }
        },
        widget: {
          type: "",
          options: [
            {
              value: "",
              label: ""
            }
          ],
          target: {
            endpoint: "",
            label: ""
          }
        },
        filterable: true
      });
    },
    remove(index) {
      this.fields.splice(index, 1);
    }
  }
};
</script>

<style>
</style>
