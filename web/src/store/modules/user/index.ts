import { Module } from 'vuex';
import { fetchUser } from './actions';
import { UpdateUser } from './mutations';
import { RootState } from '@/store';

export enum Actions {
  FETCH_USER = 'FETCH_USER'
}

export enum Mutations {
  UPDATE_USER = 'UPDATE_USER'
}

const state = {
  current: ({} as User)
};
export type UserState = typeof state;
const module: Module<UserState, RootState> = {
  state: state,
  actions: {
    [Actions.FETCH_USER]: fetchUser
  },
  mutations: {
    [Mutations.UPDATE_USER]: UpdateUser
  }
};

declare module '@/store' {
  interface State {
    user: UserState;
  }
}

export default module;
