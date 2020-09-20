<template>
  <v-container class="container">
    <report-empty v-show="!report" />
    <v-container v-show="report">
      <v-banner single-line>
        <span class="text-h4 text-capitalize">Breakdown</span>
      </v-banner>
      <v-row justify="space-around">
        <v-card flat>
          <v-card-title>{{branchTitle}}</v-card-title>
          <v-card-text class="text-capitalize text-h6">{{branch}}</v-card-text>
        </v-card>
        <v-card flat>
          <v-card-title>Files</v-card-title>
          <v-card-text>
            <ICountUp :endVal="filesCount" class="count-up text-h2" />
          </v-card-text>
        </v-card>
      </v-row>
      <v-banner single-line>
        <span class="text-h4">Coverage</span>
      </v-banner>
      <v-row justify="space-around">
        <v-card flat>
          <v-card-title>Latest Upload</v-card-title>
          <v-card-text>
            <commit-button class="mr-5" />
            {{uploadDate}}
          </v-card-text>
        </v-card>
        <v-sheet class="content">
          <v-progress-circular
            :size="100"
            :width="15"
            :rotate="-90"
            :value="$coverage"
            color="primary"
          >{{$coverage}}</v-progress-circular>
        </v-sheet>
      </v-row>
    </v-container>
  </v-container>
</template>

<script lang="ts">
import { Component, Mixins } from 'vue-property-decorator';
import ICountUp from 'vue-countup-v2';
import Vue from '@/vue';
import ReportEmpty from '@/components/ReportEmpty.vue';
import CommitButton from '@/components/CommitButton.vue';
import ReportMixin from '@/mixins/report';

@Component({
  name: 'report-overview',
  components: {
    ICountUp,
    ReportEmpty,
    CommitButton
  }
})
export default class ReportOverview extends ((Mixins(
  ReportMixin
) as typeof Vue) && ReportMixin) {
  get repo(): Repository | undefined {
    return this.$store.state.repository.current;
  }

  get report(): Report | undefined {
    return this.$store.state.report.current;
  }

  get branchTitle(): string {
    console.log(this.repo);
    console.log(this.report);
    if (
      this.repo &&
      this.report &&
      this.report.reference &&
      this.report.reference !== ''
    ) {
      return this.report.reference !== this.repo.Branch
        ? 'Current Branch'
        : 'Default Branch';
    }
    return 'Default Branch';
  }

  get branch(): string {
    if (this.report && this.report.reference && this.report.reference !== '') {
      return this.report.reference;
    }
    return this.repo ? this.repo.Branch : 'Master';
  }

  get filesCount(): number {
    return this.report && this.report.files ? this.report.files.length : 0;
  }

  get uploadDate(): string {
    if (this.report && this.report.createdAt) {
      const date = new Date(this.report.createdAt);
      return date.toLocaleDateString();
    } else {
      return 'No Report uploaded';
    }
  }
}
</script>

<style lang="scss" scoped>
@import '@/assets/styles/variables';

.container {
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
