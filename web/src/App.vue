<template>
  <v-app id="inspire">
    <v-app id="keep">
      <v-app-bar
        app
        color="primary"
        dark
        clipped-left
      >
        <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>
        <v-spacer></v-spacer>

      </v-app-bar>

      <v-navigation-drawer
        v-model="drawer"
        clipped
        app
      >
      </v-navigation-drawer>

      <v-main>
        <v-container fluid>
          <router-view></router-view>
        </v-container>
      </v-main>
    </v-app>
  </v-app>
</template>

<script lang="ts">
import { Vue, Component } from 'vue-property-decorator'

@Component
export default class AppComponent extends Vue {
  private drawer = true
}

</script>
