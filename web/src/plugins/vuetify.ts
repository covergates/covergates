import '@mdi/font/css/materialdesignicons.css';
import Vue from 'vue';
import Vuetify from 'vuetify/lib';
import colors from 'vuetify/lib/util/colors';

import GiteaIcon from '@/components/icons/GiteaIcon.vue';

Vue.use(Vuetify);

const themes = {
  light: {
    primary: colors.cyan.darken3,
    secondary: colors.grey.darken4,
    accent: colors.blueGrey.darken1,
    error: colors.red.darken4
  }
};

export default new Vuetify({
  theme: {
    themes: themes
  },
  icons: {
    iconfont: 'mdi',
    values: {
      gitea: {
        component: GiteaIcon
      }
    }
  }
});
