<template>
  <v-card flat>
    <v-card-title>OAuth Tokens</v-card-title>
    <v-divider />
    <v-card-text>
      <v-text-field label="Token Name" v-model="tokenName" outlined dense flat>
        <template v-slot:append-outer>
          <v-btn
            small
            color="accent"
            :disabled="!tokenName"
            :loading="loading"
            @click="generateToken"
          >Generate</v-btn>
        </template>
      </v-text-field>
      <v-sheet
        class="text-center text-subtitle-1 my-1 rounded"
        v-show="accessToken"
        color="primary"
        dark
      >{{accessToken}}</v-sheet>
      <v-divider />
      <v-list>
        <v-list-item v-for="(token, i) in tokens" :key="i">
          <v-list-item-icon>
            <v-icon color="primary">mdi-send</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>{{token.name}}</v-list-item-title>
            <v-list-item-subtitle>{{createTime(token)}}</v-list-item-subtitle>
          </v-list-item-content>
          <v-list-item-action>
            <v-btn small color="error" @click="clickedDelete(token)">DELETE</v-btn>
          </v-list-item-action>
        </v-list-item>
      </v-list>
    </v-card-text>
    <v-overlay :value="overlay">
      <v-card flat v-if="selectedToken" color="accent" :loading="loading">
        <v-card-title>Delete Token: {{selectedToken.name}}</v-card-title>
        <v-card-text>Deleting a token will revoke access to your account for applications using it. Continue?</v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn text color="error" @click="cancelDelete">Cancel</v-btn>
          <v-btn text @click="deleteSelectedToken" :loading="loading">Continue</v-btn>
        </v-card-actions>
      </v-card>
    </v-overlay>
  </v-card>
</template>

<script lang="ts">
import { Component } from 'vue-property-decorator';
import Vue from '@/vue';
import { Actions } from '@/store';

@Component({
  name: 'setting-tokens'
})
export default class SettingTokens extends Vue {
  // TODO: Add unittest
  private tokenName = '';
  private accessToken = '';
  private selectedToken?: Token = undefined;
  private overlay = false;
  private loading = false;
  get tokens(): Token[] {
    return this.$store.state.user.tokens;
  }

  createTime(token: Token): string {
    const date = new Date(token.createdAt);
    return date.toLocaleString();
  }

  generateToken() {
    const base = this.$store.state.base;
    const formData = new FormData();
    formData.append('name', this.tokenName);
    this.loading = true;
    this.$http
      .post<string>(`${base}/api/v1/user/tokens`, formData)
      .then(response => {
        this.accessToken = response.data;
      })
      .finally(() => {
        this.tokenName = '';
        this.loading = false;
        this.$store.dispatch(Actions.FETCH_USER_TOKENS);
      });
  }

  clickedDelete(token: Token) {
    this.overlay = true;
    this.selectedToken = token;
  }

  cancelDelete() {
    this.selectedToken = undefined;
    this.overlay = false;
  }

  deleteSelectedToken() {
    if (!this.selectedToken) {
      return;
    }
    this.loading = true;
    const base = this.$store.state.base;
    const tokenID = this.selectedToken.id;
    this.$http.delete(`${base}/api/v1/user/tokens/${tokenID}`).finally(() => {
      this.loading = false;
      this.overlay = false;
      this.selectedToken = undefined;
      this.$store.dispatch(Actions.FETCH_USER_TOKENS);
    });
  }
}
</script>
