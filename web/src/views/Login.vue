<template>
  <v-app>
    <v-main>
      <v-container class="fill-height">
        <v-row align="center" justify="center">
          <v-col cols="12" md="5">
            <v-card color="accent" elevation="12" dark>
              <div class="d-flex justify-center logo-box">
                <v-avatar color="primary" class="logo" size="56">
                  <img :src="require('@/assets/logo.png')" />
                </v-avatar>
              </div>
              <v-card-title
                primary-title
                class="justify-center text-h4 font-weight-black"
              >Covergates</v-card-title>
              <v-card-text class="d-flex flex-column">
                <v-btn
                  class="my-2 text-none text-h6"
                  tile
                  large
                  :href="login.url"
                  v-for="(login, index) in logins"
                  :key="index"
                >
                  <v-icon large>{{login.icon}}</v-icon>
                  <v-spacer />
                  Login with {{login.name}}
                  <v-spacer />
                </v-btn>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<script lang="ts">
import { Component } from 'vue-property-decorator';
import Vue from '@/vue';

interface LoginSCM {
  name: string;
  icon: string;
  url: string;
}

@Component
export default class Login extends Vue {
  private logins: LoginSCM[] = [
    {
      name: 'Gitea',
      icon: '$gitea',
      url: `${this.$store.state.base}/login/gitea`
    },
    {
      name: 'Github',
      icon: 'mdi-github',
      url: `${this.$store.state.base}/login/github`
    },
    {
      name: 'GitLab',
      icon: 'mdi-gitlab',
      url: `${this.$store.state.base}/login/gitlab`
    },
    {
      name: 'Bitbucket',
      icon: 'mdi-bitbucket',
      url: `${this.$store.state.base}/login/bitbucket`
    }
  ];

  mounted() {
    console.log(this.$store.state.user.current);
  }
}
</script>

<style lang="scss" scoped>
.logo-box {
  max-height: 28px;
  .logo {
    position: 'relative';
    top: -28px;
  }
}
</style>
