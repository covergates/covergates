<template>
  <v-container>
    <v-data-table :items="fileInfos" item-key="name" :headers="headers" :fixed-header="true">
      <template v-slot:item.name="{ item }">
        <router-link :append="true" :to="item.name">{{item.name}}</router-link>
      </template>
    </v-data-table>
  </v-container>
</template>

<script lang="ts">
import { Component } from 'vue-property-decorator';
import Vue from '@/vue';

type FileInfo = {
  name: string;
  coverage: number;
};

@Component
export default class ReportCode extends Vue {
  headers = [
    {
      text: 'File Path',
      align: 'start',
      value: 'name'
    },
    {
      text: 'Coverage',
      value: 'coverage'
    }
  ];

  get fileInfos(): FileInfo[] {
    const repo = this.$store.state.repository.current;
    const report = this.$store.state.report.current;
    if (repo === undefined || report === undefined) {
      return [];
    }
    const info = {} as { [key: string]: FileInfo };
    if (repo.Files) {
      for (const file of repo.Files) {
        info[file] = {
          name: file,
          coverage: 0
        };
      }
    }
    if (report.coverage && report.coverage.Files) {
      for (const file of report.coverage.Files) {
        if (info[file.Name]) {
          info[file.Name].coverage = file.StatementCoverage;
        } else {
          info[file.Name] = {
            name: file.Name,
            coverage: file.StatementCoverage
          };
        }
      }
    }
    return Object.values(info).sort((a, b) => {
      return a.name > b.name ? 1 : -1;
    });
  }
}
</script>

<style lang="scss" scoped>
</style>
