<template>
  <div>{{sourceCode}}</div>
</template>

<script lang="ts">
import { Component } from 'vue-property-decorator';
import Vue from '@/vue';

@Component
export default class ReportSource extends Vue {
  get sourceCode() {
    return this.$store.state.report.source;
  }
}
</script>
