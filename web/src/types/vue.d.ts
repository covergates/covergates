import { Framework } from 'vuetify';
import { AxiosStatic } from 'axios';
declare module 'vue/types/vue' {
  interface Vue {
    readonly $http: AxiosStatic;
    $vuetify: Framework;
  }
}
