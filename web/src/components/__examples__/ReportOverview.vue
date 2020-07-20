<template>
  <report-overview />
</template>

<script lang="js">
import ReportOverview from '@/components/ReportOverview.vue';
import { Mutations } from '@/store';

export default {
  components: { ReportOverview },
  mounted: function () {
    this.$store.commit(Mutations.SET_REPORT_CURRENT, {
      coverage: {
        StatementCoverage: 50
      }
    });
  }
};
</script>
