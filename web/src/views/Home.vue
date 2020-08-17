<template>
  <perfect-scrollbar class="page-container">
    <v-container>
      <v-row justify="center" align="center" class="mt-5">
        <v-avatar size="56" class="mr-5">
          <img :src="require('@/assets/logo.png')" />
        </v-avatar>
        <span class="text-h2 font-weight-bold">Covergates</span>
      </v-row>
      <v-banner class="text-center primary--text">The portal gates to coverage reports.</v-banner>
      <v-row justify="center" class="ma-10">
        <intro-card v-for="(intro, index) in introductions.slice(0,2)" :key="index" :intro="intro" />
      </v-row>
      <v-row justify="center" class="ma-10">
        <intro-card v-for="(intro, index) in introductions.slice(2,4)" :key="index" :intro="intro" />
      </v-row>
      <v-row class="mt-10 mb-5">
        <v-spacer />
        <v-btn color="primary" to="/repo" large dark>Get start!</v-btn>
        <v-spacer />
      </v-row>
    </v-container>
  </perfect-scrollbar>
</template>

<script lang="ts">
import { Component } from 'vue-property-decorator';
import Vue from '@/vue';
import IntroCard from '@/components/IntroCard.vue';

type intro = {
  title: string;
  icon: string;
  body: string;
};

@Component({
  name: 'Home',
  components: {
    IntroCard
  }
})
export default class HomeVue extends Vue {
  introductions = [
    {
      title: 'Open Source',
      icon: 'mdi-open-source-initiative',
      body: `
      Covergates is open source under GPL-3.0 license.
      It's free for self-hosted.
      Visit our <a href="https://github.com/covergates/covergates">GitHub</a>
      if your are interesting in this project!`
    },
    {
      title: 'Easy to Install',
      icon: 'mdi-clock-fast',
      body: `
      <a href="https://github.com/covergates/covergates/releases">Get the pre-built binary</a>
      for your platform. Simple run the binary and you are ready to go!
      `
    },
    {
      title: 'Simple and Clean',
      icon: 'mdi-monitor-clean',
      body: `
      Review coverage reports with the clean interface. No more distractions.
      `
    },
    {
      title: 'Cross-platform',
      icon: 'mdi-hexagon-multiple-outline',
      body: `
      Covergates is built with <a href="https://golang.org/">Go</a>.
      You can run it on Windows, Linux, macOS and
      any other <a href="https://golang.org/">Go</a> supports platform.
      `
    }
  ] as intro[];
}
</script>

<style lang="scss" scoped>
.page-container {
  overflow-y: auto;
  height: calc(100vh - 64px - 60px);
}
</style>
