<template>
  <v-container>
    <v-row align="center" justify="center">
      <v-col cols="6">
        <repo-list v-if="!loading" :repos="repositories"></repo-list>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { Component } from 'vue-property-decorator';
import Vue from '@/vue';
import RepoList from '@/components/RepoList.vue';
import { Actions } from '@/store';

@Component({
  components: {
    RepoList
  }
})
export default class Repo extends Vue {
  mounted() {
    this.$store.dispatch(Actions.FETCH_REPOSITORY_LIST);
  }

  get loading(): boolean {
    return this.$store.state.repository.loading;
  }

  get repositories(): Repository[] {
    return this.$store.state.repository.list;
  }
}
</script>
