import { Module } from 'vuex';
import {
  fetchCurrentReport,
  fetchSource,
  fetchHistory
} from './actions';
import {
  setCurrent,
  startLoading,
  stopLoading,
  setSource,
  setError,
  setHistory
} from './mutations';
import { RootState } from '@/store';

export enum Mutations {
  SET_REPORT_CURRENT = 'SET_REPORT_CURRENT',
  SET_REPORT_HISTORY = 'SET_REPORT_HISTORY',
  START_REPORT_LOADING = 'START_REPORT_LOADING',
  STOP_REPORT_LOADING = 'STOP_REPORT_LOADING',
  SET_REPORT_SOURCE = 'SET_REPORT_SOURCE',
  SET_REPORT_ERROR = 'SET_REPORT_ERROR'
}

export enum Actions {
  FETCH_REPORT_CURRENT = 'FETCH_REPORT_CURRENT',
  FETCH_REPORT_HISTORY = 'FETCH_REPORT_HISTORY',
  FETCH_REPORT_SOURCE = 'FETCH_REPORT_SOURCE'
}

export type ReportState = {
  current?: Report;
  history: Report[];
  loading: boolean;
  source?: string;
  error?: Error;
};

const module: Module<ReportState, RootState> = {
  state: {
    loading: false,
    current: undefined,
    history: [],
    source: undefined,
    error: undefined
  },
  actions: {
    [Actions.FETCH_REPORT_CURRENT]: fetchCurrentReport,
    [Actions.FETCH_REPORT_HISTORY]: fetchHistory,
    [Actions.FETCH_REPORT_SOURCE]: fetchSource
  },
  mutations: {
    [Mutations.START_REPORT_LOADING]: startLoading,
    [Mutations.STOP_REPORT_LOADING]: stopLoading,
    [Mutations.SET_REPORT_CURRENT]: setCurrent,
    [Mutations.SET_REPORT_SOURCE]: setSource,
    [Mutations.SET_REPORT_ERROR]: setError,
    [Mutations.SET_REPORT_HISTORY]: setHistory
  }
};

declare module '@/store' {
  interface State {
    report: ReportState;
  }
}

export default module;
