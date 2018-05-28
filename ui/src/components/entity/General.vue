<template>
    <div>
        {{value}}
        <v-layout class="borderwrapper" row wrap>
            <v-divider></v-divider>
            <v-layout row wrap>
                <v-flex xs10>
                    <v-subheader class="entityheader">Basic Information</v-subheader>
                </v-flex>
                <v-flex xs3>
                    <v-subheader>Name</v-subheader>
                </v-flex>

                <v-flex xs7 v-if="value !== null">
                    <v-text-field v-model="value.label" label="Name"></v-text-field>
                </v-flex>

                <v-flex xs3>
                    <v-subheader> Primary Key</v-subheader>
                </v-flex>

                <v-flex xs7 v-if="value !== null">
                    <v-text-field v-model="value.primary_key" label="Primary Key"></v-text-field>
                </v-flex>

                <v-flex xs3>
                    <v-subheader> Table</v-subheader>
                </v-flex>

                <v-flex xs7 v-if="value !== null">
                    <v-text-field v-model="value.table" label="Table Name"></v-text-field>
                </v-flex>

                <v-flex xs3>
                    <v-subheader> Description</v-subheader>
                </v-flex>

                <v-flex xs7 v-if="value !== null">
                    <v-text-field v-model="value.description" label="Description"></v-text-field>
                </v-flex>

                <v-flex xs12>
                    <v-layout column wrap>
                        <v-divider></v-divider>
                    </v-layout>
                </v-flex>

                <v-flex xs12 class="marginbottom">
                    <v-flex xs6 class=" mx-auto mb-5 mt-3">
                        <v-layout column wrap>
                            <div>
                                <v-btn class="addcustomcolor" @click="addconstraints" large>Add Table Constraints</v-btn>
                            </div>
                            <div v-if="value !== null">
                                <div v-for="(constraint, index) in value.table_constraints" :key="index">
                                    <v-list>
                                        <v-list-tile>
                                            <v-text-field v-model=" value.table_constraints[index-1]" :label="'Constraint ' + ++index"></v-text-field>
                                            <v-btn icon color="primary" v-on:click="remove(index)" dark>
                                                <i class="material-icons">clear</i>
                                            </v-btn>
                                        </v-list-tile>
                                    </v-list>
                                </div>
                            </div>
                        </v-layout>
                    </v-flex>
                </v-flex>
            </v-layout>

        </v-layout>
    </div>
</template>

<script>
export default {
  data() {
    return {
      generaldata: {
        label: "",
        primary_key: "",
        table: "",
        table_constraints: [],
        description: ""
      }
    };
  },
  props: ["value"],
  mounted() {
    if (typeof this.value === undefined) {
      this.generaldata.label = "";
      this.generaldata.primary_key = "";
      this.generaldata.table = "";
      this.generaldata.table_constraints = [];
      this.generaldata.description = "";
    }

    // if (this.value !== null) {
    //   this.generaldata.name = this.value.name;
    //   this.generaldata.primary_key = this.value.primary_key;
    //   this.generaldata.table = this.value.table;
    //   this.generaldata.description = this.value.description;
    // }
    this.$emit("input", this.generaldata);
  },
  methods: {
    // addgeneral: function() {
    //   this.$emit("input", this.generaldata);
    // },
    addconstraints() {
      this.generaldata.table_constraints.push("");
    },
    remove(index) {
      this.generaldata.table_constraints.splice(index - 1, 1);
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
  color: white;
  background-color: #c6c6c6;
}
</style>
