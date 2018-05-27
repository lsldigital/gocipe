import Vue from "vue";
import Router from "vue-router";
import Progress from "./components/Progress.vue";
import Upload from "./components/subpages/Upload.vue";

Vue.use(Router);

export default new Router({
  routes: [
    {
      path: "/",
      name: "home",
      component: Progress
    },
    {
      path: "/upload",
      name: "upload",
      component: Upload
    }
  ]
});
