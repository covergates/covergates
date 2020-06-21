import '@mdi/font/css/materialdesignicons.css';
import Vue from 'vue';
import Vuetify from 'vuetify/lib';

import GiteaIcon from '@/components/icons/GiteaIcon.vue';

Vue.use(Vuetify);

export default new Vuetify({
  icons: {
    iconfont: 'mdi',
    values: {
      gitea: {
        component: GiteaIcon
      }
    }
  }
});
