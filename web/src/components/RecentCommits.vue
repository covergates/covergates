<template>
  <v-card>
    <v-card-title class="primary white--text">
      Recent Commits
      <v-spacer />
      <v-select
        :items="branches"
        :hide-details="true"
        v-model="selectedBranch"
        flat
        dense
        solo
        @change="selectBranch"
      ></v-select>
    </v-card-title>
    <v-divider />
    <v-card-text>
      <v-card flat v-if="commits.length <= 0">
        <v-card-title>
          <v-icon size="36" class="mr-2">mdi-progress-question</v-icon>No Commits Found
        </v-card-title>
        <v-card-text class="px-5">{{hint}}</v-card-text>
      </v-card>
      <v-list v-else>
        <v-list-item v-for="commit in commits" :key="commit.sha">
          <v-list-item-avatar color="accent">
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
    this.$store.dispatch(Actions.FETCH_REPOSITORY_COMMITS, ref);
  }
}
</script>
