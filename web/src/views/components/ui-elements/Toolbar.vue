<template>
  <div>

  </div>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
import Navigation from "./Navigation";

export default {
  data() {
    return {
      drawer: true
    };
  },
  computed: {
    ...mapGetters(["settings", "user"]),
    ...mapGetters("user", ["getName", "getAuth"]),
    siteName() {
      if (this.settings.siteName.value !== "") {
        return this.settings.siteName.value;
      } else {
        return "Default";
      }
    }
  },
  methods: {
    ...mapActions("user", ["login", "logout"])
  },
  components: {
    Navigation
  }
};
</script>

<style lang="scss" scoped>
.home-link {
  text-decoration: none;
  color: white;
}
</style>
