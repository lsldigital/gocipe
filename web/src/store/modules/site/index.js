import actions from "./actions";
import getters from "./getters";
import mutations from "./mutations";

const state = {
  drawer: true
};

const namespaced = true;

export default {
  namespaced,
  state,
  actions,
  getters,
  mutations
};
