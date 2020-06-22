import { MutationTree } from 'vuex';
import { updateUser } from './user';
import { State } from '..';

export enum MutationTypes {
  UPDATE_USER = 'UPDATE_USER'
}

export function mutations(): MutationTree<State> {
  return {
    [MutationTypes.UPDATE_USER]: updateUser
  };
}
