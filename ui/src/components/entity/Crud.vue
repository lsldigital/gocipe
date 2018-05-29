<template>
    <div class="">
        <v-layout class="borderwrapper" row wrap>
            <v-layout row wrap>
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
    </div>
</template>

<script>
export default {
  data() {
    return {
      crud: {
        create: true,
        read: true,
        read_list: true,
        update: true,
        delete: true,
        merge: true,
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
  props: ["value"],
  watch: {
    value: function(query) {
      this.crud = query;
    }
  },
  mounted() {
    if (this.value !== undefined) {
      this.crud = this.value;
    }
    this.$emit("input", this.crud);
  },
  methods: {}
};
</script>

<style>
.layout.borderwrapper.row.wrap {
  padding-top: 20px;
}

.flex.marginbottom.xs5 {
  margin-bottom: 25px;
}

button.addcustomcolor.btn.btn {
  color: black;
  background-color: #c6c6c6;
}
</style>
