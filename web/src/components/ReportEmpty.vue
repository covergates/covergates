<template>
  <v-container>
    <v-row justify="center">
      <v-col sm="12" md="8">
        <v-card>
          <v-card-title class="justify-center accent white--text">
            <v-icon size="32" class="mr-5" dark>mdi-party-popper</v-icon>Commit Your First Report!
          </v-card-title>
          <v-card-text>
            <report-guide />
          </v-card-text>
          <v-card-actions>
            <v-spacer />
            <v-btn small color="primary" :href="repoURL">Navigate Repository</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { Component } from 'vue-property-decorator';
import Vue from '@/vue';
import ReportGuide from '@/components/ReportGuide.vue';

@Component({
  name: 'report-empty',
  components: {
    ReportGuide
  }
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
