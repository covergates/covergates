<template>
  <v-container>
    <RepoListItem :repo="{URL:'bear'}" />
  </v-container>
</template>

<script lang="ts">
import Vue from 'vue';
import RepoListItem from '@/components/RepoListItem.vue';
export default Vue.extend({
  name: 'HelloWorld',
  components: {
    RepoListItem
  }
});
</script>
