import Vue from 'vue';
import App from './App.vue';
import router from './router';
import store from './store';
import vuetify from './plugins/vuetify';
import { AxiosPlugin } from '@/plugins/http';
import { ActionTypes } from './store/actions';
import { makeServer } from './server';

__webpack_public_path__ = process.env.NODE_ENV === 'production' ? `${VUE_BASE}/` : process.env.BASE_URL;

Vue.use(AxiosPlugin);

if (process.env.NODE_ENV === 'development') {
  makeServer();
}

Vue.config.productionTip = false;
store.dispatch(ActionTypes.FETCH_USER).then(() => {
  new Vue({
    router,
    store,
    vuetify,
    render: h => h(App)
  }).$mount('#app');
});
