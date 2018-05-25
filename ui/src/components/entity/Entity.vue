<template>
  <div>
    <v-toolbar color="transparent" dark tabs>
      <v-tabs slot="extension" v-model="tab" color="transparent" grow>
        <v-tabs-slider color="blue"></v-tabs-slider>

        <v-tab v-for="(item, index) in items" :key="index">
          {{item}}
        </v-tab>
      </v-tabs>
    </v-toolbar>
    <v-tabs-items v-model="tab">
      <v-tab-item v-for="item in items" :key="item">
        <v-card flat>

          <div v-if="item == 'relationship'">
            <component :is="item" v-model="fields"></component>

          </div>

          <div v-else-if="item == 'crud'">
            <component :is="item" v-model="crud"></component>

          </div>

          <div v-else-if="item == 'schema'">
            <component :is="item" v-model="schema"></component>

          </div>

          <div v-else-if="item == 'widget'">
            <component :is="item" v-model="widget"></component>

          </div>

          <div v-else>
            <component :is="item" v-model="entity"></component>
          </div>

        </v-card>
      </v-tab-item>
    </v-tabs-items>

    <!-- <div> -->
    <!-- Name <input type="text" v-model="val.name" @change="update"> Description
        <input type="text" v-model="val.description" @change="update"> -->
    <!-- </div> -->

    <!-- <v-btn color="success" v-on:click=" appendentities()">Complete</v-btn> -->

  </div>
</template>

<script>
import general from "@/components/entity/General.vue";
import relationship from "@/components/entity/Relationship.vue";
import widget from "@/components/entity/Widget.vue";
import crud from "@/components/entity/Crud.vue";
import rest from "@/components/entity/Rest.vue";
import schema from "@/components/entity/Schema.vue";
export default {
  components: {
    general,
    relationship,
    widget,
    crud,
    rest,
    schema
  },
  props: ["value"],
  data() {
    return {
      tab: null,
      entity: null,
      crud: null,
      schema: null,
      widget: null,
      fields: [],
      items: ["general", "relationship", "schema", "crud", "rest", "widget"],
      val: {
        name: "",
        description: ""
      }
    };
  },
  mounted() {
    if (this.value !== null) {
      this.val = this.value;
    }

    this.entity.crud = this.crud;
    this.$emit("input", this.val);

    this.entity.schema = this.schema;
    this.$emit("input", this.val);

    this.entity.fields = this.fields;
    this.$emit("input", this.val);

    this.entity.widget = this.widget;
    this.$emit("input", this.val);

    this.$emit("input", this.entity);
  },
  methods: {
    // appendtoEntity() {
    //   this.entity.fields = this.fields;
    //   this.$emit("input", this.val);
    // },
    // appendcrud() {
    //   this.entity.crud = this.crud;
    //   this.$emit("input", this.val);
    // },
    // appendschema() {
    //   this.entity.schema = this.schema;
    //   this.$emit("input", this.val);
    // },
    // appendwidget() {
    //   this.entity.widget = this.widget;
    //   this.$emit("input", this.val);
    // },
    // appendentities() {
    //   this.$emit("input", this.entity);
    // },
    // next() {
    //   this.tab++;
    // }
  }
};
</script>

<style>
.layout.wrapper.row.wrap {
  margin-left: 40px;
}
</style>
