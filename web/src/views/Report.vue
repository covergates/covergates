<template>
  <perfect-scrollbar class="page-container">
    <v-container>
      <v-card flat>
        <v-card-title class="d-flow align-center">
          <v-icon class="mr-5" size="36">{{avatar}}</v-icon>
          <router-link
            :to="{name: 'report-overview', query: {}}"
            class="black--text text-h5 font-weight-bold"
          >{{title}}</router-link>
          <v-btn icon class="ml-5" :href="repoURL">
            <v-icon small>mdi-open-in-new</v-icon>
          </v-btn>
        </v-card-title>
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

  get owner(): boolean {
    return this.$store.state.repository.owner;
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

  get history(): Report[] {
    return this.$store.state.report.history;
  }

  get user(): User | undefined {
    return this.$store.state.user.current;
  }

  get repoURL(): string {
    return this.repo ? this.repo.URL : '';
  }

  get tabs(): tabOptions[] {
    const options: tabOptions[] = [
      {
        key: 'Overview',
        link: {
          name: 'report-overview',
          query: this.$route.query
        }
      }
    ];
    if (this.report) {
      options.push({
        key: 'Code',
        link: {
          name: 'report-code',
          query: this.$route.query
        }
      });
    }
    if (this.history.length > 0) {
      options.push({
        key: 'History',
        link: {
          name: 'report-history',
          query: this.$route.query
        }
      });
    }
    if (
      this.repo &&
      this.user &&
      this.owner &&
      Object.keys(this.$route.query).length === 0
    ) {
      options.push({
        key: 'Setting',
        link: {
          name: 'report-setting'
        }
      });
    }
    return options;
  }

  get avatar(): string {
    switch (this.repo?.SCM) {
      case 'github': {
        return 'mdi-github';
      }
      case 'gitea': {
        return '$vuetify.icons.gitea';
      }
      case 'gitlab': {
        return 'mdi-gitlab';
      }
      case 'bitbucket': {
        return 'mdi-bitbucket';
      }
      default: {
        return 'mdi-source-repository';
      }
    }
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
