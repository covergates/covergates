import { ActionContext } from 'vuex';
import Axios from 'axios';
import { RepoState, Mutations } from '.';
import { RootState } from '@/store';

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

function fetchRepository(base: string, scm: string, namespace: string, name: string): Promise<Repository> {
  return new Promise((resolve, reject) => {
    let repository: Repository;
    let error: Error;
    Axios.get<Repository>(
      `${base}/api/v1/${scm}/${namespace}/${name}`
    ).then((response) => {
      repository = response.data;
      return Axios.get<string[]>(`${base}/api/v1/${scm}/${namespace}/${name}/files`);
    }).then((response) => {
      repository.Files = response.data;
    }).catch((error) => {
      console.warn(error);
      if (error.response) {
        error = new Error(error.response.data);
      } else if (error.message) {
        error = new Error(error.message);
      } else {
        error = new Error('Unknown Error');
      }
    }).finally(() => {
      if (repository) {
        resolve(repository);
      } else {
        reject(error);
      }
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

export function updateRepositoryCurrent<S extends RepoState, R extends RootState>(
  context: ActionContext<S, R>
): Promise<void> {
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

export function updateRepositoryReportID<S extends RepoState, R extends RootState>(
  context: ActionContext<S, R>
): Promise<void> {
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

export function changeCurrentRepository<S extends RepoState, R extends RootState>(
  context: ActionContext<S, R>, params: { scm: string; namespace: string; name: string }
): Promise<void> {
  return new Promise((resolve) => {
    context.commit(Mutations.START_REPOSITORY_LOADING);
    context.commit(Mutations.SET_REPOSITORY_ERROR);
    Axios.get(`${context.rootState.base}/api/v1/repos/${params.scm}/${params.namespace}/${params.name}`)
      .then(response => {
        context.commit(Mutations.SET_REPOSITORY_CURRENT, response.data);
      })
      .catch((error) => {
        if (error.response) {
          context.commit(Mutations.SET_REPOSITORY_ERROR, new Error(error.response.data));
        } else if (error.message) {
          context.commit(Mutations.SET_REPOSITORY_ERROR, new Error(error.message));
        } else {
          context.commit(Mutations.SET_REPOSITORY_ERROR, new Error('Unknown error'));
        }
      })
      .finally(() => {
        context.commit(Mutations.STOP_REPOSITORY_LOADING);
        resolve();
      });
  });
}
