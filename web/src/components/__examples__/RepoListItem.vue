<template>
  <div>
    <repo-list-item :repo="repo1" />
    <repo-list-item :repo="repo2" />
    <repo-list-item :repo="repo3" />
  </div>
</template>

<script lang="ts">
import RepoListItem from '@/components/RepoListItem.vue';
export default {
  components: {
    RepoListItem
  },
  data() {
    return {
      repo1: {
        URL: 'http://gitea/org/repo',
        Name: 'repo',
        SCM: 'gitea'
      },
      repo2: {
        URL: 'http://github/org/repo',
        Name: 'repo',
        SCM: 'github'
      },
      repo3: {
        URL: 'http://unknon/org/repo',
        Name: 'repo',
        SCM: 'unkown'
      }
    };
  }
};
</script>
