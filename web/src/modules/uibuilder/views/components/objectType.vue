<template>
    <v-layout>
        <v-flex md2 class="root-branch">
            <h3>{{property || "root"}} </h3>
        </v-flex>
        <v-flex md10 class="children-branch">
            <template v-if="isArray(value)">
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
            </template>
            <component v-else v-for="(objectValue, objectProperty, objectIndex) in value" :key="objectIndex + '__' + objectProperty" :is="getType(objectValue)" :object="value" :property="objectProperty" :index="objectIndex" :value="objectValue"> </component>
        </v-flex>
    </v-layout>

</template>

<script>
import draggable from "vuedraggable";
import booleanType from "@/modules/uibuilder/views/components/booleanType.vue";
import stringType from "@/modules/uibuilder/views/components/stringType.vue";
import Dialog from "../dialog/Dialog.vue";

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

.root-branch {
  /* border-top: 2px dashed black; */
  /* border-right: 2px dashed black; */
}
.root-branch h3 {
  padding-left: 10px;
  text-transform: capitalize;
}

.children-branch {
  margin-bottom: 40px;
  /* border-left: 2px dashed black; */
  /* background: #dcdcdc; */
}
</style>
