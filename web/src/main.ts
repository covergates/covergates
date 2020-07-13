import Vue from 'vue';
import App from './App.vue';
import router from './router';
import vuetify from './plugins/vuetify';
import { authorize } from './router/gates';
import { makeServer } from './server';
import store, { Actions } from '@/store';
import { AxiosPlugin } from '@/plugins/http';
import { HighlightPlugin } from '@/plugins/highlight';
import '@/plugins/scrollbar';

__webpack_public_path__ = process.env.NODE_ENV === 'production' ? `${VUE_BASE}/` : process.env.BASE_URL;

Vue.use(AxiosPlugin);
Vue.use(HighlightPlugin);

if (process.env.NODE_ENV === 'development' && process.env.VUE_APP_MOCK_SERVER === 'true') {
  makeServer();
}

Vue.config.productionTip = false;
router.beforeEach(authorize(store));
store.dispatch(Actions.FETCH_USER).then(() => {
  new Vue({
    router,
    store,
    vuetify,
    render: h => h(App)
  }).$mount('#app');
});

export default store;
