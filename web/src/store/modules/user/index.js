import actions from "./actions";
import getters from "./getters";
import mutations from "./mutations";

const state = {
  auth: true,
  name: {
    first: "John",
    last: "Doe"
  },
  role: "admin"
};

const namespaced = true;

export default {
  namespaced,
  state,
  actions,
  getters,
  mutations
};
