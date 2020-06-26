<template>
  <v-app>
    <v-main>
      <v-container class="fill-height">
        <v-row align="center" justify="center">
          <v-col cols="4">
            <v-card>
              <v-card-title>Login</v-card-title>
              <v-card-text class="d-flex flex-column">
                <v-btn
                  class="my-2 text-none"
                  tile
                  :href="login.url"
                  v-for="(login, index) in logins"
                  :key="index"
                >
                  <v-icon>{{login.icon}}</v-icon>
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
      name: 'Github',
      icon: 'mdi-github',
      url: `${this.$store.state.base}/login/github`
    },
    {
      name: 'Gitea',
      icon: '$gitea',
      url: `${this.$store.state.base}/login/gitea`
    }
  ];

  mounted() {
    console.log(this.$store.state.user.current);
  }
}
</script>
