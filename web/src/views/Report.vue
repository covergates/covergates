<template>
  <div>
    <v-tabs v-model="tab">
      <v-tab v-for="tab in tabs" :key="tab.key" :append="true" :to="tab.link">{{tab.key}}</v-tab>
    </v-tabs>
    <v-tabs-items v-model="tab">
      <v-tab-item v-for="tab in tabs" :key="tab.key">
        <router-view></router-view>
      </v-tab-item>
    </v-tabs-items>
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
  tab = null;
  tabs = [
    {
      key: 'Overview',
      link: ''
    }
  ];
}
</script>
