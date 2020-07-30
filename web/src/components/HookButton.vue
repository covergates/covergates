<template>
  <v-btn :loading="loading" color="accent" @click="activate">
    <v-icon class="mr-1">{{icon}}</v-icon>activate
  </v-btn>
</template>

<script lang="ts">
import { Component } from 'vue-property-decorator';
import Vue from '@/vue';

@Component({
  name: 'hook-button'
})
export default class HookButton extends Vue {
  private loading = false;
  private icon = '';

  get repo(): Repository | undefined {
    return this.$store.state.repository.current;
  }

  activate() {
    const base = this.$store.state.base;
    if (this.repo) {
      this.loading = true;
      const { NameSpace: space, Name: name, SCM: scm } = this.repo;
      this.$http
        .post(`${base}/api/v1/repos/${scm}/${space}/${name}/hook/create`)
        .then(() => {
          this.icon = 'mdi-check';
        })
        .catch(() => {
          this.icon = 'mdi-close';
        })
        .finally(() => {
          this.loading = false;
          setTimeout(() => {
            this.icon = '';
          }, 5000);
        });
    }
  }
}
</script>
