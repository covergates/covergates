<template>
  <div></div>
</template>

<script lang="ts">
import { Component } from 'vue-property-decorator';
import Vue from '@/vue';

type ReportBrief = {
  name: string;
  coverage: string;
};

@Component
export default class ReportCode extends Vue {}
</script>
