import Vue from "vue";
import Router from "vue-router";
import Progress from "./components/Progress.vue";
import About from "./views/About.vue";

Vue.use(Router);

export default new Router({
  routes: [
    {
      path: "/",
      name: "home",
      component: Progress
    },
    {
      path: "/about",
      name: "about",
      component: About
    }
  ]
});
