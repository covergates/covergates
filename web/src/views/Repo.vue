<template>
  <perfect-scrollbar class="page-container">
    <v-container>
      <v-row align="top" justify="center" class="fill-height">
        <v-col cols="12" v-show="loading || emptyRepo"></v-col>
        <v-col cols="6" v-show="loading">
          <p class="subtitle-1 my-1">Getting Repositories</p>
          <v-progress-linear indeterminate rounded height="6"></v-progress-linear>
        </v-col>
        <v-col sm="12" md="8" v-show="!loading && !emptyRepo">
          <v-card flat>
            <v-toolbar flat>
              <v-text-field
                v-model="searchText"
                class="search-bar"
                label="Search"
                solo
                dense
                height="48"
              ></v-text-field>
              <v-btn :loading="syncing" color="accent" class="ml-5" small dark @click="synchronize">
                <v-icon class="mr-1" small>mdi-sync</v-icon>Sync
              </v-btn>
            </v-toolbar>
            <v-card-text>
              <repo-list v-if="!loading" :repos="repositories"></repo-list>
            </v-card-text>
          </v-card>
        </v-col>
        <v-col cols="12" v-show="!loading && emptyRepo" class="subtitle-1 text-center">
          Oops, there is no repository found. Try to
          <v-btn :loading="syncing" color="accent" class="mx-1" small dark @click="synchronize">
            <v-icon class="mr-1" small>mdi-sync</v-icon>Sync
          </v-btn>at first.
        </v-col>
      </v-row>
    </v-container>
  </perfect-scrollbar>
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
  syncing = false;
  searchText = '';
  mounted() {
    this.$store.dispatch(Actions.FETCH_REPOSITORY_LIST);
  }

  get loading(): boolean {
    return this.$store.state.repository.loading;
  }

  get emptyRepo(): boolean {
    return this.$store.state.repository.list.length <= 0;
  }

  get repositories(): Repository[] {
    let repos = [] as Repository[];
    repos.push(...this.$store.state.repository.list);
    repos.sort((a, b) => {
      if (a.ReportID && b.ReportID) {
        return a.URL.localeCompare(b.URL);
      } else if (a.ReportID) {
        return -1;
      } else if (b.ReportID) {
        return 1;
      } else {
        return a.URL.localeCompare(b.URL);
      }
    });
    if (this.searchText !== '') {
      const text = this.searchText.trim();
      repos = repos.filter(repo => {
        return repo.URL.toLowerCase().includes(text.toLowerCase());
      });
    }
    return repos;
  }

  synchronize() {
    this.syncing = true;
    this.$store.dispatch(Actions.SYNCHRONIZE_REPOSITORY).finally(() => {
      this.syncing = false;
    });
  }
}
</script>

<style lang="scss" scoped>
.page-container {
  overflow-y: auto;
  height: 100%;
}

::v-deep .search-bar {
  .v-input__control {
    height: 48px !important;
  }
}
</style>
