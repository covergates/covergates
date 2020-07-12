import { mock } from 'jest-mock-extended';
import { ActionContext } from 'vuex';
import { ReportState, Mutations } from '..';
import { RootState } from '@/store';
import { fetchCurrentReport, fetchSource } from '../actions';
import axios, { AxiosResponse, AxiosError } from 'axios';
import { Route } from 'vue-router';

describe('store.module.report.actions', () => {
  let context: ActionContext<ReportState, RootState>;

  beforeEach(() => {
    context = mock<ActionContext<ReportState, RootState>>();
    context.rootState.base = '';
  });

  it('commit current report when fetch success', () => {
    const mockGet = jest.spyOn(axios, 'get');
    const report = {};
    mockGet.mockResolvedValueOnce({
      status: 200,
      data: report
    } as AxiosResponse);
    return fetchCurrentReport(context, '1234').then(() => {
      expect(context.commit).toHaveBeenCalledWith(Mutations.SET_REPORT_CURRENT, report);
    });
  });

  it('commit error message when report not found', () => {
    const mockGet = jest.spyOn(axios, 'get');
    mockGet.mockRejectedValueOnce({
      response: {
        status: 404
      }
    } as AxiosError);
    return fetchCurrentReport(context, '1234').then(() => {
      expect(context.commit).toHaveBeenNthCalledWith(1, Mutations.START_REPORT_LOADING);
      expect(context.commit).toHaveBeenCalledWith(
        Mutations.SET_REPORT_ERROR, expect.any(Error));
      expect(context.commit).toHaveBeenCalledWith(
        Mutations.SET_REPORT_CURRENT, undefined);
      expect(context.commit).toHaveBeenLastCalledWith(Mutations.STOP_REPORT_LOADING);
    });
  });

  it('commit error message when fetch source fail', () => {
    const mockGet = jest.spyOn(axios, 'get');
    mockGet.mockRejectedValueOnce({
      response: {
        status: 404
      }
    } as AxiosError);
    const route = mock<Route>();
    route.params = {
      scm: '',
      namespace: '',
      name: '',
      path: ''
    };
    return fetchSource(context, route).then(() => {
      expect(context.commit).toHaveBeenNthCalledWith(1, Mutations.START_REPORT_LOADING);
      expect(context.commit).toHaveBeenCalledWith(
        Mutations.SET_REPORT_ERROR, expect.any(Error));
      expect(context.commit).toHaveBeenLastCalledWith(Mutations.STOP_REPORT_LOADING);
    });
  });
});
