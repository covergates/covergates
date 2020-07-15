import { Module } from 'vuex';
import { fetchUser } from './actions';
import {
  UpdateUser,
  ClearError,
  SetError
} from './mutations';
import { RootState } from '@/store';

export enum Actions {
  FETCH_USER = 'FETCH_USER'
}

export enum Mutations {
  UPDATE_USER = 'UPDATE_USER',
  SET_USER_ERROR = 'SET_USER_ERROR',
  CLEAR_USER_ERROR = 'CLEAR_USER_ERROR'
}

export type UserState = {
  current?: User;
  error?: Error;
};
const module: Module<UserState, RootState> = {
  state: {
    current: undefined,
    error: undefined
  },
  actions: {
    [Actions.FETCH_USER]: fetchUser
  },
  mutations: {
    [Mutations.UPDATE_USER]: UpdateUser,
    [Mutations.SET_USER_ERROR]: SetError,
    [Mutations.CLEAR_USER_ERROR]: ClearError
  }
};

declare module '@/store' {
  interface State {
    user: UserState;
  }
}

export default module;
