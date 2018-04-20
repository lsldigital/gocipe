<template>
    <v-layout>
        <v-flex md2>{{property}}</v-flex>
        <v-flex md10><v-checkbox v-model="object[property]"></v-checkbox></v-flex>
    </v-layout>
</template>

<script>
export default {
  props: ["value", "property", "index", "object"]
};
</script>
