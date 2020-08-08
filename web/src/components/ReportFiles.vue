<template>
  <v-container>
    <v-data-table
      class="rounded-t"
      :items="fileInfos"
      item-key="name"
      :headers="headers"
      :header-props="{
        sortIcon: 'mdi-chevron-up'
      }"
      :fixed-header="true"
      sort-by="coverage"
      :sort-desc="true"
    >
      <template v-slot:item.name="{ item }">
        <span class="file-icon fiv-viv mr-5" :class="[fileIcon(item.name)]"></span>
        <router-link :append="true" :to="fileLink(item)">{{item.name}}</router-link>
      </template>
    </v-data-table>
  </v-container>
</template>

<script lang="ts">
import { Component, Mixins } from 'vue-property-decorator';
import { Location } from 'vue-router';
import Vue from '@/vue';
import ReportMixin from '@/mixins/report';

type FileInfo = {
  name: string;
  coverage: number;
  hits: number;
};

@Component
export default class ReportFiles extends ((Mixins(ReportMixin) as typeof Vue) &&
  ReportMixin) {
  /**
   * Table Headers
   */
  headerClass = 'text-subtitle-1 align-center accent white--text';
  headers = [
    {
      text: 'File Path',
      align: 'start',
      value: 'name',
      class: this.headerClass
    },
    {
      text: 'Hit Lines',
      value: 'hits',
      width: '150px',
      class: this.headerClass
    },
    {
      text: 'Coverage',
      value: 'coverage',
      width: '150px',
      class: this.headerClass
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
          coverage: 0,
          hits: 0
        };
      }
    }
    for (const name in this.$sourceFiles) {
      const file = this.$sourceFiles[name];
      const coverage = Math.round(file.StatementCoverage * 10000) / 100;
      const hitLine = file.StatementHits.filter(hit => {
        return hit.Hits > 0;
      }).length;
      if (info[file.Name]) {
        info[file.Name].coverage = coverage;
        info[file.Name].hits = hitLine;
      } else {
        info[file.Name] = {
          name: file.Name,
          coverage: coverage,
          hits: hitLine
        };
      }
    }
    return Object.values(info).sort((a, b) => {
      return a.name > b.name ? 1 : -1;
    });
  }

  fileIcon(name: string): string {
    const ext = name.split('.').pop();
    return `fiv-icon-${ext}`;
  }

  fileLink(file: FileInfo): Location {
    return {
      path: file.name,
      query: this.$route.query
    };
  }
}
</script>

<style lang="scss" scoped>
@import '@/assets/styles/variables';

::v-deep table {
  border-collapse: collapse !important;
}

.file-icon {
  font-size: 20px;
}
</style>
