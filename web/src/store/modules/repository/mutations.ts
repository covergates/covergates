import { RepoState } from '.';

export function startLoading(state: RepoState): void {
  state.loading = true;
  state.list.splice(0, state.list.length);
}

export function stopLoading(state: RepoState): void {
  state.loading = false;
}

export function updateList(state: RepoState, repos: Repository[]): void {
  state.list.push(...repos);
}

export function setCurrent(state: RepoState, repo: Repository): void {
  state.current = repo;
}
