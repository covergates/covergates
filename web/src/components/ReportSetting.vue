<template>
  <v-container>
    <v-card flat>
      <v-card-title>General</v-card-title>
      <v-divider />
      <v-card-text>
        <v-row>
          <v-text-field
            label="Report ID"
            :readonly="true"
            :value="repo.ReportID"
            outlined
            dense
            flat
          ></v-text-field>
        </v-row>
        <v-row>
          <span>Use this ID to upload report. This ID is used only to identify your report. It can't be used to access repository information.</span>
          <v-spacer></v-spacer>
        </v-row>
      </v-card-text>
    </v-card>
    <v-card flat>
      <v-card-title>Webhook</v-card-title>
      <v-divider />
      <v-card-text>
        <hook-button />
      </v-card-text>
    </v-card>
    <v-card flat>
      <v-card-title>Filters</v-card-title>
      <v-divider />
      <v-card-text>
        <v-textarea name="filters" v-model="filters" :hint="hint" flat outlined></v-textarea>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn class="mr-5" @click="saveSetting" :loading="loading" :disabled="!repo" small>save</v-btn>
      </v-card-actions>
    </v-card>
  </v-container>
</template>

<script lang="ts">
import { Component, Watch } from 'vue-property-decorator';
import Vue from '@/vue';
import HookButton from '@/components/HookButton.vue';

const defaultHint = `
Filter will remove pattern from the path.
Provide a regular expression each line.
`;

@Component({
  name: 'report-setting',
  components: {
    HookButton
  }
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
