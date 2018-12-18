<template>
  <div class="preview-container">
    <div class="preview-wrapper" :class="viewport">
      <!-- <link rel="stylesheet" :href="'/themes/' + cssfile + '/style.css'"> -->
      <div class="article-wrapper" id="article-wrapper">
        <h1 class="preview__title" v-if="information && information.title">{{ information.title }}</h1>

        <!-- Lardwaz Render -->
        <Renderer class="article__body" :blocks="blocks"/>
      </div>
    </div>

    <!--<v-btn fab dark small :bottom="true" :right="true" color="primary" :absolute="true" class="responsive responsive-mobile" @click="changeViewport('mobile')">-->
    <!--<v-icon dark>phone_iphone</v-icon>-->
    <!--</v-btn>-->
    <!--<v-btn fab dark small :bottom="true" :right="true" color="primary" :absolute="true" class="responsive responsive-tablet" @click="changeViewport('tablet')">-->
    <!--<v-icon dark>tablet_mac</v-icon>-->
    <!--</v-btn>-->
    <!--<v-btn fab dark small :bottom="true" :right="true" color="primary" :absolute="true" class="responsive responsive-desktop" @click="changeViewport('desktop')">-->
    <!--<v-icon dark>desktop_mac</v-icon>-->
    <!--</v-btn>-->
  </div>
</template>

<script>
import { mapGetters } from "vuex";
import Renderer from "@lardwaz-config/Renderer";

import variables from "@lardwaz-renderer/_variables.scss";
import article from "@lardwaz-renderer/renderer.scss";
import normalize from "normalize-css";

export default {
  props: {
    information: {
      default: null,
      type: Object
    }
  },
  data() {
    return {
      currentViewport: "desktop"
    };
  },
  computed: {
    ...mapGetters({
      blocks: "lardwaz/getBlocks",
      viewport: "lardwaz/getViewport"
    }),
    cssfile() {
      const theme = this.settings.theme.value;
      return theme === "" ? "gocipe.mu" : theme;
    }
  },
  components: {
    Renderer,
    variables,
    normalize
  },
  methods: {
    changeViewport(value) {
      this.currentViewport = value;
    }
  },
  beforeMount() {}
};
</script>

<style lang="scss">
.block-wrapper {
  display: block;
  width: 100%;
}

.responsive {
  bottom: 100px !important;
  left: 150px;
  position: fixed;
}

.responsive-tablet {
  left: 100px;
}

.responsive-desktop {
  left: 50px;
}

.preview-container {
  position: relative;
  height: 100%;
}

.preview-wrapper {
  // height: 100vh;
  overflow-y: scroll;
  height: 100%;
  transition: all 0.3s ease-in-out;
  position: relative;

  &.mobile {
    max-width: 400px;
    margin: 0 auto;
    border-width: 25px;
    border-top-width: 0;
    border-bottom-width: 0;
    border-style: solid;
    border-image: linear-gradient(
        to left,
        black 1%,
        white 3%,
        rgba(0, 0, 0, 0.8) 5%,
        black 95%,
        white 97%,
        black 100%
      )
      1 100%;
  }

  &.tablet {
    max-width: 768px;
    margin: 0 auto;
    border-width: 50px;
    border-top-width: 0;
    border-bottom-width: 0;
    border-style: solid;
    border-image: linear-gradient(
        to left,
        black 1%,
        white 3%,
        rgba(0, 0, 0, 0.8) 5%,
        black 95%,
        white 97%,
        black 100%
      )
      1 100%;
  }

  &.desktop {
    max-width: 100%;
    margin: 0 auto;
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s;
}

.fade-enter,
.fade-leave-to {
  opacity: 0;
}

.article-wrapper {
  padding: 20px;
  text-align: left;
  max-width: 900px;
  margin: 0 auto;

  h1.preview__title {
    border-bottom: 5px solid black;
    margin-bottom: 20px;
    width: 100%;
  }
}

figure {
  width: 100%;
  margin: 20px 0 !important;
  line-height: 0;

  img {
    width: 100%;
    margin: 0;
  }

  figcaption {
    font-family: Arial;
    font-weight: 400;
    font-size: 15px;
    line-height: 19px;
    color: #a3a3a3;
    padding: 6px;
    text-align: center;
    background: #f3f3f3;

    p {
      margin: 0;
    }
  }
}

.related__articles {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  grid-auto-rows: auto;
  grid-gap: 10px;
  margin: 20px 0;

  article {
    overflow: hidden;
    background: #fff;
    background-position: center center !important;
    border-radius: 5px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
    position: relative;
    text-align: center;
    padding: 10px;

    a {
      text-decoration: none;

      h3 {
        color: #004a8f;
      }
    }
  }
}

.footer__container {
  text-align: center;
  width: 100%;
  display: block;
  position: relative;
  padding-top: 100px;
  padding-bottom: 100px;

  p {
    text-align: center;
    font-size: 25px;
  }
}

.gallery__container {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  grid-auto-rows: auto;
  grid-gap: 10px;

  .gallery__item {
    max-height: 200px;
    min-height: 200px;
    overflow: hidden;
    background: #cccccc;
    background-position: center center !important;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.2);
    position: relative;

    .caption {
      width: 100%;
      position: absolute;
      bottom: 0;
      background: rgba(0, 0, 0, 0.4);
      color: #fff;
      padding: 10px;
      text-align: center;
    }
  }
}
</style>
