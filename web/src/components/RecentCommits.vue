<template>
  <v-card>
    <v-card-title class="primary white--text flex-lg-nowrap">
      <span class="title">Recent Commits</span>
      <v-spacer />
      <v-select
        :items="branches"
        :hide-details="true"
        :disabled="loading"
        v-model="selectedBranch"
        flat
        dense
        solo
        @change="selectBranch"
      ></v-select>
    </v-card-title>
    <v-divider />
    <v-skeleton-loader :loading="loading" type="list-item-avatar-two-line">
      <v-card-text>
        <v-card flat v-if="commits.length <= 0">
          <v-card-title>
            <v-icon size="36" class="mr-2">mdi-progress-question</v-icon>No Commits Found
          </v-card-title>
          <v-card-text class="px-5">{{hint}}</v-card-text>
        </v-card>
        <v-list v-else>
          <v-list-item v-for="commit in commits" :key="commit.sha">
            <v-list-item-avatar class="elevation-4">
              <v-img :src="commit.committerAvatar" v-if="commit.committerAvatar"></v-img>
              <v-icon dark v-else>mdi-account</v-icon>
            </v-list-item-avatar>
            <v-list-item-content>
              <v-list-item-title>{{commit.message}}</v-list-item-title>
              <v-list-item-subtitle>
                {{commit.committer}}
                <v-chip
                  color="accent"
                  class="ml-5 px-1"
                  outlined
                  label
                  pill
                  x-small
                  dark
                >{{shortSHA(commit.sha)}}</v-chip>
              </v-list-item-subtitle>
            </v-list-item-content>
            <v-list-item-action>
              <v-btn small color="accent" :to="commitLink(commit)">Report</v-btn>
            </v-list-item-action>
          </v-list-item>
        </v-list>
      </v-card-text>
    </v-skeleton-loader>
  </v-card>
</template>

<script lang="ts">
import { Component } from 'vue-property-decorator';
import { Location } from 'vue-router';
import Vue from '@/vue';
import { Actions } from '@/store';

@Component
export default class RecentCommits extends Vue {
  selectedBranch = '';
  loadingCommits = false;

  constructor() {
    super();
    this.selectedBranch = this.currentBranch;
    this.updateCommits(this.currentBranch);
  }

  get commits(): Commit[] {
    return this.$store.state.repository.commits.slice(0, 20);
  }

  get repo(): Repository | undefined {
    return this.$store.state.repository.current;
  }

  get loading(): boolean {
    return this.$store.state.repository.loading || this.loadingCommits;
  }

  get currentBranch(): string {
    const report = this.$store.state.report.current;
    if (report && report.reference && report.reference !== '') {
      return report.reference;
    }
    return this.repo ? this.repo.Branch : '';
  }

  get branches(): string[] {
    return this.$store.state.repository.branches;
  }

  get user(): User | undefined {
    return this.$store.state.user.current;
  }

  get hint(): string {
    if (this.commits.length > 0) {
      return '';
    }
    return this.user
      ? 'Push your first commit and go back later!'
      : 'Please login to get the commits list';
  }

  shortSHA(sha: string): string {
    return sha.substring(0, 16);
  }

  commitLink(commit: Commit): Location {
    return {
      name: 'report-overview',
      query: {
        ref: commit.sha
      }
    };
  }

  selectBranch() {
    this.updateCommits(this.selectedBranch);
  }

  updateCommits(ref = '') {
    this.loadingCommits = true;
    this.$store.dispatch(Actions.FETCH_REPOSITORY_COMMITS, ref).finally(() => {
      this.loadingCommits = false;
    });
  }
}
</script>

<style lang="scss" scoped>
.title {
  min-width: 250px;
}
</style>
