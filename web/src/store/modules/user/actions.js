import types from "./types";

function login({ commit, state }) {
  commit(types.LOGIN);
}

function logout({ commit, state }) {
  commit(types.LOGOUT);
}

export default {
  login,
  logout
};
