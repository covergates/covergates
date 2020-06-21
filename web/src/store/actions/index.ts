import { ActionTree } from 'vuex';
import { fetchUser } from './user';

export enum ActionTypes {
  FETCH_USER = 'FETCH_USER'
}

export function actions<S>(): ActionTree<S, S> {
  return {
    [ActionTypes.FETCH_USER]: fetchUser
  };
}
