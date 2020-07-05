<template>
  <div>
    <v-tabs>
      <v-tab v-for="tab in tabs" :key="tab.key" :to="tab.link">{{tab.key}}</v-tab>
    </v-tabs>
    <router-view></router-view>
  </div>
</template>

<script lang="ts">
import { Component } from 'vue-property-decorator';
import Vue from '@/vue';
import ReportOverview from '@/components/ReportOverview.vue';

@Component({
  components: { ReportOverview }
})
export default class Report extends Vue {
  tabs = [
    {
      key: 'Overview',
      link: {
        name: 'report-overview'
      }
    },
    {
      key: 'Code',
      link: {
        name: 'report-code'
      }
    }
  ];
}
</script>
