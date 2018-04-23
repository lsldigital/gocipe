<template>
  <v-layout mb-2>
    <v-flex md2 class="root-branch" v-if="property && !isArray(value)">
      <header @click.stop="collapse = !collapse">
        <span>{{property || "root"}} { }
          <transition name="slide-fade">
            <v-icon v-if="collapse">expand_more</v-icon>
            <v-icon v-else>expand_less</v-icon>
          </transition>
        </span>
      </header>
    </v-flex>
    <transition name="slide-fade">
      <v-flex v-show="!collapse" class="children-branch" :class="{'md10': property, 'md12': !property, 'md12': isArray(value) }">
        <template v-if="isArray(value)">
          <header @click.stop="arrcollapse = !arrcollapse">
            <span class="array">{{property}} [ ]</span>
          </header>

          <template v-if="!arrcollapse">
            <v-expansion-panel>
              <draggable class="full-width" v-model="computedValue">
                <v-expansion-panel-content v-for="(objectValue,  objectIndex) in computedValue" :key="objectIndex + '__' + property">
                  <h4 slot="header">{{Object.keys(computedValue[objectIndex])[0]}} {{computedValue[objectIndex][Object.keys(computedValue[objectIndex])[0]]}}</h4>
                  <v-card>
                    <v-card-text>
                      <component :is="getType(objectValue)" :object="value" :index="objectIndex" :value="objectValue"> </component>
                      <v-btn outline small color="red" @click="del(value,objectIndex)">Delete</v-btn>
                    </v-card-text>
                  </v-card>

                </v-expansion-panel-content>
              </draggable>
            </v-expansion-panel>
            <add-dialog :property="property" @interface="updateArray" />
          </template>
        </template>

        <component v-else v-for="(objectValue, objectProperty, objectIndex) in value" :key="objectIndex + '__' + objectProperty" :is="getType(objectValue)" :object="value" :property="objectProperty" :index="objectIndex" :value="objectValue"> </component>
      </v-flex>
    </transition>

  </v-layout>
</template>

<script>
import draggable from "vuedraggable";
import booleanType from "@/modules/uibuilder/views/components/booleanType.vue";
import stringType from "@/modules/uibuilder/views/components/stringType.vue";
import AddDialog from "@/modules/uibuilder/views/dialog/AddDialog.vue";

export default {
  data() {
    return {
      collapse: false,
      arrcollapse: false
    };
  },
  props: ["value", "property", "index", "object"],
  computed: {
    computedValue: {
      get() {
        return this.value;
      },
      set(result) {
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
      this.computedValue.push($event);
    },
    del(value, index) {
      value.splice(index, 1);
    }
  },
  components: {
    draggable,
    booleanType,
    stringType,
    AddDialog
  }
};
</script>

<style lang="scss" scoped>
.full-width {
  width: 100%;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s;
}
.fade-enter,
.fade-leave-to {
  opacity: 0;
}

.slide-fade-enter-active {
  transition: all 0.3s ease;
}
.slide-fade-leave-active {
  transition: all 0.2s cubic-bezier(1, 0.5, 0.8, 1);
}
.slide-fade-enter {
  transform: translateY(-10px);
  opacity: 0;
}
.slide-fade-leave-to {
  transform: translateY(-30px);
  opacity: 0;
}

.root-branch {
  /* border-top: 2px dashed black; */
  /* border-right: 2px dashed black; */
}
.children-branch {
  margin-bottom: 10px;
}
.children-branch header {
  -moz-user-select: none;
  -webkit-user-select: none;
  user-select: none;
  cursor: se-resize;
}
.children-branch header span {
  display: inline-block;
  text-transform: capitalize;
  // background: white;
  color: black;
  padding: 10px 30px 10px 10px;
  box-shadow: 0 2px 1px -1px rgba(0, 0, 0, 0.2), 0 1px 1px 0 rgba(0, 0, 0, 0.14),
    0 1px 3px 0 hsla(0, 0%, 0%, 0.122);
  border-radius: 10px 50% 50% 10px;
  // border-bottom: 2px solid #ddd;

  /* margin-left: -20px; */
}

.children-branch header span.array {
  background: white;
  color: black;
  padding: 10px 30px 10px 10px;
  box-shadow: 0 2px 1px -1px rgba(0, 0, 0, 0.2), 0 1px 1px 0 rgba(0, 0, 0, 0.14),
    0 1px 3px 0 rgba(0, 0, 0, 0.12);
  border-radius: 10px 50px 0% 0px;
  border-radius: 0;
  border: none;
}

.children-branch {
  // margin-bottom: 40px;
  /* border-left: 2px dashed black; */
  /* background: #dcdcdc; */
}
</style>
