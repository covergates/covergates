import { ActionContext } from 'vuex';
import { MutationTypes } from '../mutations';
import Axios from 'axios';
import { State } from '../';

export function fetchUser<S extends State, R>(context: ActionContext<S, R>) {
  Axios.get<User>(`${context.state.base}/api/v1/user`).then((response) => {
    context.commit(MutationTypes.UPDATE_USER, response.data);
  }).catch((error) => {
    if (error.response.status === 404) {
      context.commit(MutationTypes.UPDATE_USER, { error: 'user not found' });
    } else {
      context.commit(MutationTypes.UPDATE_USER, { error: error.response.data });
    }
  });
}
