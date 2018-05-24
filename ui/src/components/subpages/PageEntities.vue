<template>
  <v-card color="dark" class="mb-5" height="100%">

    <v-divider></v-divider>
    <v-flex xs12>
      <v-layout color="black" column align-center class="pa-3">
        <div class="text-xs-center">
          <div>
            <Entity v-model="entities[index]" v-for="(entity, index) in entities" :key="index" />

            <div class="entity_btn_zone ">
              <v-btn class="addcustomcolor" @click="addEntity()" dark large>Add New Entity</v-btn>
              <v-btn color="primary" dark large>Remove</v-btn>
            </div>
          </div>
        </div>
      </v-layout>
    </v-flex>
  </v-card>
</template>

<script>
import Entity from "@/components/entity/Entity.vue";
export default {
  components: {
    Entity
  },

  data() {
    return {
      entities: []
    };
  },
  methods: {
    addEntity() {
      this.entities.push({});
    }
  }
};
</script>

<style>
.subheader.entityheader {
  font-weight: 100;
  font-size: 18px;
}

.breaker {
  width: 100%;
  padding: 10px;
}

.entity_btn_zone {
  background-color: #00000026;
  height: 75px;
  padding-top: 8px;
}
</style>
