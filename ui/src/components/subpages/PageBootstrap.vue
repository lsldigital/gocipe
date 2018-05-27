<template>
	<v-card color="dark" class="mb-5" height="100%">
		<v-list-tile avatar @click="toggle('generate')" class="pa-2">
			<v-list-tile-action>
				<v-checkbox v-model="bootstrap.generate"></v-checkbox>
			</v-list-tile-action>
			<v-list-tile-content>
				<v-list-tile-title>Generate</v-list-tile-title>
				<v-list-tile-sub-title>Generate bootstrap</v-list-tile-sub-title>
			</v-list-tile-content>
		</v-list-tile>

		<v-divider></v-divider>

		<v-layout class="borderwrapper" row wrap v-for="(row, index) in bootstrap.settings" v-bind:key="row.id" :v-model="nextid">

			<v-flex xs3>
				<v-subheader>Name</v-subheader>
			</v-flex>

			<v-flex xs7>
				<v-text-field v-model="row.name" :id="'name'+index" :name="row.name" label="Name"></v-text-field>
			</v-flex>

			<v-flex xs3>
				<v-subheader> Type</v-subheader>
			</v-flex>

			<v-flex xs7>
				<v-text-field v-model="row.type" :id="'type'+index" :name="row.type" label="Type"></v-text-field>
			</v-flex>

			<v-flex xs3>
				<v-subheader> Description</v-subheader>
			</v-flex>

			<v-flex xs7>
				<v-text-field v-model="row.description" :id="'desc'+index" :name="row.description" label="Description"></v-text-field>
			</v-flex>

			<v-flex xs3>
				<v-subheader> Env </v-subheader>
			</v-flex>
			<v-flex xs7>
				<v-list-tile avatar @click="toggle('row.env')">
					<v-list-tile-action>
						<v-checkbox v-model="row.env"></v-checkbox>
					</v-list-tile-action>
					<v-list-tile-content>
						<v-list-tile-title>Env Variable</v-list-tile-title>
					</v-list-tile-content>
				</v-list-tile>
			</v-flex>

			<v-flex xs2 class="pa-1">
				<v-btn dark @click="remEnv(index)">
					<v-icon dark left>remove_circle</v-icon>Remove
				</v-btn>
			</v-flex>
		</v-layout>

		<div style="{ color: white, border-bottom:10px }"></div>

		<v-flex xs12>
			<v-layout column align-center class="pa-3">
				<div class="text-xs-center">
					<div>
						<v-btn color="primary" @click="addEnv" dark large>Add Environment Fields</v-btn>
					</div>
				</div>
			</v-layout>
		</v-flex>

		<v-divider></v-divider>
		<v-divider></v-divider>
		<v-divider></v-divider>

		<v-flex xs12 class="text-xs-center">
			<v-btn @click="addboots" icon color="primary">
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
      bootstrap: {
        generate: false,
        settings: []
      },
      rows: [],
      nextid: 0
    };
  },
  mounted() {
    this.gocipe = this.$store.state["gocipe"];
    if (this.gocipe.bootstrap !== undefined) {
      this.bootstrap.generate = this.gocipe.bootstrap.generate;

      if (this.gocipe.bootstrap.settings !== undefined) {
        this.gocipe.bootstrap.settings.forEach(element => {
          this.bootstrap.settings.push(element);
        });
      }
    }
  },
  methods: {
    ...mapActions(["addbootstrap"]),

    toggle(name) {
      if (name == "generate") {
        this.bootstrap.generate = !this.bootstrap.generate;
      }
    },

    addboots() {
      this.addbootstrap(this.bootstrap);
    },

    addEnv() {
      this.nextid++;
      const obj = {
        name: "",
        type: "",
        description: "",
        env: false
      };
      this.bootstrap.settings.push(obj);
    },

    remEnv(index) {
      this.bootstrap.settings.splice(index, 1);
    }
  }
};
</script>

<style>
.borderwrapper:first-child {
  border-top: none;
}

.borderwrapper {
  border-top: 1px solid black;
  padding-bottom: 20px;
}
</style>
