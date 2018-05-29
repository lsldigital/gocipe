<template>
    <div>
        <v-layout class="borderwrapper" row wrap>
            <v-divider></v-divider>
            <v-layout row wrap>
                <v-flex xs10>
                    <v-subheader class="entityheader">Basic Information</v-subheader>
                </v-flex>
                <v-flex xs3>
                    <v-subheader>Name</v-subheader>
                </v-flex>

                <v-flex xs7>
                    <v-text-field v-model="generaldata.label" label="Name"></v-text-field>
                </v-flex>

                <v-flex xs3>
                    <v-subheader> Primary Key</v-subheader>
                </v-flex>

                <v-flex xs7>
                    <v-text-field v-model="generaldata.primary_key" label="Primary Key"></v-text-field>
                </v-flex>

                <v-flex xs3>
                    <v-subheader> Table</v-subheader>
                </v-flex>

                <v-flex xs7>
                    <v-text-field v-model="generaldata.table" label="Table Name"></v-text-field>
                </v-flex>

                <v-flex xs3>
                    <v-subheader> Description</v-subheader>
                </v-flex>

                <v-flex xs7>
                    <v-text-field v-model="generaldata.description" label="Description"></v-text-field>
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
                                <v-btn class="addcustomcolor" @click="addconstraints" dark large>Add Table Constraints</v-btn>
                            </div>
                            <div v-for="(constraint, index) in generaldata.table_constraints" :key="index">
                                <v-list>
                                    <v-list-tile>
                                        <v-text-field v-model=" generaldata.table_constraints[index-1]" :label="'Constraint ' + ++index"></v-text-field>
                                        <v-btn icon color="primary" v-on:click="remove(index)" dark>
                                            <i class="material-icons">clear</i>
                                        </v-btn>
                                    </v-list-tile>
                                </v-list>
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
  watch: {
    value: function(query) {
      this.generaldata = query;
    }
  },
  mounted() {
    if (this.value !== undefined) {
      this.generaldata = this.value;
    }
    this.$emit("input", this.generaldata);
  },
  methods: {
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
