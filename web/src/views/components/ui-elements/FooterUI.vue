<template>
    <v-footer app class="blue-grey darken-4" dark>
        <span v-if="env == 'development'">
            <v-chip small color="blue-grey darken-2" slot="activator" text-color="white">
                <v-avatar>
                    <v-icon>memory</v-icon>
                </v-avatar>
                MODE : {{env}}
            </v-chip>
        </span>
        <span v-else>
            <v-chip small color="blue-grey darken-2" text-color="white">
                <v-avatar>
                    <v-icon>done</v-icon>
                </v-avatar>
                Build v0.1
            </v-chip>
        </span>
        <v-spacer></v-spacer>
        <span>Gocipe UI &copy; 2018 </span>
    </v-footer>
</template>

<script>
export default {
  data() {
    return {
      env: process.env.NODE_ENV
    };
  },
  computed: {},
  methods: {},
  components: {}
};
</script>
