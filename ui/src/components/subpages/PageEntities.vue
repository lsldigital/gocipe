<template>
  <v-card color="dark" class="mb-5" height="100%">

    <v-divider></v-divider>
    <v-flex xs12>
      <v-layout color="black" column align-center class="pa-3">
        <div class="text-xs-center widthscreen">
          <div>
            <Entity v-model="entities[index]" v-for="(entity, index) in entities" :key="index">

            </Entity>
            <div class="entity_btn_zone ">

              <v-btn class="addcustomcolor" @click="addEntity()" dark large>Add New Entity</v-btn>
              <v-btn color="success" v-on:click="saveEntities()" dark large>Complete</v-btn>
            </div>
          </div>
        </div>
      </v-layout>
    </v-flex>
  </v-card>
</template>

<script>
import Entity from "@/components/entity/Entity.vue";
import { mapActions } from "vuex";
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
    ...mapActions(["addentities"]),
    addEntity() {
      this.entities.push({});
    },
    getentities() {
      this.$emit("input", this.entities);
    },
    removeEntity(index) {
      this.entities.splice(index, 1);
    },
    saveEntities() {
      this.addentities(this.entities);
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

.text-xs-center.widthscreen {
  width: 100%;
}

button.addcustomcolor.btn.btn--large.theme--dark {
  background-color: #ffffff !important;
  color: black;
}

button.addcustomcolor1.btn.btn--large {
  color: white;
  background-color: #020202 !important;
}
</style>
