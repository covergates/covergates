<template>
  <v-app class="app">
    <v-app-bar app color="primary" dark clipped-left>
      <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>
      <v-spacer></v-spacer>
    </v-app-bar>
    <v-navigation-drawer v-model="drawer" clipped app></v-navigation-drawer>

    <v-main>
      <v-container class="container pa-0" fluid>
        <router-view></router-view>
      </v-container>
    </v-main>
  </v-app>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';

@Component
export default class Main extends Vue {
  private drawer = false;
}
</script>

<style>
.app {
  overflow: hidden;
  height: 100vh;
}

.container {
  height: calc(100vh - 64px);
}
</style>
