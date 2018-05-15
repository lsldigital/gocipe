<template>
  <v-card color="dark" class="mb-5" height="100%">
    <v-list two-line subheader>
      <v-list-tile avatar @click="toggle('create')" class="pa-2">
        <v-list-tile-action>
          <v-checkbox v-model="checkedValues['create']" @click.prevent=""></v-checkbox>
        </v-list-tile-action>
        <v-list-tile-content>
          <v-list-tile-title>Create</v-list-tile-title>
          <v-list-tile-sub-title>Allow create function for generating schema</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>
			<v-divider></v-divider> 
      <v-list-tile avatar @click="toggle('drop')" class="pa-2">
        <v-list-tile-action>
          <v-checkbox v-model="checkedValues['drop']" @click.prevent=""></v-checkbox>
        </v-list-tile-action>
        <v-list-tile-content>
        <v-list-tile-title>Drop</v-list-tile-title>
        <v-list-tile-sub-title>Allow Drop function for generating schema</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>
			<v-divider></v-divider> 
      <v-list-tile avatar @click="toggle('aggregate')" class="pa-2">
        <v-list-tile-action>
          <v-checkbox v-model="checkedValues['aggregate']" @click.prevent=""></v-checkbox>
        </v-list-tile-action>
        <v-list-tile-content>
        <v-list-tile-title>Aggregate</v-list-tile-title>
        <v-list-tile-sub-title>Allow Aggregate function for generating schema</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>
    </v-list>
  </v-card>
</template>

<script>

export default {
    data () {
      return {
        checkedValues: {
          create: false,
          drop: false,
          aggregate:false
        }
      }
    },
    methods: {
      toggle (name) {
        this.checkedValues[name] = !this.checkedValues[name]
      }
    }
  }
</script>