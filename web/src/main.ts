import Vue from 'vue';
import App from './App.vue';
import router from './router';
import store from './store';
// External Libs
import axios from 'axios';
import vuetify from './plugins/vuetify';

Vue.config.productionTip = false;
Vue.prototype.$http = axios;
new Vue({
  router,
  store,
  vuetify,
  render: h => h(App)
}).$mount('#app');
