import { Module } from 'vuex';
import {
  fetchCurrentReport,
  fetchSource
} from './actions';
import {
  setCurrent,
  startLoading,
  stopLoading,
  setSource
} from './mutations';
import { RootState } from '@/store';

export enum Mutations {
  SET_REPORT_CURRENT = 'SET_REPORT_CURRENT',
  START_REPORT_LOADING = 'START_REPORT_LOADING',
  STOP_REPORT_LOADING = 'STOP_REPORT_LOADING',
  SET_REPORT_SOURCE = 'SET_REPORT_SOURCE'
}

export enum Actions {
  FETCH_REPORT_CURRENT = 'FETCH_REPORT_CURRENT',
  FETCH_REPORT_SOURCE = 'FETCH_REPORT_SOURCE'
}

export type ReportState = {
  current?: Report;
  loading: boolean;
  source?: string;
};

const module: Module<ReportState, RootState> = {
  state: {
    loading: false,
    current: undefined,
    source: undefined
  },
  actions: {
    [Actions.FETCH_REPORT_CURRENT]: fetchCurrentReport,
    [Actions.FETCH_REPORT_SOURCE]: fetchSource
  },
  mutations: {
    [Mutations.START_REPORT_LOADING]: startLoading,
    [Mutations.STOP_REPORT_LOADING]: stopLoading,
    [Mutations.SET_REPORT_CURRENT]: setCurrent,
    [Mutations.SET_REPORT_SOURCE]: setSource
  }
};

declare module '@/store' {
  interface State {
    report: ReportState;
  }
}

export default module;
