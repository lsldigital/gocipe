import types from "./types";

function toggleDrawer(state, payload) {
  console.log("new state" + payload);
  state.drawer = payload;
}

export default {
  [types.TOGGLE_DRAWER]: toggleDrawer
};
