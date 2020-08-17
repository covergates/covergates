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
            hint="Use this ID to upload report. This ID is used only to identify your report. It can't be used to access repository information."
            outlined
            dense
            flat
            persistent-hint
          ></v-text-field>
        </v-row>
        <v-row class="d-flex mt-5">
          <v-text-field label="Badge" :readonly="true" :value="badge" outlined dense flat>
            <template v-slot:append>
              <img :src="badge" alt="badge" />
            </template>
          </v-text-field>
        </v-row>
        <v-row>
          <v-textarea
            name="markdown"
            :readonly="true"
            v-model="markdown"
            hint="copy the text and paste to README.md"
            rows="2"
            flat
            no-resize
            outlined
            persistent-hint
          ></v-textarea>
        </v-row>
        <v-row class="d-flex mt-5">
          <v-text-field label="Card" :readonly="true" :value="card" outlined flat>
            <template v-slot:append>
              <v-img :src="card" alt="card" class="mb-3 d-none d-lg-block" />
            </template>
          </v-text-field>
        </v-row>
        <v-row class="d-flex d-lg-none justify-center">
          <v-col cols="12" sm="8">
            <v-img :src="card" alt="card" />
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
    <v-card flat>
      <v-card-title>Webhook</v-card-title>
      <v-divider />
      <v-card-text class="d-flex align-center">
        <hook-button />
        <v-switch
          :loading="loading"
          value
          v-model="autoMerge"
          label="Auto Merge Report"
          class="mx-5"
          @change="saveAutoMerge"
        ></v-switch>
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
        <v-btn class="mr-5" @click="saveFilters" :loading="loading" :disabled="!repo" small>save</v-btn>
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
  private autoMerge: boolean;
  private loading: boolean;
  constructor() {
    super();
    this.filters = '';
    this.hint = defaultHint;
    this.loading = false;
    this.autoMerge = false;
  }

  mounted() {
    this.filters = this.filterText;
    this.autoMerge =
      this.setting && this.setting.mergePR ? this.setting.mergePR : false;
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

  get url(): string {
    let base = this.$store.state.base;
    if (base !== '') {
      base = base.replace(/\/$/, '');
      base = base.replace(/^\//, '');
      base = `/${base}`;
    }
    const protocol = window.location.protocol;
    const host = window.location.host;
    return `${protocol}//${host}${base}`;
  }

  get badge(): string {
    if (this.repo) {
      return `${this.url}/api/v1/reports/${this.repo.ReportID}/badge`;
    }
    return '';
  }

  get card(): string {
    if (this.repo) {
      return `${this.url}/api/v1/reports/${this.repo.ReportID}/card`;
    }
    return '';
  }

  get markdown(): string {
    if (this.repo) {
      const { SCM: scm, NameSpace: space, Name: name } = this.repo;
      return `[![badge](${this.badge})](${this.url}/report/${scm}/${space}/${name})`;
    }
    return '';
  }

  @Watch('setting', { deep: true })
  onSettingUpdate() {
    this.filters = this.filterText;
    this.autoMerge =
      this.setting && this.setting.mergePR ? this.setting.mergePR : false;
  }

  saveFilters() {
    const setting = this.setting
      ? this.setting
      : ({ filters: [] } as RepositorySetting);
    setting.filters = this.filters.trim().split('\n');
    this.saveSetting(setting);
  }

  saveAutoMerge() {
    const setting = this.setting
      ? this.setting
      : ({ mergePR: false } as RepositorySetting);
    setting.mergePR = this.autoMerge;
    this.saveSetting(setting);
  }

  saveSetting(setting: RepositorySetting) {
    if (this.repo === undefined) {
      return;
    }
    const base = this.$store.state.base;
    const { SCM, NameSpace, Name } = this.repo;
    this.loading = true;
    this.$http
      .post(`${base}/api/v1/repos/${SCM}/${NameSpace}/${Name}/setting`, setting)
      .finally(() => {
        this.loading = false;
      });
  }
}
</script>
