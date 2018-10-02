<template>
    <v-layout>
        <v-flex md3>
            <v-chip small>{{typeof(value)}}</v-chip>
        </v-flex>
        <v-flex md4>{{property}}</v-flex>
        <v-flex md4><v-text-field :label="`${property}::${typeof(value)}`" v-model="object[property]"></v-text-field></v-flex>
    </v-layout>
</template>

<script>
export default {
  props: ["value", "property", "index", "object"]
};
</script>
