<template>
  <div>
    <v-tooltip top v-if="one">
      <timeago slot="activator" locale="fr-FR" :datetime="date"></timeago>
      <span>{{ date.toLocaleDateString() }} {{ date.toLocaleTimeString() }}</span>
    </v-tooltip>
    <span v-else>{{ date.toLocaleDateString() }} {{ date.toLocaleTimeString() }}</span>
  </div>
</template>

<script>
export default {
  inheritAttrs: false,
  data() {
    return {
      date: new Date()
    };
  },
  created() {
    this.date = this.$attrs.time;
  },
  computed: {
    one() {
      const oneDay = 2 * 24 * 60 * 60 * 1000;
      const today = new Date().getTime();
      const check = new Date(this.date).getTime();

      const difference = Math.abs(check - today);

      if (difference < oneDay) {
        return true;
      }
      return false;
    }
  }
};
</script>
