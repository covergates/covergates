<template>
  <div>
    <v-card flat>
      <v-card-title class="text-h6">Step 1. Get the binary</v-card-title>
      <v-divider />
      <v-card-text class="py-2">
        Download
        <a href="https://github.com/covergates/covergates/releases">covergates binary</a>
        for your platform. You will find
        <span class="font-weight-bold">covergates</span>
        (or
        <span class="font-weight-bold">covergates.exe</span> for Windows) in the package.
      </v-card-text>
    </v-card>

    <v-card flat>
      <v-card-title class="text-h6">Step 2. Update Setting</v-card-title>
      <v-divider />
      <v-card-text class="py-2">
        You can change covergates default behavior by update
        <router-link :to="{name: 'report-setting'}">setting</router-link>.
        Available options are:
        <ul>
          <li>
            Apply path
            <span class="font-weight-bold">filter</span>, which trims the file paths in coverage report.
          </li>
          <li>Activate Webhook</li>
        </ul>
      </v-card-text>
    </v-card>

    <v-card flat>
      <v-card-title class="text-h6">Step 3. Run Testing</v-card-title>
      <v-divider />
      <v-card-text class="py-2">Test your codes and collect coverage reports.</v-card-text>
    </v-card>

    <v-card flat>
      <v-card-title class="text-h6">Step 4. Commit Report</v-card-title>
      <v-divider />
      <v-card-text class="py-2">
        Almost there! Copy your report ID:
        <v-btn x-small class="mx-1" v-clipboard:copy="reportID">{{reportID}}</v-btn>
        <br />Run:
        <code class="mx-2">coveragtes upload --report "Report ID" --type "Report Type" "Report"</code>
      </v-card-text>
    </v-card>
  </div>
</template>

<script lang="ts">
import { Component } from 'vue-property-decorator';
import Vue from '@/vue';

type StepGuide = {
  title: string;
  body: string;
};

@Component
export default class ReportGuide extends Vue {
  get repo(): Repository | undefined {
    return this.$store.state.repository.current;
  }

  get reportID(): string {
    return this.repo && this.repo.ReportID ? this.repo.ReportID : '';
  }

  get reportURL(): string {
    const base = this.$store.state.base;
    if (!this.repo) {
      return '';
    }
    const { SCM: scm, Name: name, NameSpace: space } = this.repo;
    return `${base}/report/${scm}/${space}/${name}`;
  }
}
</script>
