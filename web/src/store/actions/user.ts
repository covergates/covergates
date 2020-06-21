import { ActionContext } from 'vuex';
import { MutationTypes } from '../mutations';
import Axios from 'axios';

export function fetchUser<S, R>(context: ActionContext<S, R>) {
  Axios.get<User>('').then((response) => {
    context.commit(MutationTypes.UPDATE_USER, response.data);
  }).catch((error) => {
    if (error.response.status === 404) {
      context.commit(MutationTypes.UPDATE_USER, { error: 'user not found' });
    } else {
      context.commit(MutationTypes.UPDATE_USER, { error: error.response.data });
    }
  });
}
