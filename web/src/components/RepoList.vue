<template>
  <v-list>
    <repo-list-item v-for="(repo, index) in repos" :key="index" :repo="repo"></repo-list-item>
  </v-list>
</template>

<script lang="ts">
import { Component, Vue, Prop } from 'vue-property-decorator';
import RepoListItem from '@/components/RepoListItem.vue';
@Component({
  name: 'repo-list',
  components: {
    RepoListItem
  }
})
export default class RepoList extends Vue {
  @Prop(Array) readonly repos: Repository[] | undefined;
}
</script>

<style lang="scss">
</style>

<docs>

### Examples

```[import](./__examples__/RepoList.vue)
```

</docs>
