<template>
    <div>
        <v-layout v-if="isArray(value)">
            <v-flex md2>{{property}} </v-flex>
            <v-flex md10>
                <v-expansion-panel>
                    <draggable class="full-width" v-model="computedValue">
                        <v-expansion-panel-content v-for="(objectValue,  objectIndex) in computedValue" :key="objectIndex + '__' + property">
                            <h4 slot="header">{{Object.keys(computedValue[objectIndex])[0]}} {{computedValue[objectIndex][Object.keys(computedValue[objectIndex])[0]]}}</h4>
                            <v-card>
                                <v-card-text>
                                    <component :is="getType(objectValue)" :object="value" :property="objectIndex" :index="objectIndex" :value="objectValue"> </component>
                                    <v-btn outline small color="red">Delete</v-btn>
                                </v-card-text>
                            </v-card>
                        </v-expansion-panel-content>
                    </draggable>
                </v-expansion-panel>
                <Dialog :property="property" @interface="updateArray" />
            </v-flex>
        </v-layout>
        <v-layout v-else v-for="(objectValue, objectProperty, objectIndex) in value" :key="objectIndex + '__' + objectProperty">
            <v-flex md12>
                <component :is="getType(objectValue)" :object="value" :property="objectProperty" :index="objectIndex" :value="objectValue"> </component>
            </v-flex>
        </v-layout>
    </div>

</template>

<script>
import draggable from "vuedraggable";
import booleanType from "@/modules/uibuilder/views/components/booleanType.vue";
import stringType from "@/modules/uibuilder/views/components/stringType.vue";
import Dialog from "../dialog/Dialog.vue";

import Vue from "vue";

export default {
  data() {
    return {};
  },
  props: ["value", "property", "index", "object"],
  computed: {
    computedValue: {
      get() {
        return this.value;
      },
      set(result) {
        console.log(this.object);
        console.log(this.property);
        console.log(result);
        this.object[this.property] = result;
      }
    }
  },
  methods: {
    getType: value => {
      let TYPE_OF = typeof value;
      //   if (Array.isArray(value)) {
      //     TYPE_OF = "array";
      //   }
      return TYPE_OF + "Type";
    },
    isArray: value => {
      if (Array.isArray(value)) {
        return true;
      }
      return false;
    },
    updateArray($event) {
      const copy = Object.assign($event);
      console.log(copy);
      this.computedValue.push(copy);
    }
  },
  components: {
    draggable,
    booleanType,
    stringType,
    Dialog
  }
};
</script>

<style>
.full-width {
  width: 100%;
}
</style>
