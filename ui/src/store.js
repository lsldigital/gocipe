import Vue from "vue";
import Vuex from "vuex";

Vue.use(Vuex);

const state = {
  gocipe: {
    bootstrap: {},
    http: {},
    schema: {}
  }
};

const mutations = {
  addbootstrap(state, payload) {
    // state.gocipe.bootstrap.bootstrap = payload
    Vue.set(state.gocipe, "bootstrap", payload)
  },
  addschema(state, payload) {
    Vue.set(state.gocipe, "schema", payload)
  }
};

const actions = {
  addbootstrap({ commit }, payload) {
    commit("addbootstrap", payload)
  },
  addschema({commit}, payload) {
    commit("addschema", payload)
  }
};

const store = new Vuex.Store({
  state:state,
  mutations:mutations,
  actions:actions
});

export default store;
