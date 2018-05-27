<template>
  <v-card color="dark" class="mb-5" height="100%">
    <v-divider></v-divider>
    <v-flex xs12>
      <v-layout color="black" column align-center class="pa-3">
        <div class="text-xs-center widthscreen">
          <div v-for="(entity, index) in entities" :key="index">
            <Entity v-model="entities[index]"> </Entity>
            <v-btn class=" mx-auto" dark v-on:click="removeEntities(index)">
              <v-icon dark left>remove_circle</v-icon>Delete
            </v-btn>
          </div>
          <div class="entity_btn_zone ">
            <v-btn class="addcustomcolor" @click="addEntity()" dark large>Add New Entity</v-btn>

            <v-btn color="success" v-on:click="saveEntities()" dark large>Complete</v-btn>
          </div>
        </div>
      </v-layout>
    </v-flex>
  </v-card>
</template>

<script>
import gocipe from "../../assets/gocipe.json";
import Entity from "@/components/entity/Entity.vue";
import { mapActions } from "vuex";
export default {
  components: {
    Entity
  },

  data() {
    return {
      gocipe,
      entities: []
    };
  },
  mounted() {
    this.gocipe = this.$store.state["gocipe"];

    if (this.gocipe.entities !== undefined) {
      this.gocipe.entities.forEach(element => {
        this.entities.push(element);
      });
    } else {
      this.entities.push({});
    }
  },
  methods: {
    ...mapActions(["addentities"]),
    addEntity() {
      this.entities.push({});
    },
    getentities() {
      this.$emit("input", this.entities);
    },
    saveEntities() {
      this.addentities(this.entities);
    },
    removeEntities(index) {
      if (index !== -1) {
        this.entities.splice(index, 1);
      }
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
