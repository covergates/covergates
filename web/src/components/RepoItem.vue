<template>
    <div>
        {{this.repo.URL}}
    </div>
</template>

<script lang="ts">

import { Component, Vue, Prop } from 'vue-property-decorator';

@Component
export default class RepoItem extends Vue {
    /**
     * repository to show in this item
     */
    @Prop(Object) readonly repo: Repository | undefined;
}

</script>

<style lang="sass" scoped>

</style>>

<docs>

## Examples

```[import](./__examples__/RepoItem.vue)
```

</docs>
