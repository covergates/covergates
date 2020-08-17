<template>
  <v-card flat class="d-flex flex-column justify-begin align-center intro-card" width="400">
    <v-avatar size="96" tile>
      <v-icon size="96" color="secondary">{{intro.icon}}</v-icon>
    </v-avatar>
    <v-card-title class="text-h4 font-weight-medium">{{intro.title}}</v-card-title>
    <v-card-text class="text-center font-weight-medium">
      <p v-html="intro.body"></p>
    </v-card-text>
  </v-card>
</template>

<script lang="ts">
import { Component, Prop } from 'vue-property-decorator';
import Vue from '@/vue';

@Component
export default class IntroCard extends Vue {
  @Prop() readonly intro: {} | undefined;
}
</script>
