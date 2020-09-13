import { ActionContext } from 'vuex';
import Axios from 'axios';
import { UserState, Mutations } from '.';
import { RootState } from '@/store';
import { reasonToError } from '@/plugins/http';

export function fetchUser<S extends UserState, R extends RootState>(context: ActionContext<S, R>) {
  return new Promise<void>((resolve) => {
    context.commit(Mutations.CLEAR_USER_ERROR);
    Axios.get<User>(`${context.rootState.base}/api/v1/user`).then((response) => {
      context.commit(Mutations.UPDATE_USER, response.data);
    }).catch(reason => {
      if (reason.response && reason.response.status === 404) {
        context.commit(Mutations.SET_USER_ERROR, new Error('user not found'));
      } else {
        context.commit(Mutations.SET_USER_ERROR, reasonToError(reason));
      }
    }).finally(() => {
      resolve();
    });
  });
}

export function fetchTokens<S extends UserState, R extends RootState>(context: ActionContext<S, R>): Promise<void> {
  return new Promise<void>((resolve) => {
    const base = context.rootState.base;
    context.commit(Mutations.START_USER_LOADING);
    context.commit(Mutations.UPDATE_USER_TOKENS, []);
    Axios.get<Token[]>(`${base}/api/v1/user/tokens`)
      .then(response => {
        context.commit(Mutations.UPDATE_USER_TOKENS, response.data);
      })
      .catch(reason => {
        console.warn(reasonToError(reason));
      })
      .finally(() => {
        context.commit(Mutations.STOP_USER_LOADING);
        resolve();
      });
  });
}

export function fetchSCM<S extends UserState, R extends RootState>(context: ActionContext<S, R>): Promise<void> {
  return new Promise<void>((resolve) => {
    const base = context.rootState.base;
    context.commit(Mutations.START_USER_LOADING);
    Axios.get<Record<string, boolean>>(`${base}/api/v1/user/scm`).then(response => {
      context.commit(Mutations.UPDATE_USER_SCM, response.data);
    }).catch(reason => {
      console.warn(reasonToError(reason));
    }).finally(() => {
      context.commit(Mutations.STOP_USER_LOADING);
      resolve();
    });
  });
}
