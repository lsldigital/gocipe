import Vue from "vue";
import Vuex from "vuex";

Vue.use(Vuex);

const state = {
  gocipe: {
    bootstrap: {},
    http: {},
    schema: {},
    crud: {},
    rest: {},
    vuetify: {},
    entities: []
  }
};

const getters = {
  gocipe(state) {
    return state.gocipe;
  }
};
const mutations = {
  addbootstrap(state, payload) {
    // state.gocipe.bootstrap.bootstrap = payload
    Vue.set(state.gocipe, "bootstrap", payload);
  },
  addschema(state, payload) {
    Vue.set(state.gocipe, "schema", payload);
  },
  addhttp(state, payload) {
    Vue.set(state.gocipe, "http", payload);
  },
  addcrud(state, payload) {
    Vue.set(state.gocipe, "crud", payload);
  },
  addrest(state, payload) {
    Vue.set(state.gocipe, "rest", payload);
  },
  addvuetify(state, payload) {
    Vue.set(state.gocipe, "vuetify", payload);
  },
  addentities(state, payload) {
    Vue.set(state.gocipe, "entities", payload);
  },
  addexisting(state, payload) {
    Vue.set(state, "gocipe", payload);
  }
};

const actions = {
  addbootstrap({ commit }, payload) {
    commit("addbootstrap", payload);
  },
  addschema({ commit }, payload) {
    commit("addschema", payload);
  },
  addhttp({ commit }, payload) {
    commit("addhttp", payload);
  },
  addcrud({ commit }, payload) {
    commit("addcrud", payload);
  },
  addrest({ commit }, payload) {
    commit("addrest", payload);
  },
  addvuetify({ commit }, payload) {
    commit("addvuetify", payload);
  },
  addentities({ commit }, payload) {
    commit("addentities", payload);
  },
  addexisting({ commit }, payload) {
    commit("addexisting", payload);
  }
};

const store = new Vuex.Store({
  state: state,
  getters: getters,
  mutations: mutations,
  actions: actions
});

export default store;
