<template>
  <v-container>
    <v-data-table
      :items="fileInfos"
      item-key="name"
      :headers="headers"
      :fixed-header="true"
      sort-by="coverage"
      :sort-desc="true"
    >
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
  /**
   * Table Headers
   */
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

  get report(): Report | undefined {
    return this.$store.state.report.current;
  }

  get fileInfos(): FileInfo[] {
    if (!this.report) {
      return [];
    }
    const info = {} as { [key: string]: FileInfo };
    if (this.report.files) {
      for (const file of this.report.files) {
        info[file] = {
          name: file,
          coverage: 0
        };
      }
    }
    if (this.report.coverage && this.report.coverage.Files) {
      for (const file of this.report.coverage.Files) {
        const coverage = Math.round(file.StatementCoverage * 10000) / 100;
        if (info[file.Name]) {
          info[file.Name].coverage = coverage;
        } else {
          info[file.Name] = {
            name: file.Name,
            coverage: coverage
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
