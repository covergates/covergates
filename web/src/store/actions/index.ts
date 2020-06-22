import { ActionTree } from 'vuex';
import { fetchUser } from './user';
import { State } from '../';

export enum ActionTypes {
  FETCH_USER = 'FETCH_USER'
}

export function actions<S extends State>(): ActionTree<S, S> {
  return {
    [ActionTypes.FETCH_USER]: fetchUser
  };
}
