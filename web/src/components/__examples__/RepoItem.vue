<template>
  <RepoItem :repo="repo"/>
</template>

<script lang="ts">
import RepoItem from '@/components/RepoItem.vue';
export default {
  components: {
    RepoItem
  },
  data() {
    return {
      repo: {
        URL: 'bear'
      }
    };
  }
};
</script>
