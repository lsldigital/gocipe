<template>
    <v-card class="mb-5" color="dark" height="100%">
        <v-toolbar color="#212121">
            <i class="material-icons">settings</i>
            <v-toolbar-title>Choose Desired Crud Operations</v-toolbar-title>
            <v-spacer></v-spacer>
        </v-toolbar>
        <v-layout class="borderwrapper" row wrap>
            <v-layout row wrap class="page_space">
                <v-flex xs10>
                    <v-subheader class="entityheader">Basic Crud operations</v-subheader>
                </v-flex>

                <v-flex xs4>
                    <v-checkbox v-model="crud.create" label="Create"></v-checkbox>
                </v-flex>

                <v-flex xs4>
                    <v-checkbox v-model="crud.hooks.pre_save" label="Pre Save"></v-checkbox>
                </v-flex>

                <v-flex xs4>
                    <v-checkbox v-model="crud.hooks.post_save" label="Post Save"></v-checkbox>
                </v-flex>

                <v-flex xs12>
                    <v-layout column wrap>
                        <v-divider></v-divider>
                    </v-layout>
                </v-flex>

                <v-flex xs4>
                    <v-checkbox v-model="crud.read" label="Read"></v-checkbox>
                </v-flex>

                <v-flex xs4>
                    <v-checkbox v-model="crud.hooks.pre_read" label="Pre Read"></v-checkbox>
                </v-flex>

                <v-flex xs4>
                    <v-checkbox v-model="crud.hooks.post_read" label="Post Read"></v-checkbox>
                </v-flex>

                <v-flex xs12>
                    <v-layout column wrap>
                        <v-divider></v-divider>
                    </v-layout>
                </v-flex>

                <v-flex xs4>
                    <v-checkbox v-model="crud.read_list" label="Read List"></v-checkbox>
                </v-flex>

                <v-flex xs4>
                    <v-checkbox v-model="crud.hooks.pre_list" label="Pre List"></v-checkbox>
                </v-flex>

                <v-flex xs4>
                    <v-checkbox v-model="crud.hooks.post_list" label="Post List"></v-checkbox>
                </v-flex>

                <v-flex xs12>
                    <v-layout column wrap>
                        <v-divider></v-divider>
                    </v-layout>
                </v-flex>

                <v-flex xs4>
                    <v-checkbox v-model="crud.delete" label="Delete"></v-checkbox>
                </v-flex>

                <v-flex xs4>
                    <v-checkbox v-model="crud.hooks.pre_delete" label="Pre Delete"></v-checkbox>
                </v-flex>

                <v-flex xs4>
                    <v-checkbox v-model="crud.hooks.post_delete" label="Post Delete"></v-checkbox>
                </v-flex>

                <v-flex xs12>
                    <v-layout column wrap>
                        <v-divider></v-divider>
                    </v-layout>
                </v-flex>

                <v-flex xs4>
                    <v-checkbox v-model="crud.merge" label="Merge"></v-checkbox>
                </v-flex>

                <v-flex xs12>
                    <v-layout column wrap>
                        <v-divider></v-divider>
                    </v-layout>
                </v-flex>

            </v-layout>
        </v-layout>

        <v-flex xs12 class="text-xs-center">
            <v-btn @click="pushcrud" icon color="primary">
                <i class="material-icons"> check_circle </i>
            </v-btn>
        </v-flex>
    </v-card>
</template>
<script>
import { mapActions } from "vuex";
import gocipe from "../../assets/gocipe.json";
export default {
  data() {
    return {
      gocipe,
      crud: {
        create: true,
        read: true,
        read_list: true,
        update: true,
        delete: true,
        merge: false,
        hooks: {
          pre_save: false,
          post_save: false,
          pre_read: false,
          post_read: false,
          pre_list: false,
          post_list: false,
          pre_delete: false,
          post_delete: false
        }
      }
    };
  },
  mounted() {
    if (this.gocipe.crud !== undefined) {
      console.log(this.gocipe.hooks);

      if (this.gocipe.crud.hooks !== undefined) {
        this.crud.hooks.pre_save =
          this.gocipe.crud.hooks.pre_save !== undefined
            ? this.gocipe.crud.hooks.pre_save
            : false;
        this.crud.hooks.post_save =
          this.gocipe.crud.hooks.post_save !== undefined
            ? this.gocipe.crud.hooks.post_save
            : false;
        this.crud.hooks.pre_read =
          this.gocipe.crud.hooks.pre_read !== undefined
            ? this.gocipe.crud.hooks.pre_read
            : false;
        this.crud.hooks.post_read =
          this.gocipe.crud.hooks.post_read !== undefined
            ? this.gocipe.crud.hooks.post_read
            : false;
        this.crud.hooks.pre_list =
          this.gocipe.crud.hooks.pre_list !== undefined
            ? this.gocipe.crud.hooks.pre_list
            : false;
        this.crud.hooks.post_list =
          this.gocipe.crud.hooks.post_list !== undefined
            ? this.gocipe.crud.hooks.post_list
            : false;
        this.crud.hooks.pre_delete =
          this.gocipe.crud.hooks.pre_delete !== undefined
            ? this.gocipe.crud.hooks.pre_delete
            : false;
        this.crud.hooks.post_delete =
          this.gocipe.crud.hooks.post_delete !== undefined
            ? this.gocipe.crud.hooks.post_delete
            : false;
      }

      this.crud.create =
        this.gocipe.crud.create === undefined ? false : this.gocipe.crud.create;
      this.crud.read =
        this.gocipe.crud.read === undefined ? false : this.gocipe.crud.read;
      this.crud.read_list =
        this.gocipe.crud.read_list === undefined
          ? false
          : this.gocipe.crud.read_list;
      this.crud.update =
        this.gocipe.crud.update === undefined ? false : this.gocipe.crud.update;
      this.crud.delete =
        this.gocipe.crud.delete === undefined ? false : this.gocipe.crud.delete;
      this.crud.merge =
        this.gocipe.crud.merge === undefined ? false : this.gocipe.crud.merge;
    }
  },
  methods: {
    ...mapActions(["addcrud"]),
    pushcrud() {
      this.addcrud(this.crud);
    }
  }
};
</script>

<style>
.layout.page_space.row.wrap {
  padding: 40px;
}
</style>
