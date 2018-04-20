function getName(state) {
  return `${state.name.first}  ${state.name.last}`;
}

function getAuth(state) {
  return state.auth;
}

export default {
  getName,
  getAuth
};
