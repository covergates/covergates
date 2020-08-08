<template>
  <v-container>
    <v-timeline dense clipped class="ml-5">
      <v-timeline-item
        fill-dot
        color="accent"
        small
        right
        v-for="(report, index) in history"
        :key="index"
      >
        <v-row class="pt-1">
          <v-col cols="6" md="2" xl="1" align-self="center" class="text-h4">{{coverage(report)}}%</v-col>
          <v-col cols="6" md="3" xl="1" class="text-md-center text-xl-left" align-self="center">
            <v-btn small color="primary" :to="reportLink(report)">{{shortSHA(report)}}</v-btn>
          </v-col>
          <v-col cols="12" md="2" align-self="center">
            <strong>{{uploadData(report)}}</strong>
          </v-col>
        </v-row>
      </v-timeline-item>
    </v-timeline>
  </v-container>
</template>

<script lang="ts">
import { Location } from 'vue-router';
import { Component } from 'vue-property-decorator';
import Vue from '@/vue';

@Component
export default class ReportHistory extends Vue {
  get history(): Report[] {
    return this.$store.state.report.history;
  }

  coverage(report: Report): number {
    if (report.coverages.length <= 0) {
      return 0;
    }
    let sum = 0;
    for (const coverage of report.coverages) {
      sum += coverage.statementCoverage;
    }
    return Math.round((sum / report.coverages.length) * 10000) / 100;
  }

  shortSHA(report: Report): string {
    return report.commit.substr(0, 10);
  }

  uploadData(report: Report): string {
    if (report.createdAt) {
      const date = new Date(report.createdAt);
      return date.toLocaleDateString();
    } else {
      return '';
    }
  }

  coverageColor(report: Report): string {
    const cov = this.coverage(report);
    if (cov > 80) {
      return 'green';
    } else if (cov > 40) {
      return 'amber';
    } else {
      return 'red';
    }
  }

  reportLink(report: Report): Location {
    return {
      name: 'report-overview',
      query: {
        ref: report.commit
      }
    };
  }
}
</script>

<style lang="scss">
.v-timeline-item__opposite {
  max-width: 100px !important;
}
</style>
