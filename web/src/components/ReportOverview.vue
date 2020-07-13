<template>
  <v-container class="container">
    <report-empty v-show="!report" />
    <v-container v-show="report">
      <v-banner single-line>
        <span class="text-h4">Files</span>
      </v-banner>
      <v-sheet class="content">
        <ICountUp :endVal="filesCount" class="count-up text-h2" />
      </v-sheet>
      <v-banner single-line>
        <span class="text-h4">Coverage</span>
      </v-banner>
      <v-sheet class="content">
        <v-progress-circular
          :size="100"
          :width="15"
          :rotate="-90"
          :value="coverage"
          color="primary"
        >{{coverage}}</v-progress-circular>
      </v-sheet>
    </v-container>
  </v-container>
</template>

<script lang="ts">
import { Component } from 'vue-property-decorator';
import ICountUp from 'vue-countup-v2';
import Vue from '@/vue';
import ReportEmpty from '@/components/ReportEmpty.vue';

@Component({
  name: 'report-overview',
  components: {
    ICountUp,
    ReportEmpty
  }
})
export default class ReportOverview extends Vue {
  get report(): Report | undefined {
    return this.$store.state.report.current;
  }

  get coverage(): number {
    if (this.report !== undefined && this.report.coverage !== undefined) {
      return Math.round(this.report.coverage.StatementCoverage * 10000) / 100;
    }
    return 0;
  }

  get filesCount(): number {
    return this.report && this.report.files ? this.report.files.length : 0;
  }
}
</script>

<style lang="scss" scoped>
@import '@/assets/styles/variables';

.container {
  overflow-y: auto;
  max-height: 100%;
  .content {
    padding: 20px;
  }
  .count-up {
    color: $content-color;
  }
}
</style>

<docs>

### Examples
```[import](./__examples__/ReportOverview.vue)
```
</docs>
