<template>
  <v-card>
    <v-card-title class="primary white--text">Recent Commits</v-card-title>
    <v-divider />
    <v-card-text>
      <v-list>
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

@Component
export default class RecentCommits extends Vue {
  get commits(): Commit[] {
    return this.$store.state.repository.commits.slice(0, 20);
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
}
</script>
