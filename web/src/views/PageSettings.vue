<template>
  <div class="page">
    <v-container>
      <h1 class="headline">Settings</h1>
      <template v-for="setting in settings">

        <component :is="setting.type" v-model="setting.value" :key="setting.attr.label" v-bind="setting.attr"></component>

        <!-- <div :key="setting.attr.label" v-if="setting.type == 'v-select'">
            <v-select v-bind="setting.attr" :item-value="setting.value" v-model="setting.value"></v-select>
          </div> -->
      </template>
    </v-container>
  </div>
</template>

<script>
import Counter from "./components/Counter";
import { mapGetters } from "vuex";

export default {
  components: {
    Counter
  },
  computed: {
    ...mapGetters(["settings"])
  }
};
</script>
