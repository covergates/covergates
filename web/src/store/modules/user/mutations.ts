
import { UserState } from '.';

export function UpdateUser(state: UserState, user: User): void {
  state.current = user;
}

export function updateSCM(state: UserState, scm: Record<string, boolean>): void {
  state.scm = scm;
}

export function startLoading(state: UserState): void {
  state.loading = true;
}

export function stopLoading(state: UserState): void {
  state.loading = false;
}

export function SetError(state: UserState, err: Error): void {
  state.error = err;
}

export function ClearError(state: UserState): void {
  state.error = undefined;
}
