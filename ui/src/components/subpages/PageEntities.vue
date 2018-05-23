<template>
  <v-card color="dark" class="mb-5" height="100%">
    <v-toolbar color="#212121">
      <i class="material-icons">settings</i>
      <v-toolbar-title>Entities specification</v-toolbar-title>
      <v-spacer></v-spacer>
    </v-toolbar>

    <v-layout class="borderwrapper" row wrap v-for="entity in entities" :key="entity.name">
      <v-divider></v-divider>
      <v-layout row wrap>
        <v-flex xs10>
          <v-subheader class="entityheader">Basic Information</v-subheader>
        </v-flex>
        <v-flex xs3>
          <v-subheader>Name</v-subheader>
        </v-flex>

        <v-flex xs7>
          <v-text-field v-model="entity.name" :name="entity.name" label="Name"></v-text-field>
        </v-flex>

        <v-flex xs3>
          <v-subheader> Primary Key</v-subheader>
        </v-flex>

        <v-flex xs7>
          <v-text-field v-model="entity.primary_key" :name="entity.type" label="Primary Key"></v-text-field>
        </v-flex>

        <v-flex xs3>
          <v-subheader> Table</v-subheader>
        </v-flex>

        <v-flex xs7>
          <v-text-field v-model="entity.table" :name="entity.type" label="Table Name"></v-text-field>
        </v-flex>

        <v-flex xs3>
          <v-subheader> Description</v-subheader>
        </v-flex>

        <v-flex xs7>
          <v-text-field v-model="entity.description" :name="entity.description" label="Description"></v-text-field>
        </v-flex>
      </v-layout>

      <v-flex xs12>
        <v-subheader> Table constraints</v-subheader>
      </v-flex>

      <div class="breaker"></div>

      <v-layout v-for="(constraint, index) in constraints" :key="index">
        <v-flex xs5>
          <v-text-field :label="'Constraint '+ ++index" :id="'const'+index" v-model="constraints[index]" class="ml-5"></v-text-field>
        </v-flex>

        <v-flex xs1 class="text-xs-center ml-4">
          <v-btn icon color="primary">
            <i class="material-icons"> check_circle </i>
          </v-btn>
        </v-flex>

        <v-flex xs1 class="text-xs-center">
          <v-btn @click="remconst(index)" icon color="primary ">
            <i class="material-icons "> remove_circle_outline </i>
          </v-btn>
        </v-flex>

      </v-layout>

      <v-flex xs12>
        <v-layout column align-center class="pa-3 ">
          <div class="text-xs-center ">
            <div>
              <v-btn color="primary " @click="addconstraints" dark large>Add Constraints</v-btn>
            </div>
          </div>
        </v-layout>
      </v-flex>

      <v-layout row wrap justify-center align-center>
        <entity :myProp="entity"></entity>
      </v-layout>

    </v-layout>
    <v-divider></v-divider>
    <v-flex xs12>
      <v-layout column align-center class="pa-3">
        <div class="text-xs-center">
          <div>
            <v-btn color="primary" @click="addentity" dark large>Add Entity</v-btn>
          </div>
        </div>
      </v-layout>
    </v-flex>
  </v-card>
</template>

<script>
import Entity from "@/components/subpages/Entity.vue";
export default {
  components: {
    Entity
  },

  data() {
    return {
      entities: [],
      constraints: [],
      entity: {
        name: "",
        primary_key: "",
        table: "",
        table_constraints: [],
        description: "",
        fields: [
          {
            label: "",
            serialized: "",
            property: {
              name: "",
              type: ""
            },
            schema: {
              field: "",
              type: "",
              nullable: false,
              default: ""
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
          }
        ],
        schema: {
          create: true,
          drop: true,
          aggregate: true,
          path: ""
        },
        schema: {},
        crud: {},
        rest: {},
        vuetify: {}
      }
    };
  },
  methods: {
    toggle(name) {
      if (name == "generate") {
        this.http.generate = !this.http.generate;
      }
    },
    addconstraints() {
      this.constraints.push("");
    },
    remconst(index) {
      console.log(index);
      this.constraints.splice(index - 1, 1);
    },
    addentity() {
      const obj = {
        name: "",
        primary_key: "",
        table: "",
        table_constraints: [],
        description: "",
        fields: [
          {
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
          }
        ],
        schema: {
          create: true,
          drop: true,
          aggregate: true,
          path: ""
        },
        schema: {},
        crud: {},
        rest: {},
        vuetify: {}
      };
      this.entities.push(obj);
    }
  }
};
</script>

<style>
.subheader.entityheader {
  font-weight: 100;
  font-size: 18px;
}

.breaker {
  width: 100%;
  padding: 10px;
  /* border: 2px solid; */
}
</style>
