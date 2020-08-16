<template>
  <v-card flat>
    <v-card-title class="primary white--text">Default Branch Timeline</v-card-title>
    <v-card-text>
      <v-timeline dense clipped>
        <v-timeline-item
          fill-dot
          color="accent"
          small
          right
          v-for="(report, index) in history"
          :key="index"
        >
          <v-card>
            <v-list-item>
              <v-progress-circular
                size="48"
                :value="coverage(report)"
                class="d-none d-sm-flex"
                color="accent"
              >{{coverage(report)}}</v-progress-circular>
              <v-list-item-content class="ml-5">
                <v-list-item-title>{{shortSHA(report)}}</v-list-item-title>
                <v-list-item-subtitle class="d-flex d-sm-none primary--text">{{coverage(report)}}%</v-list-item-subtitle>
                <v-list-item-subtitle>{{uploadData(report)}}</v-list-item-subtitle>
              </v-list-item-content>
              <v-list-item-action class="d-flex align-center justify-center">
                <v-btn small color="accent" class="d-none d-sm-flex" :to="reportLink(report)">Report</v-btn>
                <v-btn class="d-flex d-sm-none" fab x-small color="accent" :to="reportLink(report)">
                  <v-icon>mdi-file-chart</v-icon>
                </v-btn>
              </v-list-item-action>
            </v-list-item>
          </v-card>
        </v-timeline-item>
      </v-timeline>
    </v-card-text>
  </v-card>
</template>

<script lang="ts">
import { Location } from 'vue-router';
import { Component } from 'vue-property-decorator';
import Vue from '@/vue';

@Component
export default class ReportTimeline extends Vue {
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
    return Math.round((sum / report.coverages.length) * 1000) / 10;
  }

  shortSHA(report: Report): string {
    return report.commit.substr(0, 16);
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
