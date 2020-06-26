import { ActionContext } from 'vuex';
import { RootState } from '@/store';
import { RepoState, Mutations } from '.';
import Axios from 'axios';

const errUndefinedCurrentRepo = new Error('current repository is undefined');

function fetchSCMRepositories(context: ActionContext<RepoState, RootState>, scm: string): Promise<void> {
  return new Promise((resolve) => {
    Axios.get<Repository[]>(`${context.rootState.base}/api/v1/repos/${scm}`)
      .then((response) => {
        context.commit(Mutations.UPDATE_REPOSITORY_LIST, response.data);
      }).catch((error) => {
        console.warn(error);
      }).finally(() => {
        resolve();
      });
  });
}

export function fetchRepositoryList<S extends RepoState, R extends RootState>(context: ActionContext<S, R>): Promise<void> {
  context.commit(Mutations.START_REPOSITORY_LOADING);
  const providers = ['github'];
  const jobs: Promise<void>[] = [];
  for (const scm of providers) {
    jobs.push(fetchSCMRepositories(context, scm));
  }
  return Promise.all(jobs).then(() => {
    context.commit(Mutations.STOP_REPOSITORY_LOADING);
  });
}

export function updateRepositoryCurrent<S extends RepoState, R extends RootState>(context: ActionContext<S, R>): Promise<void> {
  return new Promise<void>((resolve, reject) => {
    context.commit(Mutations.START_REPOSITORY_LOADING);
    if (context.state.current === undefined) {
      context.commit(Mutations.STOP_REPOSITORY_LOADING);
      reject(errUndefinedCurrentRepo);
    }
    const repo = context.state.current;
    Axios.get<Repository>(
      `${context.rootState.base}/api/v1/${repo?.SCM}/${repo?.NameSpace}/${repo?.Name}`
    ).then((response) => {
      context.commit(Mutations.SET_REPOSITORY_CURRENT, response.data);
    }).catch((error) => {
      // TODO: add error mutation
      reject(error);
    }).finally(() => {
      context.commit(Mutations.STOP_REPOSITORY_LOADING);
      resolve();
    });
  });
}

export function updateRepositoryReportID<S extends RepoState, R extends RootState>(context: ActionContext<S, R>): Promise<void> {
  return new Promise<void>((resolve, reject) => {
    context.commit(Mutations.START_REPOSITORY_LOADING);
    const repo = context.state.current;
    if (repo === undefined) {
      context.commit(Mutations.STOP_REPOSITORY_LOADING);
      reject(errUndefinedCurrentRepo);
    }
    Axios.patch<Repository>(
      `${context.rootState.base}/api/v1/${repo?.SCM}/${repo?.NameSpace}/${repo?.Name}`
    ).then((response) => {
      context.commit(Mutations.SET_REPOSITORY_CURRENT, response.data);
    }).catch((error) => {
      // TODO: add error mutation
      reject(error);
    }).finally(() => {
      context.commit(Mutations.STOP_REPOSITORY_LOADING);
      resolve();
    });
  });
}
