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
    <v-footer min-height="60" app fixed dark color="secondary">
      &copy; 2020 covergates
      <v-divider vertical class="ma-4" />
      <v-btn text small>license</v-btn>
      <v-spacer />
      <v-btn text href="https://github.com/covergates/covergates">
        <v-icon left large class="mr-5">mdi-github</v-icon>Github
      </v-btn>
    </v-footer>
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
    this.$router.push(to).catch(() => {
      // do nothing
    });
  }
}
</script>

<style>
.app {
  overflow: hidden;
  height: 100vh;
}

.container {
  height: calc(100vh - 64px - 60px);
}
</style>
