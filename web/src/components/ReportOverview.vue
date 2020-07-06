<template>
  <v-container class="container">
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
</template>

<script lang="ts">
import { Component } from 'vue-property-decorator';
import ICountUp from 'vue-countup-v2';
import Vue from '@/vue';

@Component({
  name: 'report-overview',
  components: {
    ICountUp
  }
})
export default class ReportOverview extends Vue {
  get coverage(): number {
    const report = this.$store.state.report.current;
    if (report !== undefined && report.coverage !== undefined) {
      return report.coverage.StatementCoverage * 100;
    }
    return 0;
  }

  get filesCount(): number {
    const repo = this.$store.state.repository.current;
    if (repo && repo.Files) {
      return repo.Files.length;
    }
    return 0;
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
