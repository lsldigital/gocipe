<template>
  <div>
      <v-layout>
          <v-flex sm9>
            <!-- <draggable v-model="myArray">
                <transition-group>
                    <div v-for="element in myArray" :key="element.id">
                        {{element.id}} - {{element.name}}
                        <v-switch :label="element.name" v-model="element.bobo"></v-switch>
                    </div>
                </transition-group>
            </draggable> -->

          
          <div v-for="(value, property, index) in gocipe" :key="property+index">
              <component :is="getType(value)" :object="gocipe" :property="property" :index="index" :value="value"> </component>
          </div>

          </v-flex>
          <v-flex sm3>
              <!-- <pre>
                {{ myArray }}
              </pre> -->
              <pre>
                {{ gocipe }}
              </pre>
          </v-flex>
      </v-layout>
    
  </div>
</template>

<script>
import draggable from "vuedraggable";
// import gocipe from "../gocipe.json";
import gocipe from "../simple.gocipe.json";
import booleanType from "./components/booleanType.vue";
import stringType from "./components/stringType.vue";

export default {
  data() {
    return {
      gocipe
    };
  },
  methods: {
    getType: value => {
      let TYPEOF = typeof value;
      return TYPEOF + "Type";
    }
  },
  components: {
    draggable,
    booleanType,
    stringType
  }
};
</script>
