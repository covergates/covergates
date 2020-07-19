<template>
  <v-container>
    <v-card>
      <v-card-title>Filters</v-card-title>
      <v-card-text>
        <v-textarea name="filters" v-model="filters" :hint="hint" flat outlined></v-textarea>
      </v-card-text>
      <v-divider></v-divider>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn class="mr-5" text @click="saveSetting" :loading="loading" :disabled="!repo">save</v-btn>
      </v-card-actions>
    </v-card>
  </v-container>
</template>

<script lang="ts">
import { Component, Watch } from 'vue-property-decorator';
import Vue from '@/vue';

const defaultHint = `
Filter will remove pattern from the path.
Provide a regular expression each line.
`;

@Component({
  name: 'report-setting'
})
export default class ReportSetting extends Vue {
  private filters: string;
  private hint: string;
  private loading: boolean;
  constructor() {
    super();
    this.filters = '';
    this.hint = defaultHint;
    this.loading = false;
  }

  mounted() {
    this.filters = this.filterText;
  }

  get repo(): Repository | undefined {
    return this.$store.state.repository.current;
  }

  get setting(): RepositorySetting | undefined {
    return this.$store.state.repository.setting;
  }

  get filterText(): string {
    if (this.setting && this.setting.filters) {
      return this.setting.filters.join('\n');
    } else {
      return '';
    }
  }

  @Watch('setting', { deep: true })
  onSettingUpdate() {
    this.filters = this.filterText;
  }

  saveSetting() {
    if (this.repo === undefined) {
      return;
    }
    const base = this.$store.state.base;
    const setting = this.setting
      ? this.setting
      : ({ filters: [] } as RepositorySetting);
    const { SCM, NameSpace, Name } = this.repo;
    setting.filters = this.filters.trim().split('\n');
    this.loading = true;
    this.$http
      .post(`${base}/api/v1/repos/${SCM}/${NameSpace}/${Name}/setting`, setting)
      .finally(() => {
        this.loading = false;
      });
  }
}
</script>
