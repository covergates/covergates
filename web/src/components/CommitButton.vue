<template>
  <v-btn small v-clipboard:copy="report?report.commit:''" v-clipboard:success="onCopied">
    <v-icon left>{{buttonIcon}}</v-icon>
    {{commit}}
  </v-btn>
</template>

<script lang="ts">
import { Component } from 'vue-property-decorator';
import Vue from '@/vue';

@Component({
  name: 'commit-button'
})
export default class CommitButton extends Vue {
  private copied = false;

  get buttonIcon(): string {
    return this.copied ? 'mdi-check' : 'mdi-clipboard-outline';
  }

  get report(): Report | undefined {
    return this.$store.state.report.current;
  }

  get commit(): string {
    if (this.report) {
      return this.report.commit.substring(0, 10);
    } else {
      return '';
    }
  }

  onCopied() {
    this.copied = true;
    setTimeout(() => {
      this.copied = false;
    }, 2000);
  }
}
</script>
