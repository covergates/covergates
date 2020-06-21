import { State } from '@/store';

export function updateUser(state: State, user: User) {
  state.user = user;
}
