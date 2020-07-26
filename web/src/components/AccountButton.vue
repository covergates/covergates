<template>
  <v-menu left bottom offset-y>
    <template v-slot:activator="{ on, attrs }">
      <v-btn class="mr-5" v-bind="attrs" v-on="on" icon>
        <v-avatar size="36">
          <v-icon v-if="avatar===''">mdi-account</v-icon>
          <img :src="avatar" :alt="`${name}-avatar`" v-else />
        </v-avatar>
      </v-btn>
    </template>
    <v-list>
      <v-list-item v-for="(item, index) in actions" :key="index" dense @click="actionClick(item)">
        <v-list-item-icon>
          <v-icon>{{item.icon}}</v-icon>
        </v-list-item-icon>
        <v-list-item-title>{{ item.name }}</v-list-item-title>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script lang="ts">
import { Component } from 'vue-property-decorator';
import Vue from '@/vue';

type actionItem = {
  name: string;
  icon: string;
  to: string;
};

@Component({
  name: 'account-button'
})
export default class AccountButton extends Vue {
  get user(): User | undefined {
    return this.$store.state.user.current;
  }

  get avatar(): string {
    return this.user && this.user.avatar ? this.user.avatar : '';
  }

  get name(): string {
    return this.user && this.user.login ? this.user.login : 'user';
  }

  get actions(): actionItem[] {
    const items = [] as actionItem[];
    if (this.user) {
      items.push({
        name: 'Setting',
        icon: 'mdi-cog',
        to: '/user'
      });
      items.push({
        name: 'Logout',
        icon: 'mdi-logout',
        to: '/logoff'
      });
    } else {
      items.push({
        name: 'Login',
        icon: 'mdi-login',
        to: '/login'
      });
    }
    return items;
  }

  actionClick(item: actionItem) {
    if (item.name === 'Logout') {
      const base = this.$store.state.base;
      window.location.href = `${base}${item.to}`;
    } else {
      this.$router.push(item.to);
    }
  }
}
</script>
