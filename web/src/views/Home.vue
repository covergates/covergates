<template>
  <v-container class="mt-5">
    <v-row justify="center" align="center">
      <v-avatar size="56" class="mr-5">
        <img :src="require('@/assets/logo.png')" />
      </v-avatar>
      <span class="text-h2 font-weight-bold">Covergates</span>
    </v-row>
    <v-banner class="text-center primary--text">The portal gates to good code quality.</v-banner>
    <v-row justify="center" class="ma-10">
      <v-col
        cols="4"
        class="d-flex justify-center"
        v-for="(intro, index) in introductions"
        :key="index"
      >
        <v-card flat class="d-flex flex-column justify-center align-center">
          <v-card-title class="text-h4 font-weight-medium">{{intro.title}}</v-card-title>
          <v-avatar size="128" tile>
            <v-icon size="128">{{intro.icon}}</v-icon>
          </v-avatar>
        </v-card>
      </v-col>
    </v-row>
    <v-row class="mt-10">
      <v-spacer />
      <v-btn color="primary" to="/repo" large dark>Get start!</v-btn>
      <v-spacer />
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { Component } from 'vue-property-decorator';
import Vue from '@/vue';

type intro = {
  title: string;
  icon: string;
};

@Component({
  name: 'Home'
})
export default class HomeVue extends Vue {
  introductions = [
    {
      title: 'Open Source',
      icon: 'mdi-open-source-initiative'
    },
    {
      title: 'Self-Hosted',
      icon: 'mdi-server'
    }
  ] as intro[];
}
</script>
