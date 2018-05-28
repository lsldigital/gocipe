<template>
    <div>
        <v-layout class="borderwrapper" row wrap>
            <v-layout class="wrapper" row wrap>

                <v-flex xs12>
                    <v-subheader class="entityheader"> Schema Options</v-subheader>
                </v-flex>

                <v-flex xs2 class="mx-auto" v-if="value !== null">
                    <v-checkbox v-model="value.create" label="Create" hint="Create whether or not to generate CREATE TABLE"></v-checkbox>
                </v-flex>
                <v-flex xs2 class="mx-auto" v-else>
                    <v-checkbox v-model="schema.create" label="Create" hint="Create whether or not to generate CREATE TABLE"></v-checkbox>
                </v-flex>

                <v-flex xs2 class="mx-auto" v-if="value !== null">
                    <v-checkbox v-model="value.drop" label="Drop" hint="Drop whether or not to generate DROP IF EXISTS before CREATE"></v-checkbox>
                </v-flex>
                <v-flex xs2 class="mx-auto" v-else>
                    <v-checkbox v-model="schema.drop" label="Drop" hint="Drop whether or not to generate DROP IF EXISTS before CREATE"></v-checkbox>
                </v-flex>

                <v-flex xs2 class="mx-auto" v-if="value !== null">
                    <v-checkbox v-model="value.aggregate" label="Aggregate" hint="Aggregate whether or not to generate schema into single file instead of separate files"></v-checkbox>
                </v-flex>

                <v-flex xs2 class="mx-auto" v-else>
                    <v-checkbox v-model="schema.aggregate" label="Aggregate" hint="Aggregate whether or not to generate schema into single file instead of separate files"></v-checkbox>
                </v-flex>

                <v-flex xs10 class="mx-auto" v-if="value !== null">
                    <v-text-field v-model="value.path" label="Path" hint="Path indicates in which path to generate the schema sql file"></v-text-field>
                </v-flex>
                <v-flex xs10 class="mx-auto" v-else>
                    <v-text-field v-model="schema.path" label="Path" hint="Path indicates in which path to generate the schema sql file"></v-text-field>
                </v-flex>

            </v-layout>

            <!-- <v-flex xs12>
                <v-layout column wrap>
                    <v-divider></v-divider>
                    <v-flex class=" mx-auto">
                        <v-btn color="info" v-on:click="addschema()" dark large>Next</v-btn>
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
      schema: {
        create: true,
        drop: true,
        aggregate: true,
        path: ""
      }
    };
  },
  props: ["value"],
  mounted() {
    if (typeof this.value === undefined) {
      this.schema.create = true;
      this.schema.drop = true;
      this.schema.aggregate = true;
      this.schema.path = "";
    }

    this.$emit("input", this.schema);
  },
  methods: {
    addschema: function() {
      //   this.$emit("input", this.schema);
    }
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
