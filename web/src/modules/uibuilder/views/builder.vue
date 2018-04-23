<template>
  <v-layout>
    <v-flex sm12>
      <component :is="getType(gocipe)" :value="gocipe"> </component>
    </v-flex>
    <!-- <v-flex sm12 class="code">
      <pre>
          {{ gocipe }}
        </pre>
    </v-flex> -->
  </v-layout>
</template>

<script>
import draggable from "vuedraggable";
import gocipe from "../gocipe.json";
// import gocipe from "../simple.gocipe.json";
// import gocipe from "../types.json";
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


<style lang="scss" scoped>
.code {
  background: red;
}
</style>




