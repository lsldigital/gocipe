<template>
    <v-text-field :label="property" v-model="object[property]"></v-text-field>
</template>

<script>
export default {
  props: ["value", "property", "index", "object"]
};
</script>
