<template>
  <v-card flat>
    <v-card-title>Information</v-card-title>
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
</template>

<script lang="ts">
import { Component } from 'vue-property-decorator';
import Vue from '@/vue';

@Component({
  name: 'setting-info'
})
export default class SettingInfo extends Vue {
  get repo(): Repository | undefined {
    return this.$store.state.repository.current;
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
}
</script>
