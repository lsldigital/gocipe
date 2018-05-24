<template>
  <div class="someclass">
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
          <!-- <general v-model="entity" /> -->
          <!-- <relationship v-model="entity" /> -->
          <component :is="item" v-model="entity"></component>
        </v-card>
      </v-tab-item>
    </v-tabs-items>

    <!-- <div> -->
    <!-- Name <input type="text" v-model="val.name" @change="update"> Description
        <input type="text" v-model="val.description" @change="update"> -->
    <!-- </div> -->
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
  },
  methods: {
    update() {
      this.$emit("input", this.val);
    }
  }
};
</script>
