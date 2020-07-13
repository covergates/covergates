<template>
  <v-container>
    <v-row justify="center">
      <v-col cols="6">
        <v-card>
          <v-card-title>No Report Found</v-card-title>
          <v-card-text></v-card-text>
          <v-card-actions>
            <v-spacer />
            <v-btn :href="repoURL">Navigate Repository</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { Component } from 'vue-property-decorator';
import Vue from '@/vue';

@Component({
  name: 'report-empty'
})
export default class ReportEmpty extends Vue {
  get repository(): Repository | undefined {
    return this.$store.state.repository.current;
  }

  get repoURL(): string {
    return this.repository ? this.repository.URL : '';
  }
}
</script>

<style lang="scss" scoped>
</style>
