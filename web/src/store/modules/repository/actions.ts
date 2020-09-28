import { ActionContext } from 'vuex';
import Axios from 'axios';
import { RepoState, Mutations, Actions } from '.';
import { RootState } from '@/store';
import { reasonToError } from '@/plugins/http';

const errUndefinedCurrentRepo = new Error('current repository is undefined');

function fetchRepositories(context: ActionContext<RepoState, RootState>): Promise<void> {
  return new Promise((resolve) => {
    Axios.get<Repository[]>(`${context.rootState.base}/api/v1/user/repos`)
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
    Axios.get<Repository>(
      `${base}/api/v1/repos/${scm}/${namespace}/${name}`
    ).then((response) => {
      repository = response.data;
    }).catch((reason) => {
      let error: Error;
      if (reason.response) {
        error = new Error(reason.response.data);
      } else if (reason.message) {
        error = new Error(reason.message);
      } else {
        error = new Error('Unknown Error');
      }
      if (!repository) {
        reject(error);
      }
    }).finally(() => {
      resolve(repository);
    });
  });
}

export function fetchRepositoryList<S extends RepoState, R extends RootState>(context: ActionContext<S, R>): Promise<void> {
  context.commit(Mutations.START_REPOSITORY_LOADING);
  return fetchRepositories(context).then(() => {
    context.commit(Mutations.STOP_REPOSITORY_LOADING);
  });
}

export function synchronizeRepository(context: ActionContext<RepoState, RootState>): Promise<void> {
  return new Promise(resolve => {
    Axios.patch(`${context.rootState.base}/api/v1/user/repos`)
      .then(() => {
        context.dispatch(Actions.FETCH_REPOSITORY_LIST);
      }).catch(reason => {
        console.warn(reasonToError(reason));
      }).finally(() => {
        resolve();
      });
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
    const { SCM, Name, NameSpace } = (repo as Repository);
    fetchRepository(context.rootState.base, SCM, NameSpace, Name)
      .then((response) => {
        context.commit(Mutations.SET_REPOSITORY_CURRENT, response);
      }).catch(error => {
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
    }).catch(reason => {
      const error = reasonToError(reason);
      context.commit(Mutations.SET_REPOSITORY_ERROR, error);
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
    fetchRepository(context.rootState.base, params.scm, params.namespace, params.name)
      .then((response) => {
        context.commit(Mutations.SET_REPOSITORY_CURRENT, response);
      }).catch(error => {
        context.commit(Mutations.SET_REPOSITORY_ERROR, error);
      }).finally(() => {
        context.commit(Mutations.STOP_REPOSITORY_LOADING);
        resolve();
      });
  });
}

export function fetchRepositorySetting<S extends RepoState, R extends RootState>(
  context: ActionContext<S, R>
): Promise<void> {
  return new Promise((resolve) => {
    const base = context.rootState.base;
    const repo = context.state.current;
    if (repo === undefined) {
      context.commit(Mutations.SET_REPOSITORY_SETTING, undefined);
      resolve();
      return;
    }
    context.commit(Mutations.START_REPOSITORY_LOADING);
    const { SCM, Name, NameSpace } = (repo as Repository);
    Axios.get<RepositorySetting>(`${base}/api/v1/repos/${SCM}/${NameSpace}/${Name}/setting`)
      .then(response => {
        context.commit(Mutations.SET_REPOSITORY_SETTING, response.data);
      })
      .catch(() => {
        context.commit(Mutations.SET_REPOSITORY_SETTING, undefined);
      }).finally(() => {
        context.commit(Mutations.STOP_REPOSITORY_LOADING);
        resolve();
      });
  });
}

export function fetchRepositoryOwner<S extends RepoState, R extends RootState>(
  context: ActionContext<S, R>
): Promise<void> {
  return new Promise((resolve) => {
    const base = context.rootState.base;
    const repo = context.state.current;
    if (repo === undefined) {
      context.commit(Mutations.SET_REPOSITORY_OWNER, false);
      resolve();
      return;
    }
    const { SCM, Name, NameSpace } = (repo as Repository);
    Axios.get<RepositorySetting>(`${base}/api/v1/user/owner/${SCM}/${NameSpace}/${Name}`)
      .then(() => {
        context.commit(Mutations.SET_REPOSITORY_OWNER, true);
      })
      .catch(() => {
        context.commit(Mutations.SET_REPOSITORY_OWNER, false);
      }).finally(() => {
        resolve();
      });
  });
}

export function fetchRepositoryCommits<S extends RepoState, R extends RootState>(
  context: ActionContext<S, R>, ref = ''
): Promise<void> {
  return new Promise(resolve => {
    const base = context.rootState.base;
    const repo = context.state.current;
    if (repo === undefined) {
      context.commit(Mutations.SET_REPOSITORY_COMMITS, []);
      resolve();
      return;
    }
    const { SCM, Name, NameSpace } = (repo as Repository);
    Axios.get<Commit[]>(
      `${base}/api/v1/repos/${SCM}/${NameSpace}/${Name}/commits`,
      { params: { ref: ref } })
      .then(response => {
        context.commit(Mutations.SET_REPOSITORY_COMMITS, response.data);
      })
      .catch(() => {
        context.commit(Mutations.SET_REPOSITORY_COMMITS, []);
      })
      .finally(() => {
        resolve();
      });
  });
}

export function fetchRepositoryBranches<S extends RepoState, R extends RootState>(
  context: ActionContext<S, R>
): Promise<void> {
  return new Promise(resolve => {
    const base = context.rootState.base;
    const repo = context.state.current;
    if (repo === undefined) {
      context.commit(Mutations.SET_REPOSITORY_BRANCHES, []);
      resolve();
      return;
    }
    const { SCM: scm, Name: name, NameSpace: nameSpace } = (repo as Repository);
    Axios.get<string[]>(`${base}/api/v1/repos/${scm}/${nameSpace}/${name}/branches`)
      .then(response => {
        context.commit(Mutations.SET_REPOSITORY_BRANCHES, response.data);
      })
      .catch(() => {
        context.commit(Mutations.SET_REPOSITORY_BRANCHES, []);
      })
      .finally(() => {
        resolve();
      });
  });
}
