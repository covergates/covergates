import { Component } from 'vue-property-decorator';
import Vue from '@/vue';

@Component
export default class ReportMixin extends Vue {
  get $report(): Report | undefined {
    return this.$store.state.report.current;
  }

  get $coverage(): number {
    if (this.$report === undefined) {
      return 0;
    }
    if (this.$report.coverages === undefined || this.$report.coverages.length <= 0) {
      return 0;
    }
    let sum = 0;
    for (const coverage of this.$report.coverages) {
      sum += coverage.statementCoverage;
    }
    return Math.round(sum / this.$report.coverages.length * 10000) / 100;
  }

  get $sourceFiles(): Record<string, SourceFile> {
    if (this.$report?.coverages === undefined) {
      return {};
    }
    const records = {} as Record<string, SourceFile>;
    for (const coverage of this.$report?.coverages) {
      if (coverage.files === undefined) {
        continue;
      }
      for (const file of coverage.files) {
        records[file.Name] = file;
      }
    }
    return records;
  }

  findSourceFile(name: string): SourceFile | undefined {
    return this.$sourceFiles[name];
  }
}
