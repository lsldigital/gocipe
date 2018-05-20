<template>
    <v-card class="mb-5" color="dark" height="100%">
        <v-toolbar color="#212121">
            <i class="material-icons">settings</i>
            <v-toolbar-title>Choose Desired Crud Operations</v-toolbar-title>
            <v-spacer></v-spacer>
        </v-toolbar>

        <v-list subheader="subheader" two-line="two-line">
            <v-list-tile @click="toggle('create')" avatar="avatar" class="pa-2">
                <v-list-tile-action>
                    <v-checkbox  v-model="crud.create"></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Create</v-list-tile-title>
                    <v-list-tile-sub-title>Allow Create function for CRUD operations</v-list-tile-sub-title>
                </v-list-tile-content>
            </v-list-tile>

            <v-list-tile @click="toggle('pre_create')" avatar="avatar" class="pa-2">
                <v-list-tile-action>
                    <v-checkbox  v-model="crud.hooks.pre_create"></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Pre Create</v-list-tile-title>
                    <v-list-tile-sub-title>Allow Pre Create function for CRUD operations</v-list-tile-sub-title>
                </v-list-tile-content>
            </v-list-tile>

            <v-list-tile @click="toggle('post_create')" avatar="avatar" class="pa-2">
                <v-list-tile-action>
                    <v-checkbox  v-model="crud.hooks.post_create"></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Post Create</v-list-tile-title>
                    <v-list-tile-sub-title>Allow Post Create function for CRUD operations</v-list-tile-sub-title>
                </v-list-tile-content>
            </v-list-tile>

            <v-divider></v-divider>
            <v-list-tile @click="toggle('read')" avatar="avatar" class="pa-2">
                <v-list-tile-action>
                    <v-checkbox  v-model="crud.read"></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Read</v-list-tile-title>
                    <v-list-tile-sub-title>Allow Read function for CRUD operations</v-list-tile-sub-title>
                </v-list-tile-content>
            </v-list-tile>

            <v-list-tile @click="toggle('pre_read')" avatar="avatar" class="pa-2">
                <v-list-tile-action>
                    <v-checkbox  v-model="crud.hooks.pre_read"></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Pre Read</v-list-tile-title>
                    <v-list-tile-sub-title>Allow Pre Read function for CRUD operations</v-list-tile-sub-title>
                </v-list-tile-content>
            </v-list-tile>

            <v-list-tile @click="toggle('post_read')" avatar="avatar" class="pa-2">
                <v-list-tile-action>
                    <v-checkbox  v-model="crud.hooks.post_read"></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Post Read</v-list-tile-title>
                    <v-list-tile-sub-title>Allow Post Read function for CRUD operations</v-list-tile-sub-title>
                </v-list-tile-content>
            </v-list-tile>
            <v-divider></v-divider>
            <v-list-tile @click="toggle('read_list')" avatar="avatar" class="pa-2">
                <v-list-tile-action>
                    <v-checkbox  v-model="crud.read_list"></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Read List</v-list-tile-title>
                    <v-list-tile-sub-title>Allow Read List function for CRUD operations</v-list-tile-sub-title>
                </v-list-tile-content>
            </v-list-tile>

            <v-list-tile @click="toggle('pre_read_list')" avatar="avatar" class="pa-2">
                <v-list-tile-action>
                    <v-checkbox  v-model="crud.hooks.pre_read_list"></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Pre Read List</v-list-tile-title>
                    <v-list-tile-sub-title>Allow Pre Read List function for CRUD operations</v-list-tile-sub-title>
                </v-list-tile-content>
            </v-list-tile>

            <v-list-tile @click="toggle('post_read_list')" avatar="avatar" class="pa-2">
                <v-list-tile-action>
                    <v-checkbox  v-model="crud.hooks.post_read_list"></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Post Read List</v-list-tile-title>
                    <v-list-tile-sub-title>Allow Post Read List function for CRUD operations</v-list-tile-sub-title>
                </v-list-tile-content>
            </v-list-tile>

            <v-divider></v-divider>
            <v-list-tile @click="toggle('update')" avatar="avatar" class="pa-2">
                <v-list-tile-action>
                    <v-checkbox  v-model="crud.update"></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Update</v-list-tile-title>
                    <v-list-tile-sub-title>Allow Update function for CRUD operations</v-list-tile-sub-title>
                </v-list-tile-content>
            </v-list-tile>

            <v-list-tile @click="toggle('pre_update')" avatar="avatar" class="pa-2">
                <v-list-tile-action>
                    <v-checkbox  v-model="crud.hooks.pre_update"></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Pre Update</v-list-tile-title>
                    <v-list-tile-sub-title>Allow Pre Update function for CRUD operations</v-list-tile-sub-title>
                </v-list-tile-content>
            </v-list-tile>

            <v-list-tile @click="toggle('post_update')" avatar="avatar" class="pa-2">
                <v-list-tile-action>
                    <v-checkbox  v-model="crud.hooks.post_update"></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Post Update</v-list-tile-title>
                    <v-list-tile-sub-title>Allow Post Update function for CRUD operations</v-list-tile-sub-title>
                </v-list-tile-content>
            </v-list-tile>

            <v-divider></v-divider>

            <v-list-tile @click="toggle('delete')" avatar="avatar" class="pa-2">
                <v-list-tile-action>
                    <v-checkbox  v-model="crud.delete"></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Delete</v-list-tile-title>
                    <v-list-tile-sub-title>Allow Delete function for CRUD operations</v-list-tile-sub-title>
                </v-list-tile-content>
            </v-list-tile>

            <v-list-tile @click="toggle('pre_delete')" avatar="avatar" class="pa-2">
                <v-list-tile-action>
                    <v-checkbox  v-model="crud.hooks.pre_delete"></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Pre Delete</v-list-tile-title>
                    <v-list-tile-sub-title>Allow Pre Delete function for CRUD operations</v-list-tile-sub-title>
                </v-list-tile-content>
            </v-list-tile>

            <v-list-tile @click="toggle('post_delete')" avatar="avatar" class="pa-2">
                <v-list-tile-action>
                    <v-checkbox  v-model="crud.hooks.post_delete"></v-checkbox>
                </v-list-tile-action>
                <v-list-tile-content>
                    <v-list-tile-title>Post Delete</v-list-tile-title>
                    <v-list-tile-sub-title>Allow Post Delete function for CRUD operations</v-list-tile-sub-title>
                </v-list-tile-content>
            </v-list-tile>
        </v-list>
    </v-card>
</template>
<script>

    export default {
        data() {
            return {
                crud: {
                    "create": true,
                    "read": true,
                    "read_list": true,
                    "update": true,
                    "delete": true,
                    "hooks": {
                        "pre_create": false,
                        "post_create": false,
                        "pre_read": false,
                        "post_read": false,
                        "pre_list": false,
                        "post_list": false,
                        "pre_update": false,
                        "post_update": false,
                        "pre_delete": false,
                        "post_delete": false
                    }
                }
            }
        },
        methods: {
            toggle(name) {
              console.log(name)
                if (name == "create") {
                    this.crud.create = !this.crud.create
                }
                if (name == "read") {
                    this.crud.read = !this.crud.read
                }
                if (name == "read_list") {
                    this.crud.read_list = !this.crud.read_list
                }
                if (name == "update") {
                    this.crud.update = !this.crud.update
                }
                if (name == "delete") {
                    this.crud.delete = !this.crud.delete
                }
                if (name == "pre_create") {
                    this.crud.hooks.pre_create = !this.crud.hooks.pre_create
                }
                if (name == "post_create") {
                    this.crud.hooks.post_create = !this.crud.hooks.post_create
                }
                if (name == "pre_read") {
                    this.crud.pre_read = !this.crud.hooks.pre_read
                }
                if (name == "post_read") {
                    this.crud.hooks.post_read = !this.crud.hooks.post_read
                }
                if (name == "pre_list") {
                    this.crud.pre_list = !this.crud.hooks.pre_list
                }
                if (name == "post_list") {
                    this.crud.hooks.post_list = !this.crud.hooks.post_list
                }
                if (name == "pre_update") {
                    this.crud.pre_update = !this.crud.hooks.pre_update
                }
                if (name == "post_update") {
                    this.crud.hooks.post_update = !this.crud.hooks.post_update
                }
                if (name == "pre_delete") {
                    this.crud.pre_delete = !this.crud.hooks.pre_delete
                }
                if (name == "post_delete") {
                    this.crud.hooks.post_delete = !this.crud.hooks.post_delete
                }
            }
        }
    }
</script>