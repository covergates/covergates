import { Module } from 'vuex';
import {
  updateList,
  startLoading,
  stopLoading,
  setCurrent,
  setError
} from './mutations';
import {
  fetchRepositoryList,
  updateRepositoryCurrent,
  updateRepositoryReportID,
  changeCurrentRepository
} from './actions';
import { RootState } from '@/store';

export enum Mutations {
  UPDATE_REPOSITORY_LIST = 'UPDATE_REPOSITORY_LIST',
  START_REPOSITORY_LOADING = 'START_REPOSITORY_LOADING',
  STOP_REPOSITORY_LOADING = 'STOP_REPOSITORY_LOADING',
  SET_REPOSITORY_CURRENT = 'SET_REPOSITORY_CURRENT',
  SET_REPOSITORY_ERROR = 'SET_REPOSITORY_ERROR'
}

export enum Actions {
  FETCH_REPOSITORY_LIST = 'FETCH_REPOSITORY_LIST',
  UPDATE_REPOSITORY_CURRENT = 'UPDATE_REPOSITORY_CURRENT',
  UPDATE_REPOSITORY_REPORT_ID = 'UPDATE_REPOSITORY_REPORT_ID',
  CHANGE_CURRENT_REPOSITORY = 'CHANGE_CURRENT_REPOSITORY'
}

export type RepoState = {
  loading: boolean;
  current?: Repository;
  list: Repository[];
  error?: Error;
};

const module: Module<RepoState, RootState> = {
  state: {
    loading: false,
    current: undefined,
    list: [],
    error: undefined
  },
  mutations: {
    [Mutations.UPDATE_REPOSITORY_LIST]: updateList,
    [Mutations.START_REPOSITORY_LOADING]: startLoading,
    [Mutations.STOP_REPOSITORY_LOADING]: stopLoading,
    [Mutations.SET_REPOSITORY_CURRENT]: setCurrent,
    [Mutations.SET_REPOSITORY_ERROR]: setError
  },
  actions: {
    [Actions.FETCH_REPOSITORY_LIST]: fetchRepositoryList,
    [Actions.UPDATE_REPOSITORY_CURRENT]: updateRepositoryCurrent,
    [Actions.UPDATE_REPOSITORY_REPORT_ID]: updateRepositoryReportID,
    [Actions.CHANGE_CURRENT_REPOSITORY]: changeCurrentRepository
  }
};

declare module '@/store' {
  interface State {
    repository: RepoState;
  }
}

export default module;
