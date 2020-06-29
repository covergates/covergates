<template>
  <v-container>
    <v-row>
      <v-col md="auto" class="d-flex align-stretch">
        <v-card width="200">
          <v-card-title>Files</v-card-title>
          <v-card-text>
            <h1>{{filesCount}}</h1>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col md="auto">
        <v-card>
          <v-card-title>Test Coverage</v-card-title>
          <v-card-text>
            <v-progress-circular size="100" width="15" :value="coverage">{{coverage}}</v-progress-circular>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from '@/vue';
import { Component } from 'vue-property-decorator';

@Component({
  name: 'report-overview'
})
export default class ReportOverview extends Vue {
  get filesCount(): number {
    return 10;
  }

  get coverage(): number {
    return 80;
  }
}
</script>

<docs>

### Examples
```[import](./__examples__/ReportOverview.vue)
```
</docs>
