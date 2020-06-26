
import { UserState } from '.';

export function UpdateUser(state: UserState, user: User): void {
  state.current = user;
}
