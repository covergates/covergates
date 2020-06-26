import { ActionContext } from 'vuex';
import Axios from 'axios';
import { UserState, Mutations } from '.';
import { RootState } from '@/store';

export function fetchUser<S extends UserState, R extends RootState>(context: ActionContext<S, R>) {
  return new Promise<void>((resolve) => {
    Axios.get<User>(`${context.rootState.base}/api/v1/user`).then((response) => {
      context.commit(Mutations.UPDATE_USER, response.data);
    }).catch((error) => {
      if (error.response.status === 404) {
        context.commit(Mutations.UPDATE_USER, { error: 'user not found' });
      } else {
        context.commit(Mutations.UPDATE_USER, { error: error.response.data });
      }
    }).finally(() => {
      resolve();
    });
  });
}
