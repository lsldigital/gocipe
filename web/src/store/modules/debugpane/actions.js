import types from "./types";

function togglePaneVisibility({ commit, state }, value) {
  let pane = !state.paneVisibility;
  console.log("setting pane to" + pane);
  commit(types.TOGGLE_PANE_VISIBILITY, pane);
}

function closePane({ commit }) {
  commit(types.CLOSE_PANE);
}

function openPane({ commit }) {
  commit(types.OPEN_PANE);
}

function addEvent({ commit }, message) {
  let event = {
    type: "error",
    text: "this is my message"
  };
  commit(types.ADD_EVENT, event);
}
export default {
  [types.TOGGLE_PANE_VISIBILITY]: togglePaneVisibility,
  [types.CLOSE_PANE]: closePane,
  [types.OPEN_PANE]: openPane,
  [types.ADD_EVENT]: addEvent
};
