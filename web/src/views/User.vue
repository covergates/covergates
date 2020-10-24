<template>
  <perfect-scrollbar class="page-container">
    <v-container fill-height class="d-flex justify-center">
      <v-card flat width="600">
        <div class="d-md-flex flex-no-wrap justify-space-around align-center">
          <div class="d-flex flex-column justify-center align-center">
            <v-avatar size="128" color="grey">
              <v-icon size="64" v-if="avatar===''">mdi-account</v-icon>
              <img :src="avatar" :alt="`${name}-avatar`" v-else />
            </v-avatar>
            <v-card-title class="text-h4">{{name}}</v-card-title>
            <v-card-subtitle class="text-subtitles">{{mail}}</v-card-subtitle>
          </div>
          <div class="user-content">
            <v-card-text>
              <v-list>
                <v-banner class="text-h5">SCM Accounts</v-banner>
                <v-list-item v-for="(bind, index) in bindings" :key="index">
                  <v-list-item-icon>
                    <v-icon size="28">{{scmIcon(bind.name)}}</v-icon>
                  </v-list-item-icon>
                  <v-list-item-content class="text-uppercase">{{bind.name}}</v-list-item-content>
                  <v-list-item-action>
                    <v-icon v-if="bind.active">mdi-check</v-icon>
                    <v-btn text v-else color="red" @click="bindAccount(bind.name)">Link</v-btn>
                  </v-list-item-action>
                </v-list-item>
              </v-list>
            </v-card-text>
          </div>
        </div>
      </v-card>
    </v-container>
  </perfect-scrollbar>
</template>

<script lang="ts">
import { Component } from 'vue-property-decorator';
import Vue from '@/vue';

type bind = {
  name: string;
  active: boolean;
};

@Component({})
export default class UserClass extends Vue {
  get user(): User | undefined {
    return this.$store.state.user.current;
  }

  get scm(): Record<string, boolean> | undefined {
    return this.$store.state.user.scm;
  }

  get avatar(): string {
    return this.user && this.user.avatar ? this.user.avatar : '';
  }

  get name(): string {
    return this.user && this.user.login ? this.user.login : '';
  }

  get mail(): string {
    return this.user && this.user.email ? this.user.email : '';
  }

  get bindings(): bind[] {
    const result = [] as bind[];
    if (this.scm) {
      for (const key in this.scm) {
        result.push({
          name: key,
          active: this.scm[key]
        });
      }
    }
    return result;
  }

  scmIcon(scm: string): string {
    switch (scm) {
      case 'github': {
        return 'mdi-github';
      }
      case 'gitea': {
        return '$vuetify.icons.gitea';
      }
      case 'gitlab': {
        return 'mdi-gitlab';
      }
      case 'bitbucket': {
        return 'mdi-bitbucket';
      }
      default: {
        return 'mdi-source-repository';
      }
    }
  }

  bindAccount(scm: string) {
    const base = this.$store.state.base;
    window.location.href = `${base}/login/${scm}?bind`;
  }
}
</script>

<style lang="scss" scoped>
.user-content {
  min-width: 300px;
}
</style>
