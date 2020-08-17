<template>
  <v-app class="app">
    <v-app-bar app color="primary" flat dark clipped-left>
      <v-app-bar-nav-icon class="ml-1" @click="routeTo('/')">
        <v-avatar tile rounded="false">
          <img :src="require('@/assets/logo.png')" />
        </v-avatar>
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
      <v-btn text small href="https://github.com/covergates/covergates/blob/master/LICENSE">license</v-btn>
      <v-spacer />
      <v-btn icon href="https://docs.covergates.com/" class="d-flex d-sm-none">
        <v-icon class="mr-1">mdi-book-open-variant</v-icon>
      </v-btn>
      <v-btn icon href="https://github.com/covergates/covergates" class="d-flex d-sm-none">
        <v-icon class="mr-1">mdi-github</v-icon>
      </v-btn>
      <v-btn text href="https://docs.covergates.com/" class="d-none d-sm-flex">
        <v-icon left size="32" class="mr-5">mdi-book-open-variant</v-icon>Docs
      </v-btn>
      <v-btn text href="https://github.com/covergates/covergates" class="d-none d-sm-flex">
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
