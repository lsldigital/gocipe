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

      <v-card flat>
        <v-tab-item>
          <general v-model="entity"></general>
        </v-tab-item>
        <v-tab-item>
          <relationship v-model="entity.fields"></relationship>
        </v-tab-item>
        <v-tab-item>
          <schema v-model="entity.schema"></schema>
        </v-tab-item>

        <v-tab-item>
          <crud v-model="entity.crud"></crud>
        </v-tab-item>
        <v-tab-item>
          <rest v-model="entity.rest"></rest>
        </v-tab-item>
        <v-tab-item>
          <widget v-model="entity.widget"></widget>
        </v-tab-item>
      </v-card>

    </v-tabs-items>

    <!-- <div> -->
    <!-- Name <input type="text " v-model="val.name " @change="update "> Description
        <input type="text " v-model="val.description " @change="update "> -->
    <!-- </div> -->

    <!-- <v-btn color="success " v-on:click=" appendentities() ">Complete</v-btn> -->

  </div>
</template>

<script>
import gocipe from "../../assets/gocipe.json";
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
      gocipe,
      tab: null,
      entity: {
        crud: null,
        schema: null,
        widget: null,
        fields: []
      },
      items: ["general", "relationship", "schema", "crud", "rest", "widget"]
    };
  },
  mounted() {
    if (typeof this.value !== undefined) {
      this.entity = this.value;
    }
    this.$emit("input", this.entity);
  },
  methods: {},
  watch: {
    // entity: {
    //   handler(val) {
    //     console.log(val);
    //   },
    //   deep: true
    // }
  }
};
</script>

<style>
.layout.wrapper.row.wrap {
  margin-left: 40px;
}
</style>
