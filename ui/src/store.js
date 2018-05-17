import Vue from "vue";
import Vuex from "vuex";

Vue.use(Vuex);

const state = {
  gocipe: {
    bootstrap: {},
    http: {}
  }
};

const mutations = {
  addbootstrap(state, payload) {
    state.gocipe.bootstrap = payload
  },
  addhttp(state, payload) {
    state.gocipe.http = payload
  },
};

const actions = {
  addbootstrap({ commit }, payload) {
    commit("addbootstrap", payload)
  },
  addhttp({ commit }, payload) {
    commit("addhttp", payload)
  }
};

const store = new Vuex.Store({
  state:state,
  mutations:mutations,
  actions:actions
});

export default store;
