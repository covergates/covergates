<template>
  <v-card class="ma-5">
    <v-toolbar flat>
      <v-toolbar-title class="grey--text">{{filePath}}</v-toolbar-title>
    </v-toolbar>
    <v-divider></v-divider>
    <v-card-text>
      <table cellspacing="0" cellpadding="0">
        <tbody>
          <tr v-for="(line, index) in codeLines" :key="index">
            <td class="line-number">{{index+1}}</td>
            <td>
              <pre v-html="line"></pre>
            </td>
          </tr>
        </tbody>
      </table>
    </v-card-text>
  </v-card>
</template>

<script lang="ts">
import { Component } from 'vue-property-decorator';
import Vue from '@/vue';

@Component
export default class ReportSource extends Vue {
  get filePath(): string {
    return this.$route.params.path;
  }

  get sourceCode() {
    const source = this.$store.state.report.source;
    return source ? this.$highlight(source) : '';
  }

  get codeLines(): string[] {
    return this.sourceCode.split(/\r?\n/);
  }
}
</script>

<style lang="scss" scoped>
@import '@/assets/styles/variables';

table {
  border: none;
  .line-number {
    user-select: none;
    width: 55px;
    color: $line-number-color;
    font-size: 12px;
  }
}
</style>
