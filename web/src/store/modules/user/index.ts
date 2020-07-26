import { Module } from 'vuex';
import {
  fetchUser,
  fetchSCM
} from './actions';
import {
  UpdateUser,
  ClearError,
  SetError,
  updateSCM,
  startLoading,
  stopLoading
} from './mutations';
import { RootState } from '@/store';

export enum Actions {
  FETCH_USER = 'FETCH_USER',
  FETCH_USER_SCM = 'FETCH_USER_SCM'
}

export enum Mutations {
  UPDATE_USER = 'UPDATE_USER',
  UPDATE_USER_SCM = 'UPDATE_USER_SCM',
  SET_USER_ERROR = 'SET_USER_ERROR',
  CLEAR_USER_ERROR = 'CLEAR_USER_ERROR',
  START_USER_LOADING = 'START_USER_LOADING',
  STOP_USER_LOADING = 'STOP_USER_LOADING'
}

export type UserState = {
  current?: User;
  error?: Error;
  scm?: Record<string, boolean>;
  loading: boolean;
};
const module: Module<UserState, RootState> = {
  state: {
    current: undefined,
    error: undefined,
    scm: undefined,
    loading: false
  },
  actions: {
    [Actions.FETCH_USER]: fetchUser,
    [Actions.FETCH_USER_SCM]: fetchSCM
  },
  mutations: {
    [Mutations.UPDATE_USER]: UpdateUser,
    [Mutations.SET_USER_ERROR]: SetError,
    [Mutations.CLEAR_USER_ERROR]: ClearError,
    [Mutations.START_USER_LOADING]: startLoading,
    [Mutations.STOP_USER_LOADING]: stopLoading,
    [Mutations.UPDATE_USER_SCM]: updateSCM
  }
};

declare module '@/store' {
  interface State {
    user: UserState;
  }
}

export default module;
