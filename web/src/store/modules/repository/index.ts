import { Module } from 'vuex';
import {
  updateList,
  startLoading,
  stopLoading,
  setCurrent,
  setError,
  setSetting,
  setOwner,
  setCommits,
  setBranches
} from './mutations';
import {
  fetchRepositoryList,
  updateRepositoryCurrent,
  updateRepositoryReportID,
  changeCurrentRepository,
  fetchRepositorySetting,
  fetchRepositoryOwner,
  fetchRepositoryCommits,
  fetchRepositoryBranches,
  synchronizeRepository
} from './actions';
import { RootState } from '@/store';

export enum Mutations {
  UPDATE_REPOSITORY_LIST = 'UPDATE_REPOSITORY_LIST',
  START_REPOSITORY_LOADING = 'START_REPOSITORY_LOADING',
  STOP_REPOSITORY_LOADING = 'STOP_REPOSITORY_LOADING',
  SET_REPOSITORY_CURRENT = 'SET_REPOSITORY_CURRENT',
  SET_REPOSITORY_ERROR = 'SET_REPOSITORY_ERROR',
  SET_REPOSITORY_SETTING = 'SET_REPOSITORY_SETTING',
  SET_REPOSITORY_OWNER = 'SET_REPOSITORY_OWNER',
  SET_REPOSITORY_COMMITS = 'SET_REPOSITORY_COMMITS',
  SET_REPOSITORY_BRANCHES = 'SET_REPOSITORY_BRANCHES'
}

export enum Actions {
  FETCH_REPOSITORY_LIST = 'FETCH_REPOSITORY_LIST',
  SYNCHRONIZE_REPOSITORY = 'SYNCHRONIZE_REPOSITORY',
  UPDATE_REPOSITORY_CURRENT = 'UPDATE_REPOSITORY_CURRENT',
  UPDATE_REPOSITORY_REPORT_ID = 'UPDATE_REPOSITORY_REPORT_ID',
  CHANGE_CURRENT_REPOSITORY = 'CHANGE_CURRENT_REPOSITORY',
  FETCH_REPOSITORY_SETTING = 'FETCH_REPOSITORY_SETTING',
  FETCH_REPOSITORY_OWNER = 'FETCH_REPOSITORY_OWNER',
  FETCH_REPOSITORY_COMMITS = 'FETCH_REPOSITORY_COMMITS',
  FETCH_REPOSITORY_BRANCHES = 'FETCH_REPOSITORY_BRANCHES'
}

export type RepoState = {
  loading: boolean;
  current?: Repository;
  commits: Commit[];
  branches: string[];
  owner: boolean;
  setting?: RepositorySetting;
  list: Repository[];
  error?: Error;
};

const module: Module<RepoState, RootState> = {
  state: {
    loading: false,
    current: undefined,
    commits: [],
    branches: [],
    setting: undefined,
    list: [],
    error: undefined,
    owner: false
  },
  mutations: {
    [Mutations.UPDATE_REPOSITORY_LIST]: updateList,
    [Mutations.START_REPOSITORY_LOADING]: startLoading,
    [Mutations.STOP_REPOSITORY_LOADING]: stopLoading,
    [Mutations.SET_REPOSITORY_CURRENT]: setCurrent,
    [Mutations.SET_REPOSITORY_ERROR]: setError,
    [Mutations.SET_REPOSITORY_SETTING]: setSetting,
    [Mutations.SET_REPOSITORY_OWNER]: setOwner,
    [Mutations.SET_REPOSITORY_COMMITS]: setCommits,
    [Mutations.SET_REPOSITORY_BRANCHES]: setBranches
  },
  actions: {
    [Actions.FETCH_REPOSITORY_LIST]: fetchRepositoryList,
    [Actions.UPDATE_REPOSITORY_CURRENT]: updateRepositoryCurrent,
    [Actions.UPDATE_REPOSITORY_REPORT_ID]: updateRepositoryReportID,
    [Actions.CHANGE_CURRENT_REPOSITORY]: changeCurrentRepository,
    [Actions.FETCH_REPOSITORY_SETTING]: fetchRepositorySetting,
    [Actions.FETCH_REPOSITORY_OWNER]: fetchRepositoryOwner,
    [Actions.FETCH_REPOSITORY_COMMITS]: fetchRepositoryCommits,
    [Actions.FETCH_REPOSITORY_BRANCHES]: fetchRepositoryBranches,
    [Actions.SYNCHRONIZE_REPOSITORY]: synchronizeRepository
  }
};

declare module '@/store' {
  interface State {
    repository: RepoState;
  }
}

export default module;
