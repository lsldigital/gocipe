import Vue from "vue";
import Vuex from "vuex";

import site from "./modules/site/index";
import user from "./modules/user/index";
import objectType from "@/modules/uibuilder/views/components/objectType.vue";

import SiteSettings from "./site-settings.js";

Vue.use(Vuex);
Vue.component("objectType", objectType);

const state = {
  page: {},
  settings: SiteSettings
};

const mutations = {
  RESTORED(state, payload) {
    state.lardwaz.blocks = Vue.set(state, "blocks", payload.blocks);
    state.settings = Vue.set(state, "settings", payload.settings);
    state.meta = Vue.set(state, "meta", payload.meta);
    console.log(payload.meta);
  }
};

const getters = {
  settings: state => state.settings,
  meta: state => state.meta
};

const actions = {
  save({ state }) {
    let blocks = state.lardwaz.blocks;
    let settings = state.settings;
    let meta = state.meta;
    let data = JSON.stringify({
      blocks,
      settings,
      meta
    });
    localStorage.setItem("lardwaz", data);
  },
  restore({ commit }) {
    let data = JSON.parse(localStorage.getItem("lardwaz"));
    commit("RESTORED", data);
  }
};

export default new Vuex.Store({
  state,
  modules: {
    user,
    site
  },
  mutations,
  actions,
  getters
});
