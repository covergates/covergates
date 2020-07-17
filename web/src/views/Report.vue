<template>
  <perfect-scrollbar class="page-container">
    <v-container>
      <v-card flat>
        <v-card-title>{{title}}</v-card-title>
        <v-tabs v-show="!loading">
          <v-tab v-for="tab in tabs" :key="tab.key" :to="tab.link">{{tab.key}}</v-tab>
        </v-tabs>
      </v-card>
      <v-progress-linear :active="loading" :indeterminate="loading"></v-progress-linear>
      <div class="router-container" v-show="!loading">
        <router-view></router-view>
      </div>
    </v-container>
  </perfect-scrollbar>
</template>

<script lang="ts">
import { Component } from 'vue-property-decorator';
import { Location } from 'vue-router';
import Vue from '@/vue';
import ReportOverview from '@/components/ReportOverview.vue';

type tabOptions = {
  key: string;
  link: Location;
};

@Component({
  components: { ReportOverview }
})
export default class ReportView extends Vue {
  get loading(): boolean {
    return (
      this.$store.state.repository.loading || this.$store.state.report.loading
    );
  }

  get repo(): Repository | undefined {
    return this.$store.state.repository.current;
  }

  get title(): string {
    if (this.repo) {
      return `${this.repo.NameSpace}/${this.repo.Name}`;
    }
    return 'Report';
  }

  get report(): Report | undefined {
    return this.$store.state.report.current;
  }

  get tabs(): tabOptions[] {
    const options: tabOptions[] = [
      {
        key: 'Overview',
        link: {
          name: 'report-overview'
        }
      }
    ];
    if (this.report) {
      options.push({
        key: 'Code',
        link: {
          name: 'report-code'
        }
      });
    }
    if (this.repo) {
      options.push({
        key: 'Setting',
        link: {
          name: 'report-setting'
        }
      });
    }
    return options;
  }
}
</script>

<style lang="scss" scoped>
.router-container {
  height: calc(100% - 48px);
}
.page-container {
  overflow-y: auto;
  height: 100%;
}
</style>
