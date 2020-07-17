<template>
  <v-card class="ma-5">
    <v-toolbar flat>
      <v-toolbar-title class="grey--text">{{filePath}}</v-toolbar-title>
    </v-toolbar>
    <v-divider></v-divider>
    <v-card-text>
      <table cellspacing="0" cellpadding="0">
        <tbody>
          <tr v-for="(line, index) in codeLines" :key="index" :class="[hitClass(index+1)]">
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
import { Component, Watch } from 'vue-property-decorator';
import Vue from '@/vue';

@Component
export default class ReportSource extends Vue {
  hitMap = {} as { [key: number]: boolean };

  mounted() {
    this.updateHitMap();
  }

  get report(): Report | undefined {
    return this.$store.state.report.current;
  }

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

  updateHitMap() {
    this.hitMap = {};
    if (this.report && this.report.coverage) {
      const file = this.report.coverage.Files.find(file => {
        return file.Name === this.filePath;
      });
      if (file) {
        for (const hit of file.StatementHits) {
          this.hitMap[hit.LineNumber] = hit.Hits > 0;
        }
      }
    }
  }

  hitClass(i: number): string {
    return this.hitMap[i] ? 'statement-hit' : 'statement-miss';
  }

  @Watch('report')
  onReportChange() {
    this.updateHitMap();
  }
}
</script>

<style lang="scss" scoped>
@import '@/assets/styles/variables';

table {
  border: none;
  table-layout: fixed;
  word-wrap: break-word;
  width: 100%;
  pre {
    white-space: pre-wrap;
    word-wrap: break-word;
  }
  .line-number {
    user-select: none;
    width: 55px;
    color: $line-number-color;
    font-size: 12px;
  }
}

.statement-hit {
  background-color: $hit-statement-color;
}
</style>
