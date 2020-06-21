import Vue from 'vue';
import { Framework } from 'vuetify';
declare module 'vue/types/vue' {
  interface Vue {
    $vuetify: Framework;
  }
}
