<template>
  <v-tabs :vertical="verticalTabs" class="mt-1" slider-color="accent">
    <v-tab class="accent--text">
      <v-icon left>mdi-cog</v-icon>General
    </v-tab>
    <v-tab class="accent--text">
      <v-icon left>mdi-key</v-icon>Tokens
    </v-tab>

    <v-tab-item>
      <setting-general />
    </v-tab-item>
    <v-tab-item>
      <setting-tokens />
    </v-tab-item>
  </v-tabs>
</template>

<script lang="ts">
import { Component } from 'vue-property-decorator';
import Vue from '@/vue';
import SettingGeneral from '@/components/SettingGenertal.vue';
import SettingTokens from '@/components/SettingTokens.vue';

@Component({
  name: 'report-setting',
  components: {
    SettingGeneral,
    SettingTokens
  }
})
export default class ReportSetting extends Vue {
  get verticalTabs(): boolean {
    switch (this.$vuetify.breakpoint.name) {
      case 'xs':
        return false;

      case 'sm':
        return false;

      default:
        return true;
    }
  }
}
</script>
