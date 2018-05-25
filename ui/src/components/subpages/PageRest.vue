<template>
    <v-card class="mb-5" color="dark" height="100%">
        <v-toolbar color="#212121">
            <i class="material-icons">settings</i>
            <v-toolbar-title>Choose Desired REST Operations</v-toolbar-title>
            <v-spacer></v-spacer>
        </v-toolbar>

        <v-list subheader="subheader" two-line="two-line">
            <v-list-tile @click="toggle('create')" avatar="avatar" class="pa-2">
                <v-list-tile-action>
                    <v-checkbox v-model="rest.create"></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Create</v-list-tile-title>
                    <v-list-tile-sub-title>Allow Create function for REST operations</v-list-tile-sub-title>
                </v-list-tile-content>
            </v-list-tile>

            <v-list-tile @click="toggle('pre_create')" avatar="avatar" class="pa-2">
                <v-list-tile-action>
                    <v-checkbox v-model="rest.hooks.pre_create"></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Pre Create</v-list-tile-title>
                    <v-list-tile-sub-title>Allow Pre Create function for REST operations</v-list-tile-sub-title>
                </v-list-tile-content>
            </v-list-tile>

            <v-list-tile @click="toggle('post_create')" avatar="avatar" class="pa-2">
                <v-list-tile-action>
                    <v-checkbox v-model="rest.hooks.post_create"></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Post Create</v-list-tile-title>
                    <v-list-tile-sub-title>Allow Post Create function for REST operations</v-list-tile-sub-title>
                </v-list-tile-content>
            </v-list-tile>

            <v-divider></v-divider>
            <v-list-tile @click="toggle('read')" avatar="avatar" class="pa-2">
                <v-list-tile-action>
                    <v-checkbox v-model="rest.read"></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Read</v-list-tile-title>
                    <v-list-tile-sub-title>Allow Read function for REST operations</v-list-tile-sub-title>
                </v-list-tile-content>
            </v-list-tile>

            <v-list-tile @click="toggle('pre_read')" avatar="avatar" class="pa-2">
                <v-list-tile-action>
                    <v-checkbox v-model="rest.hooks.pre_read"></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Pre Read</v-list-tile-title>
                    <v-list-tile-sub-title>Allow Pre Read function for REST operations</v-list-tile-sub-title>
                </v-list-tile-content>
            </v-list-tile>

            <v-list-tile @click="toggle('post_read')" avatar="avatar" class="pa-2">
                <v-list-tile-action>
                    <v-checkbox v-model="rest.hooks.post_read"></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Post Read</v-list-tile-title>
                    <v-list-tile-sub-title>Allow Read function for REST operations</v-list-tile-sub-title>
                </v-list-tile-content>
            </v-list-tile>
            <v-divider></v-divider>
            <v-list-tile @click="toggle('read_list')" avatar="avatar" class="pa-2">
                <v-list-tile-action>
                    <v-checkbox v-model="rest.read_list"></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Read List</v-list-tile-title>
                    <v-list-tile-sub-title>Allow Read List function for REST operations</v-list-tile-sub-title>
                </v-list-tile-content>
            </v-list-tile>

            <v-divider></v-divider>
            <v-list-tile @click="toggle('update')" avatar="avatar" class="pa-2">
                <v-list-tile-action>
                    <v-checkbox v-model="rest.update"></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Update</v-list-tile-title>
                    <v-list-tile-sub-title>Allow Update function for REST operations</v-list-tile-sub-title>
                </v-list-tile-content>
            </v-list-tile>

            <v-list-tile @click="toggle('pre_update')" avatar="avatar" class="pa-2">
                <v-list-tile-action>
                    <v-checkbox v-model="rest.hooks.pre_update"></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Pre Update</v-list-tile-title>
                    <v-list-tile-sub-title>Allow Pre Update function for REST operations</v-list-tile-sub-title>
                </v-list-tile-content>
            </v-list-tile>

            <v-list-tile @click="toggle('post_update')" avatar="avatar" class="pa-2">
                <v-list-tile-action>
                    <v-checkbox v-model="rest.hooks.post_update"></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Post Update</v-list-tile-title>
                    <v-list-tile-sub-title>Allow Post Update function for REST operations</v-list-tile-sub-title>
                </v-list-tile-content>
            </v-list-tile>

            <v-divider></v-divider>

            <v-list-tile @click="toggle('delete')" avatar="avatar" class="pa-2">
                <v-list-tile-action>
                    <v-checkbox v-model="rest.delete"></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Delete</v-list-tile-title>
                    <v-list-tile-sub-title>Allow Delete function for REST operations</v-list-tile-sub-title>
                </v-list-tile-content>
            </v-list-tile>

            <v-list-tile @click="toggle('pre_delete')" avatar="avatar" class="pa-2">
                <v-list-tile-action>
                    <v-checkbox v-model="rest.hooks.pre_delete"></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Pre Delete</v-list-tile-title>
                    <v-list-tile-sub-title>Allow Pre Delete function for REST operations</v-list-tile-sub-title>
                </v-list-tile-content>
            </v-list-tile>

            <v-list-tile @click="toggle('post_delete')" avatar="avatar" class="pa-2">
                <v-list-tile-action>
                    <v-checkbox v-model="rest.hooks.post_delete"></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Post Delete</v-list-tile-title>
                    <v-list-tile-sub-title>Allow Post Delete function for REST operations</v-list-tile-sub-title>
                </v-list-tile-content>
            </v-list-tile>
        </v-list>

        <v-divider></v-divider>
        <v-divider></v-divider>
        <v-divider></v-divider>
        <v-flex xs12 class="text-xs-center">
            <v-btn @click="pushrest" icon color="primary">
                <i class="material-icons"> check_circle </i>
            </v-btn>
        </v-flex>
    </v-card>
</template>
<script>
import gocipe from "../../assets/gocipe.json";
import { mapActions } from "vuex";
export default {
  data() {
    return {
      gocipe,
      rest: {
        create: true,
        read: true,
        read_list: true,
        update: true,
        delete: true,
        hooks: {
          pre_create: false,
          post_create: false,
          pre_read: false,
          post_read: false,
          pre_list: false,
          post_list: false,
          pre_update: false,
          post_update: false,
          pre_delete: false,
          post_delete: false
        }
      }
    };
  },

  mounted() {
    if (this.gocipe.json == true) {
      if (this.gocipe.rest !== undefined) {
        this.rest.create =
          this.gocipe.rest.create === undefined
            ? true
            : this.gocipe.rest.create;

        this.rest.read =
          this.gocipe.rest.read === undefined ? true : this.gocipe.rest.read;

        this.rest.read_list =
          this.gocipe.rest.read_list === undefined
            ? true
            : this.gocipe.rest.read_list;

        this.rest.update =
          this.gocipe.rest.update === undefined
            ? true
            : this.gocipe.rest.update;

        this.rest.delete =
          this.gocipe.rest.delete === undefined
            ? true
            : this.gocipe.rest.delete;

        if (this.gocipe.rest.hooks !== undefined) {
          this.rest.hooks.pre_create =
            this.gocipe.rest.hooks.pre_create === undefined
              ? false
              : this.gocipe.rest.hooks.pre_create;
          this.rest.hooks.post_create =
            this.gocipe.rest.hooks.post_create === undefined
              ? false
              : this.gocipe.rest.hooks.post_create;
          this.rest.hooks.pre_read =
            this.gocipe.rest.hooks.pre_read === undefined
              ? false
              : this.gocipe.rest.hooks.pre_read;
          this.rest.hooks.post_read =
            this.gocipe.rest.hooks.post_read === undefined
              ? false
              : this.gocipe.rest.hooks.post_read;
          this.rest.hooks.pre_list =
            this.gocipe.rest.hooks.pre_list === undefined
              ? false
              : this.gocipe.rest.hooks.pre_list;
          this.rest.hooks.post_list =
            this.gocipe.rest.hooks.post_list === undefined
              ? false
              : this.gocipe.rest.hooks.post_list;
          this.rest.hooks.pre_update =
            this.gocipe.rest.hooks.pre_update === undefined
              ? false
              : this.gocipe.rest.hooks.pre_update;
          this.rest.hooks.post_update =
            this.gocipe.rest.hooks.post_update === undefined
              ? false
              : this.gocipe.rest.hooks.post_update;
          this.rest.hooks.pre_delete =
            this.gocipe.rest.hooks.pre_delete === undefined
              ? false
              : this.gocipe.rest.hooks.pre_delete;
          this.rest.hooks.post_delete =
            this.gocipe.rest.hooks.post_delete === undefined
              ? false
              : this.gocipe.rest.hooks.post_delete;
        }
      }
    }
  },
  methods: {
    ...mapActions(["addrest"]),
    toggle(name) {
      if (name == "create") {
        this.rest.create = !this.rest.create;
      }
      if (name == "read") {
        this.rest.read = !this.rest.read;
      }
      if (name == "read_list") {
        this.rest.read_list = !this.rest.read_list;
      }
      if (name == "update") {
        this.rest.update = !this.rest.update;
      }
      if (name == "delete") {
        this.rest.delete = !this.rest.delete;
      }
      if (name == "pre_create") {
        this.rest.hooks.pre_create = !this.rest.hooks.pre_create;
      }
      if (name == "post_create") {
        this.rest.hooks.post_create = !this.rest.hooks.post_create;
      }
      if (name == "pre_read") {
        this.rest.hooks.pre_read = !this.rest.hooks.pre_read;
      }
      if (name == "post_read") {
        this.rest.hooks.post_read = !this.rest.hooks.post_read;
      }
      if (name == "pre_list") {
        this.rest.hooks.pre_list = !this.rest.hooks.pre_list;
      }
      if (name == "post_list") {
        this.rest.hooks.post_list = !this.rest.hooks.post_list;
      }
      if (name == "pre_update") {
        this.rest.hooks.pre_update = !this.rest.hooks.pre_update;
      }
      if (name == "post_update") {
        this.rest.hooks.post_update = !this.rest.hooks.post_update;
      }
      if (name == "pre_delete") {
        this.rest.hooks.pre_delete = !this.rest.hooks.pre_delete;
      }
      if (name == "post_delete") {
        this.rest.hooks.post_delete = !this.rest.hooks.post_delete;
      }
    },
    pushrest() {
      this.addrest(this.rest);
    }
  }
};
</script>