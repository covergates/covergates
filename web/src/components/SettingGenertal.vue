<template>
  <div>
    <setting-info />
    <v-card flat>
      <v-card-title>General</v-card-title>
      <v-divider />
      <v-card-text>
        <v-simple-table>
          <template v-slot:default>
            <tbody>
              <tr>
                <td>Project Setting</td>
                <td class="d-flex align-center">
                  <v-switch
                    :loading="loading"
                    value
                    v-model="projectProtected"
                    label="Protected"
                    @change="saveProtected"
                  ></v-switch>
                  <span class="mx-5 text-caption">Only authorized user can upload report</span>
                </td>
              </tr>
              <tr>
                <td>Project Webhooks</td>
                <td class="d-flex align-center">
                  <hook-button />
                  <v-switch
                    :loading="loading"
                    value
                    v-model="autoMerge"
                    label="Auto Merge Report"
                    class="mx-5"
                    @change="saveAutoMerge"
                  ></v-switch>
                </td>
              </tr>
            </tbody>
          </template>
        </v-simple-table>
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
  </div>
</template>

<script lang="ts">
import { Component, Watch } from 'vue-property-decorator';
import Vue from '@/vue';
import HookButton from '@/components/HookButton.vue';
import SettingInfo from '@/components/SettingInfo.vue';

const defaultHint = `
Filter will remove pattern from the path.
Provide a regular expression each line.
`;

@Component({
  name: 'setting-general',
  components: {
    HookButton,
    SettingInfo
  }
})
export default class SettingGeneral extends Vue {
  private filters: string;
  private hint: string;
  private autoMerge: boolean;
  private projectProtected: boolean;
  private loading: boolean;
  constructor() {
    super();
    this.filters = '';
    this.hint = defaultHint;
    this.loading = false;
    this.autoMerge = false;
    this.projectProtected = false;
  }

  mounted() {
    this.syncSetting();
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
    this.syncSetting();
  }

  syncSetting() {
    this.filters = this.filterText;
    if (this.setting) {
      this.autoMerge =
        this.setting.mergePR !== undefined ? this.setting.mergePR : false;
      this.projectProtected =
        this.setting.protected !== undefined ? this.setting.protected : false;
    }
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

  saveProtected() {
    const setting = this.setting
      ? this.setting
      : ({ projectProtected: false } as RepositorySetting);
    setting.protected = this.projectProtected;
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
