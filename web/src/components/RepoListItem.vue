<template>
  <div>
    <v-list-item>
      <v-list-item-avatar>
        <v-icon>{{avatar}}</v-icon>
      </v-list-item-avatar>
      <v-list-item-content>
        <v-list-item-title>{{name}}</v-list-item-title>
        <v-list-item-subtitle>{{repoURL}}</v-list-item-subtitle>
      </v-list-item-content>
      <v-list-item-action>
        <v-btn icon :to="routeLink">
          <v-icon color="grey lighten-1">mdi-chevron-right</v-icon>
        </v-btn>
      </v-list-item-action>
    </v-list-item>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from 'vue-property-decorator';

@Component({
  name: 'repo-list-item'
})
export default class RepoListItem extends Vue {
  /**
   * repository to show in this item
   */
  @Prop(Object) readonly repo: Repository | undefined;

  get avatar(): string {
    switch (this.repo?.SCM) {
      case 'github': {
        return 'mdi-github';
      }
      case 'gitea': {
        return '$vuetify.icons.gitea';
      }
      default: {
        return 'mdi-source-repository';
      }
    }
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
    return `/${this.repo.SCM}/${this.repo.NameSpace}/${this.repo.Name}`;
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
