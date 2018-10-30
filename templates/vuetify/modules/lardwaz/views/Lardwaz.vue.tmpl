<template>
  <div class="page">
    <div class="dashboard-wrapper">
      <!--<v-btn fab dark right fixed color="teal" :absolute="returnTrue()" class="save-button" @click="saveArticle()">-->
        <!--<v-icon dark>save</v-icon>-->
      <!--</v-btn>-->
      <!--<v-btn fab dark right fixed color="blue-grey" :absolute="returnTrue()" class="restore-button" @click="$store.dispatch('restore')">-->
        <!--<v-icon dark>settings_backup_restore</v-icon>-->
      <!--</v-btn>-->
      <v-layout row wrap fill-height="true" v-if="ready" class="editor">
        <v-flex xs4 class="editor-left-wrapper">
          <BlockEditor :information="information" />
        </v-flex>
        <v-flex xs8 class="editor-right-wrapper">
          <Preview transition="fade-transition" :information="information"/>
        </v-flex>
      </v-layout>
      <v-layout  row wrap fill-height v-else>
        <v-container fluid fill-height justify-center align-center>
          <v-progress-circular indeterminate :size="50" color="primary"></v-progress-circular>
        </v-container>
      </v-layout>


      <!--<v-snackbar-->
        <!--:timeout="3000"-->
        <!--color="success"-->
        <!--vertical="vertical"-->
        <!--v-model="saveArticleResult"-->
      <!--&gt;-->
        <!--Article saved.-->
        <!--<v-btn dark flat @click.native="saveArticleResult = false">Close</v-btn>-->
      <!--</v-snackbar>-->
    </div>

  </div>
</template>

<script>
import BlockPicker from "./components/BlockPicker";
import BlockEditor from "./components/BlockEditor";
import Information from "./components/Information";
import Preview from "@lardwaz-config/Preview";

export default {
  data() {
    return {
      ready: false,
      saveArticleResult: false
    };
  },
  props: {
    information: {
      default: null,
      type: Object
    }
  },
  mounted() {
    setTimeout(() => {
      this.ready = true;
    }, 500);
  },
  methods: {
    // returnTrue() {
    //   return true;
    // },
    // saveArticle() {
    //   this.$store.dispatch("lardwaz/save");
    //   this.saveArticleResult = true;
    // }
  },
  components: {
    BlockPicker,
    BlockEditor,
    Preview,
    Information
  }
};
</script>

<style lang="scss">
.dashboard-wrapper {
  overflow: hidden;
  min-height: 100vh;
}
.editor {
  overflow: hidden;
  min-height: 100vh;
  position: relative;
}
.editor-left-wrapper {
  overflow-y: scroll;
  background: #ebebeb;
}
.editor-right-wrapper {
  overflow-y: scroll;
}
.save-button {
  top: 80px;
  right: 120px;
}
.restore-button {
  top: 80px;
  right: 50px;
}
.options-tab {
  height: calc(100vh - 84px);
}
</style>
