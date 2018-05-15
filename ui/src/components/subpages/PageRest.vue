<template>
  <v-card color="dark" class="mb-5" height="100%">
	  <v-toolbar color="#212121">
      <i class="material-icons">settings</i>
      <v-toolbar-title>Choose Desired REST Operations</v-toolbar-title>
      <v-spacer></v-spacer>
    </v-toolbar>

    <v-list two-line subheader>
      <v-list-tile avatar @click="toggle('rest_create')" class="pa-2">
        <v-list-tile-action>
          <v-checkbox v-model="checkedValues['rest_create']" @click.prevent=""></v-checkbox>
        </v-list-tile-action>
        <v-list-tile-content>
          <v-list-tile-title>Create</v-list-tile-title>
          <v-list-tile-sub-title>Allow Create function for REST operations</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>

			<v-list-tile avatar @click="toggle('rest_pre_create')" class="pa-2">
        <v-list-tile-action>
          <v-checkbox v-model="checkedValues['rest_pre_create']" @click.prevent=""></v-checkbox>
        </v-list-tile-action>
        <v-list-tile-content>
          <v-list-tile-title>Pre Create</v-list-tile-title>
          <v-list-tile-sub-title>Allow Pre Create function for REST operations</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>

			      <v-list-tile avatar @click="toggle('rest_post_create')" class="pa-2">
        <v-list-tile-action>
          <v-checkbox v-model="checkedValues['rest_post_create']" @click.prevent=""></v-checkbox>
        </v-list-tile-action>
        <v-list-tile-content>
          <v-list-tile-title>Post Create</v-list-tile-title>
          <v-list-tile-sub-title>Allow Post Create function for REST operations</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>

      <v-divider></v-divider> 
      <v-list-tile avatar @click="toggle('rest_read')" class="pa-2">
        <v-list-tile-action>
          <v-checkbox v-model="checkedValues['rest_read']" @click.prevent=""></v-checkbox>
        </v-list-tile-action>
        <v-list-tile-content>
        <v-list-tile-title>Read</v-list-tile-title>
        <v-list-tile-sub-title>Allow Read function for REST operations</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>

	      <v-list-tile avatar @click="toggle('rest_pre_read')" class="pa-2">
        <v-list-tile-action>
          <v-checkbox v-model="checkedValues['rest_pre_read']" @click.prevent=""></v-checkbox>
        </v-list-tile-action>
        <v-list-tile-content>
        <v-list-tile-title>Pre Read</v-list-tile-title>
        <v-list-tile-sub-title>Allow Pre Read function for REST operations</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>

			<v-list-tile avatar @click="toggle('rest_post_read')" class="pa-2">
        <v-list-tile-action>
          <v-checkbox v-model="checkedValues['rest_post_read']" @click.prevent=""></v-checkbox>
        </v-list-tile-action>
        <v-list-tile-content>
        <v-list-tile-title>Post Read</v-list-tile-title>
        <v-list-tile-sub-title>Allow Post Read function for REST operations</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>		
      <v-divider></v-divider> 
      <v-list-tile avatar @click="toggle('rest_read_list')" class="pa-2">
        <v-list-tile-action>
          <v-checkbox v-model="checkedValues['rest_read_list']" @click.prevent=""></v-checkbox>
        </v-list-tile-action>
        <v-list-tile-content>
        <v-list-tile-title>Read List</v-list-tile-title>
        <v-list-tile-sub-title>Allow  Read List function for REST operations</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>

			<v-list-tile avatar @click="toggle('rest_pre_read_list')" class="pa-2">
        <v-list-tile-action>
          <v-checkbox v-model="checkedValues['rest_pre_read_list']" @click.prevent=""></v-checkbox>
        </v-list-tile-action>
        <v-list-tile-content>
        <v-list-tile-title>Pre Read List</v-list-tile-title>
        <v-list-tile-sub-title>Allow Pre Read List function for REST operations</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>

			<v-list-tile avatar @click="toggle('rest_post_read_list')" class="pa-2">
        <v-list-tile-action>
          <v-checkbox v-model="checkedValues['rest_post_read_list']" @click.prevent=""></v-checkbox>
        </v-list-tile-action>
        <v-list-tile-content>
        <v-list-tile-title>Post Read List</v-list-tile-title>
        <v-list-tile-sub-title>Allow Post Read List function for REST operations</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>

      <v-divider></v-divider> 
      <v-list-tile avatar @click="toggle('rest_update')" class="pa-2">
        <v-list-tile-action>
          <v-checkbox v-model="checkedValues['rest_update']" @click.prevent=""></v-checkbox>
        </v-list-tile-action>
        <v-list-tile-content>
        <v-list-tile-title>Update</v-list-tile-title>
        <v-list-tile-sub-title>Allow Update function for REST operations</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>

			<v-list-tile avatar @click="toggle('rest_update')" class="pa-2">
        <v-list-tile-action>
          <v-checkbox v-model="checkedValues['rest_update']" @click.prevent=""></v-checkbox>
        </v-list-tile-action>
        <v-list-tile-content>
        <v-list-tile-title>Update</v-list-tile-title>
        <v-list-tile-sub-title>Allow Update function for REST operations</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>

			<v-list-tile avatar @click="toggle('rest_pre_update')" class="pa-2">
        <v-list-tile-action>
          <v-checkbox v-model="checkedValues['rest_pre_update']" @click.prevent=""></v-checkbox>
        </v-list-tile-action>
        <v-list-tile-content>
        <v-list-tile-title>Pre Update</v-list-tile-title>
        <v-list-tile-sub-title>Allow Pre Update function for REST operations</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>

		<v-list-tile avatar @click="toggle('rest_post_update')" class="pa-2">
        <v-list-tile-action>
          <v-checkbox v-model="checkedValues['rest_post_update']" @click.prevent=""></v-checkbox>
        </v-list-tile-action>
        <v-list-tile-content>
        <v-list-tile-title>Post Update</v-list-tile-title>
        <v-list-tile-sub-title>Allow Post Update function for REST operations</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>

      <v-divider></v-divider> 

      <v-list-tile avatar @click="toggle('rest_delete')" class="pa-2">
        <v-list-tile-action>
          <v-checkbox v-model="checkedValues['rest_delete']" @click.prevent=""></v-checkbox>
        </v-list-tile-action>
        <v-list-tile-content>
        <v-list-tile-title>Delete</v-list-tile-title>
        <v-list-tile-sub-title>Allow Delete function for REST operations</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>

			<v-list-tile avatar @click="toggle('rest_pre_delete')" class="pa-2">
        <v-list-tile-action>
          <v-checkbox v-model="checkedValues['rest_pre_delete']" @click.prevent=""></v-checkbox>
        </v-list-tile-action>
        <v-list-tile-content>
        <v-list-tile-title>Pre Delete</v-list-tile-title>
        <v-list-tile-sub-title>Allow Pre Delete function for REST operations</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>

			<v-list-tile avatar @click="toggle('rest_post_delete')" class="pa-2">
        <v-list-tile-action>
          <v-checkbox v-model="checkedValues['rest_post_delete']" @click.prevent=""></v-checkbox>
        </v-list-tile-action>
        <v-list-tile-content>
        <v-list-tile-title>Post Delete</v-list-tile-title>
        <v-list-tile-sub-title>Allow Post Delete function for REST operations</v-list-tile-sub-title>
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
          rest_create: false,
          rest_read: false,
          rest_read_list:false,
          rest_update: false,
          rest_delete: false,
					
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