import Vue from "vue";
import Router from "vue-router";
import Home from "../views/Home.vue";
import PageSettings from "../views/PageSettings.vue";
import PageUsers from "../views/PageUsers.vue";
import PageGocipe from "@/modules/uibuilder/views/PageGocipe.vue";

Vue.use(Router);

export default new Router({
  mode: "hash",
  routes: [
    {
      path: "/",
      name: "gocipe",
      component: PageGocipe
    }
  ]
});
