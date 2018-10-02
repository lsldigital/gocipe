<template>
    <!-- <v-layout> -->
    <!-- <v-flex md12> -->
    <v-checkbox :label="property" v-model="object[property]"></v-checkbox>
    <!-- </v-flex> -->
    <!-- </v-layout> -->
</template>

<script>
export default {
  props: ["value", "property", "index", "object"]
};
</script>
