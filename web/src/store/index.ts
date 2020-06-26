import Vue from 'vue';
import Vuex from 'vuex';
import UserModule, {
  Actions as UserActions,
  Mutations as UserMutations
} from './modules/user';
import RepoModule, {
  Actions as RepoActions
} from './modules/repository';
Vue.use(Vuex);

function getBaseURL(): string {
  const base = process.env.NODE_ENV === 'production' ? VUE_BASE : process.env.BASE_URL;
  if (base === '/') {
    return '';
  }
  return base;
}

/**
 * Enum for Vux Mutations
 */
export const Mutations = {
  ...UserMutations
};

/**
 * Enum for Vuex Actions
 * @readonly
 * @enum
 */
export const Actions = {
  ...UserActions,
  ...RepoActions
};

/**
 * State of the the Vux, includes the state from modules.
 * Modules should declare it's own state to this interface.
 */
export interface State {
  base: string;
}

const rootState = {
  base: getBaseURL()
};

/**
 * RootState is the state that belongs the root store,
 * excludes state from modules.
 */
export type RootState = typeof rootState;
export default new Vuex.Store<RootState>({
  state: rootState,
  mutations: {},
  actions: {},
  modules: {
    user: UserModule,
    repository: RepoModule
  }
});
