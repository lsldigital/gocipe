<template>
  <v-app id="app">

    <v-content>
      <v-slide-x-transition mode="out-in">
        <router-view></router-view>
      </v-slide-x-transition>
    </v-content>

    <FooterUI />

  </v-app>
</template>

<script>
import FooterUI from "./views/components/ui-elements/FooterUI";

export default {
  data: () => ({}),
  props: {
    source: String
  },
  components: {
    FooterUI
  }
};
</script>

<style lang="scss">
* {
  //   margin:0;padding:0;
}

body {
  //   margin: 0;
  //   font-size: 2rem;
  //   font-family: 'Roboto', sans-serif,
  //               -apple-system, BlinkMacSystemFont,
  //                'avenir next', avenir,
  //                helvetica, 'helvetica neue',
  //                Ubuntu,
  //    'segoe ui', arial,
  //                sans-serif;
  //  flex:1 1 auto;
}

.page {
  //   text-align: center;
  height: 100%;
  overflow: hidden;
  // -webkit-box-flex: 1;
  // display: flex;
  // flex:1 1 auto;
  /* nesting for the need to test postcss */
  code {
    background-color: #f0f0f0;
    padding: 3px 5px;
    border-radius: 2px;
  }
}

.btns {
  padding: 10px;
  text-transform: uppercase;
  font-size: 14px;
  cursor: pointer;
  background: #fff;
  border: 0;
  box-shadow: 0 0 5px #999;
  border-radius: 5px;
  &:hover {
    box-shadow: 0 0 10px #999;
  }
  &.info {
    background: lightskyblue;
  }
  &.success {
    background: lightgreen;
  }
}

.tooling-header {
  font-size: 14px;
  text-transform: uppercase;
  text-align: left;
  padding: 10px 20px;
}

.top-toolbar {
  .btn-wrapper {
    position: absolute;
    right: 0;
    padding: 0;
  }
}

.footer {
  text-align: center;
  // align-items: center;
  justify-content: flex-end;
  padding-right: 20px;
}

.placeholder-text {
  color: #eee !important;
}

.single-component-wrapper {
  --borderColor: rgb(226, 226, 226);
  --fontSize: 14px;
  --fieldHeight: 40px;
  padding: 10px 20px;
  //   border-bottom: 1px solid var(--borderColor);
  // label {
  //     font-size: 11px;
  //     display: block;
  //     text-transform: uppercase;
  //     font-weight: 600;
  //     padding-bottom: 5px;
  //     color: #555;
  // }
  // input, select, textarea {
  //     display: block;
  //     height: var(--fieldHeight);
  //     width: 100%;
  //     border: 1px solid var(--borderColor);
  //     background: white;
  //     padding: 0 10px;
  //     box-sizing: border-box;
  //     font-size: var(--fontSize);
  //     margin-bottom: 10px;
  //     color: #555;
  // }
  // textarea {
  //     display: block;
  //     height: auto;
  //     min-height: var(--fieldHeight);
  //     // width: 100%;
  //     border: 1px solid  var(--borderColor);
  //     padding: 0 10px;
  //     font-size: var(--fontSize);
  //     resize: vertical;
  // }
}
</style>
