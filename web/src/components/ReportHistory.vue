<template>
  <v-container>
    <v-row>
      <v-col sm="12" md="6">
        <recent-commits />
      </v-col>
      <v-col sm="12" md="6">
        <report-timeline />
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { Component } from 'vue-property-decorator';
import Vue from '@/vue';
import RecentCommits from '@/components/RecentCommits.vue';
import ReportTimeline from '@/components/ReportTimeline.vue';

@Component({
  components: {
    RecentCommits,
    ReportTimeline
  }
})
export default class ReportHistory extends Vue {}
</script>
