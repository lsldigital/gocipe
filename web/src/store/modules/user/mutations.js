import types from "./types";

function logout(state) {
  state.auth = false;
}

function login(state) {
  state.auth = true;
}

export default {
  [types.LOGOUT]: logout,
  [types.LOGIN]: login
};
