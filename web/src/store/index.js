import Vue from "vue";
import Vuex from "vuex";

import site from "./modules/site/index";
import user from "./modules/user/index";
import debugpane from "./modules/debugpane/index";
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
  },
  INCREMENT(state) {
    state.count++;
  },
  DECREMENT(state) {
    state.count--;
  }
};

const getters = {
  settings: state => state.settings,
  meta: state => state.meta
};

const actions = {
  incrementAsync({ commit }) {
    setTimeout(() => {
      commit("INCREMENT");
    }, 200);
  },
  setBlock(context, { index, value }) {
    if (index > 0 && index < context.store.data.length) {
      context.commit("BLOCK_CHANGED", { index, value });
    }
  },
  editContent(context, nid) {
    console.log(nid);
    let record = context.state.posts.filter(item => item.nid === nid)[0];
    console.log(record.meta);
    state.meta = Vue.set(state, "meta", record.meta);

    console.log(record);
  },
  save({ state, commit }) {
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
  restore({ state, commit }) {
    let data = JSON.parse(localStorage.getItem("lardwaz"));
    commit("RESTORED", data);
  }
};

export default new Vuex.Store({
  state,
  modules: {
    user,
    debugpane,
    site
  },
  mutations,
  actions,
  getters
});
