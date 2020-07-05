import 'vuetify/dist/vuetify.min.css';
import Vue from 'vue';
import Vuex from 'vuex';
import Vuetify from '../src/plugins/vuetify';
import { AxiosPlugin } from '@/plugins/http';
Vue.use(Vuetify);
Vue.use(Vuex);
Vue.use(AxiosPlugin);
