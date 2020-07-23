<template>
  <v-app class="app">
    <v-app-bar app color="primary" flat dark clipped-left>
      <v-app-bar-nav-icon @click="routeTo('/')">
        <v-icon>mdi-home</v-icon>
      </v-app-bar-nav-icon>
      <v-btn class="ml-5" text @click="routeTo('/repo')">Repositories</v-btn>
      <v-spacer></v-spacer>
      <account-button />
    </v-app-bar>
    <v-main>
      <v-container class="container pa-0" fluid>
        <router-view></router-view>
      </v-container>
    </v-main>
  </v-app>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import AccountButton from '@/components/AccountButton.vue';

@Component({
  components: {
    AccountButton
  }
})
export default class Main extends Vue {
  routeTo(to: string) {
    this.$router.push(to);
  }
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
