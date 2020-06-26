import { Module } from 'vuex';
import { RootState } from '@/store';
import { fetchCurrentReport } from './actions';
import {
  setCurrent,
  startLoading,
  stopLoading
} from './mutations';

export enum Mutations {
  SET_REPORT_CURRENT = 'SET_REPORT_CURRENT',
  START_REPORT_LOADING = 'START_REPORT_LOADING',
  STOP_REPORT_LOADING = 'STOP_REPORT_LOADING'
}

export enum Actions {
  FETCH_REPORT_CURRENT = 'FETCH_REPORT_CURRENT'
}

export type ReportState = {
  current?: Report;
  loading: boolean;
};

const module: Module<ReportState, RootState> = {
  state: {
    loading: false
  },
  actions: {
    [Actions.FETCH_REPORT_CURRENT]: fetchCurrentReport
  },
  mutations: {
    [Mutations.START_REPORT_LOADING]: startLoading,
    [Mutations.STOP_REPORT_LOADING]: stopLoading,
    [Mutations.SET_REPORT_CURRENT]: setCurrent
  }
};

export default module;
