import types from "./types";

function toggleDrawer({ commit, state }) {
  let drawer = !state.drawer;
  commit(types.TOGGLE_DRAWER, drawer);
}

export default {
  [types.TOGGLE_DRAWER]: toggleDrawer
};
