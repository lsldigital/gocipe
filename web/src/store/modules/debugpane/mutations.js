import types from "./types";

function togglePaneVisibility(state, payload) {
  state.paneVisibility = payload;
}

function closePane(state) {
  state.paneVisibility = false;
}

function openPane(state) {
  state.paneVisibility = true;
}

function addEvent(state, event) {
  state.events.push(event);
}

export default {
  [types.TOGGLE_PANE_VISIBILITY]: togglePaneVisibility,
  [types.CLOSE_PANE]: closePane,
  [types.OPEN_PANE]: openPane,
  [types.ADD_EVENT]: addEvent
};
