import { ActionContext } from 'vuex';
import Axios from 'axios';
import { ReportState, Mutations } from '.';
import { RootState } from '@/store';

export function fetchCurrentReport<S extends ReportState, R extends RootState>(context: ActionContext<S, R>, reportID: string): Promise<void> {
  return new Promise((resolve) => {
    context.commit(Mutations.START_REPORT_LOADING);
    Axios.get<Report>(`${context.rootState.base}/api/v1/reports/${reportID}`)
      .then((response) => {
        context.commit(Mutations.SET_REPORT_CURRENT, response.data);
      })
      .catch((error) => {
        // TODO: add error mutation
        console.warn(error);
      })
      .finally(() => {
        context.commit(Mutations.STOP_REPORT_LOADING);
        resolve();
      });
  });
}
