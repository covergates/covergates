<template>
  <repo-list :repos="repositories"></repo-list>
</template>

<script lang="ts">
import RepoList from '@/components/RepoList.vue';
export default {
  components: {
    RepoList
  },
  data() {
    return {
      repositories: [
        {
          URL: 'http://repo1',
          Name: 'repo1'
        },
        {
          URL: 'http://repo2',
          Name: 'repo2'
        }
      ]
    };
  }
};
</script>
