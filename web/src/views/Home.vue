<template>
  <div class="page">
    <!-- <counter></counter> -->
    <v-container fluid fill-height>
        <v-layout align-center justify-center>
          <v-fade-transition>
            <v-flex md3 text-xs-center>
              <div class="logo">
                <!-- <img src="@/assets/logo-src/logo.svg" alt=""> -->
              </div>
            </v-flex>
          </v-fade-transition>
      </v-layout>
    </v-container>

    
    
  </div>
</template>

<script>
// import Counter from './components/Counter'

export default {
  components: {
    // Counter
  }
};
</script>


<style lang="scss" scoped>
@keyframes spin {
  0% {
    // transform: rotateY(0deg);
    transform: scale(0.2) translateY(20px) rotate3d(-2000deg);
  }
  100% {
    // transform: rotateY(360deg);
    transform: scale(1) translateY(0) rotate3d(0deg);
  }
}
.logo {
  animation: spin 2s ease-out 1 alternate;
}
</style>
