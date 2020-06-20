<template>
  <v-app>
    <v-main>
      <v-container class="fill-height" fluid>
        <v-row align="center" justify="center">
          <v-col cols="4">
            <v-card>
              <v-card-title>Login With</v-card-title>
              <v-card-text class="d-flex flex-column justify-space-around">
                <v-btn class="mt-5" tile large color="teal">
                  <v-icon>mdi-vuetify</v-icon>Hello
                </v-btn>
                <v-btn class="mt-5" tile large color="teal">
                  <v-icon>mdi-vuetify</v-icon>Hello
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
import { Component, Vue } from 'vue-property-decorator';

class LoginSCM {
  name: string;
  icon: string;
  url: string;
  constructor(name: string, icon: string, url: string) {
    this.name = name;
    this.icon = icon;
    this.url = url;
  }
}

@Component
export default class Login extends Vue {
  private logins: LoginSCM[] = [];
  mounted() {
    const ok = 'ok';
    console.log(`hi ${ok}`);
  }
}
</script>>

<style lang="scss" scoped>
</style>>
