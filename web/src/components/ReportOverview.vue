<template>
  <v-container>
    <v-banner single-line>
      <span class="text-h4">Files</span>
    </v-banner>
    <v-sheet class="content">
      <ICountUp :endVal="filesCount" class="count-up text-h2" />
    </v-sheet>
    <v-banner single-line>
      <span class="text-h4">Coverage</span>
    </v-banner>
    <v-sheet class="content">
      <v-progress-circular
        :size="100"
        :width="15"
        :rotate="-90"
        :value="coverage"
        color="primary"
      >{{coverage}}</v-progress-circular>
    </v-sheet>
  </v-container>
</template>

<script lang="ts">
import { Component, Watch } from 'vue-property-decorator';
import ICountUp from 'vue-countup-v2';
import Vue from '@/vue';

@Component({
  name: 'report-overview',
  components: {
    ICountUp
  }
})
export default class ReportOverview extends Vue {
  filesCount = 0;

  @Watch('repository', { immediate: true })
  onRepositoryChange() {
    this.updateFileCount(this.repository);
  }

  mounted() {
    this.updateFileCount(this.repository);
  }

  get repository(): Repository | undefined {
    return this.$store.state.repository.current;
  }

  get coverage(): number {
    const report = this.$store.state.report.current;
    if (report !== undefined && report.coverage !== undefined) {
      return report.coverage.StatementCoverage * 100;
    }
    return 0;
  }

  updateFileCount(repo: Repository | undefined) {
    if (repo === undefined) {
      return;
    }
    const scm = repo.SCM;
    const name = `${repo.NameSpace}/${repo.Name}`;
    this.$http
      .get<string[]>(
        `${this.$store.state.base}/api/v1/repos/${scm}/${name}/files`
      )
      .then(response => {
        this.filesCount = response.data.length;
      })
      .catch(error => {
        console.warn(error);
        this.filesCount = 0;
      });
  }
}
</script>

<style lang="scss" scoped>
@import '@/assets/styles/variables';

.content {
  padding: 20px;
}

.count-up {
  color: $content-color;
}
</style>

<docs>

### Examples
```[import](./__examples__/ReportOverview.vue)
```
</docs>
