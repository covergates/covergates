
import { UserState } from '.';

export function UpdateUser(state: UserState, user: User): void {
  state.current = user;
}

export function SetError(state: UserState, err: Error): void {
  state.error = err;
}

export function ClearError(state: UserState): void {
  state.error = undefined;
}
