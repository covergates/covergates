import Vue from 'vue';
import Vuex from 'vuex';
import * as actions from './actions';
import { mutations } from './mutations';

Vue.use(Vuex);

function getBaseURL(): string {
  const base = process.env.NODE_ENV === 'production' ? VUE_BASE : process.env.BASE_URL;
  if (base === '/') {
    return '';
  }
  return base;
}

const user: User = {};

export const state = {
  base: getBaseURL(),
  user: user
};

export type State = typeof state

export default new Vuex.Store<State>({
  state: state,
  mutations: mutations(),
  actions: actions.actions<State>(),
  modules: {
  }
});
