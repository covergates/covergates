<template>
  <div>
    <v-list-item>
      <v-list-item-avatar>
        <v-icon size="28">{{avatar}}</v-icon>
      </v-list-item-avatar>
      <v-list-item-content>
        <v-list-item-title>{{name}}</v-list-item-title>
        <v-list-item-subtitle>{{repoURL}}</v-list-item-subtitle>
      </v-list-item-content>
      <v-list-item-action>
        <v-btn
          ref="activate"
          small
          v-if="!activated"
          :loading="loading"
          @click="activateRepository"
          class="d-none d-md-flex"
        >Activate</v-btn>
        <v-btn
          ref="activate"
          v-if="!activated"
          :loading="loading"
          @click="activateRepository"
          class="d-flex d-md-none align-center"
          icon
        >
          <v-icon>mdi-plus-box-multiple</v-icon>
        </v-btn>
        <v-btn ref="goto" icon :to="routeLink" v-if="activated">
          <v-icon color="grey lighten-1">mdi-chevron-right</v-icon>
        </v-btn>
      </v-list-item-action>
    </v-list-item>
    <v-snackbar v-model="showSnackbar">
      {{ error }}
      <template v-slot:action="{ attrs }">
        <v-btn color="pink" text v-bind="attrs" @click="showSnackbar = undefined">Close</v-btn>
      </template>
    </v-snackbar>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Watch } from 'vue-property-decorator';
import { AxiosResponse } from 'axios';
import Vue from '@/vue';

@Component({
  name: 'repo-list-item'
})
export default class RepoListItem extends Vue {
  /**
   * repository to show in this item
   */
  @Prop(Object) readonly repo: Repository | undefined;

  protected loading = false;
  protected activated: boolean;
  protected error = '';
  protected showSnackbar: boolean;

  constructor() {
    super();
    this.activated = false;
    this.showSnackbar = false;
  }

  @Watch('repo')
  onRepoChanged() {
    this.setActivated();
  }

  mounted() {
    this.setActivated();
  }

  get avatar(): string {
    switch (this.repo?.SCM) {
      case 'github': {
        return 'mdi-github';
      }
      case 'gitea': {
        return '$vuetify.icons.gitea';
      }
      case 'gitlab': {
        return 'mdi-gitlab';
      }
      case 'bitbucket': {
        return 'mdi-bitbucket';
      }
      default: {
        return 'mdi-source-repository';
      }
    }
  }

  get base(): string {
    return this.$store.state.base;
  }

  get name(): string {
    return this.repo ? this.repo.Name : '';
  }

  get repoURL(): string {
    return this.repo ? this.repo.URL : '';
  }

  get routeLink(): string {
    if (
      this.repo === undefined ||
      this.repo.Name === undefined ||
      this.repo.NameSpace === undefined
    ) {
      return '/';
    }
    return `/report/${this.repo.SCM}/${this.repo.NameSpace}/${this.repo.Name}`;
  }

  private createRepositoryIfNotExists(): Promise<void | AxiosResponse> {
    if (this.repo === undefined) {
      return Promise.reject(new Error('repository is undefined'));
    }
    const { SCM: scm, NameSpace: namespace, Name: name, URL: url } = this.repo;
    return this.$http
      .get(`${this.base}/api/v1/repos/${scm}/${namespace}/${name}`)
      .catch(reason => {
        if (reason.response.status === 404) {
          return this.$http.post(`${this.base}/api/v1/repos`, {
            name: name,
            namespace: namespace,
            scm: scm,
            url: url
          });
        }
      });
  }

  private showError(msg: string) {
    this.error = msg;
    this.showSnackbar = true;
  }

  activateRepository() {
    if (this.repo !== undefined) {
      this.loading = true;
      const { SCM: scm, NameSpace: namespace, Name: name } = this.repo;
      this.createRepositoryIfNotExists()
        .then(() => {
          return this.$http
            .patch(
              `${this.base}/api/v1/repos/${scm}/${namespace}/${name}/report`
            )
            .then(() => {
              this.activated = true;
            });
        })
        .catch(reason => {
          this.showError(this.$httpError(reason).message);
        })
        .finally(() => {
          this.loading = false;
        });
    }
  }

  setActivated() {
    this.activated =
      this.repo !== undefined &&
      this.repo.ReportID !== undefined &&
      this.repo.ReportID !== '';
  }
}
</script>

<style lang="scss" scoped>
</style>>

<docs>

### Examples

```[import](./__examples__/RepoListItem.vue)
```

</docs>
