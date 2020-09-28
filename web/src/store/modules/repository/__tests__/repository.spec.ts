import { mock } from 'jest-mock-extended';
import { ActionContext } from 'vuex';
import axios, { AxiosResponse, AxiosError } from 'axios';
import { RepoState, Mutations, Actions } from '..';
import { fetchRepositoryList, changeCurrentRepository, fetchRepositorySetting, synchronizeRepository } from '@/store/modules/repository/actions';
import { makeServer, MockServer } from '@/server';
import { RootState } from '@/store';

describe('store.module.repository', () => {
  let server: MockServer;
  beforeEach(() => {
    server = makeServer();
    server.logging = false;
  });
  afterEach(() => {
    server.shutdown();
  });

  it('mark repository loading when fetch', async () => {
    const context = mock<ActionContext<RepoState, RootState>>();
    context.rootState.base = '';
    await fetchRepositoryList(context);
    expect(context.commit).toBeCalledTimes(3);
    const calls = context.commit.mock.calls;
    expect(calls[0][0]).toEqual(Mutations.START_REPOSITORY_LOADING);
    expect(calls[1][0]).toEqual(Mutations.UPDATE_REPOSITORY_LIST);
    expect(calls[calls.length - 1][0]).toEqual(Mutations.STOP_REPOSITORY_LOADING);
  });

  it('fetch repository from /user/repos', async () => {
    const context = mock<ActionContext<RepoState, RootState>>();
    context.rootState.base = '';
    jest.mock('axios');
    const spy = jest.spyOn(axios, 'get');
    spy.mockResolvedValueOnce({
      data: []
    });
    await fetchRepositoryList(context);
    expect(spy).toHaveBeenCalledWith('/api/v1/user/repos');
  });

  it('fetch repository after synchronize', async () => {
    const context = mock<ActionContext<RepoState, RootState>>();
    context.rootState.base = '';
    jest.mock('axios');
    const spy = jest.spyOn(axios, 'patch');
    spy.mockResolvedValueOnce({});
    await synchronizeRepository(context);
    expect(context.dispatch).toHaveBeenCalledWith(Actions.FETCH_REPOSITORY_LIST);
  });

  it('update current when change repository successfully', async () => {
    jest.mock('axios');
    const spy = jest.spyOn(axios, 'get');
    const repository = {
      SCM: 'github'
    } as Repository;
    spy.mockResolvedValueOnce({
      data: repository
    });
    const context = mock<ActionContext<RepoState, RootState>>();
    context.rootState.base = '';
    await changeCurrentRepository(context, {
      name: 'repo',
      namespace: 'repo',
      scm: 'github'
    });
    expect(spy).toBeCalled();
    expect(context.commit).toHaveBeenCalledWith(Mutations.SET_REPOSITORY_ERROR);
    expect(context.commit).toHaveBeenCalledWith(Mutations.SET_REPOSITORY_CURRENT, repository);
  });

  it('set error when change repository fail', async () => {
    jest.mock('axios');
    const spy = jest.spyOn(axios, 'get');
    const msg = 'error';
    spy.mockRejectedValueOnce({
      response: {
        data: msg
      }
    });
    const context = mock<ActionContext<RepoState, RootState>>();
    await changeCurrentRepository(context, {
      name: 'repo',
      namespace: 'repo',
      scm: 'github'
    });
    expect(context.commit).toHaveBeenCalledWith(Mutations.SET_REPOSITORY_ERROR, new Error(msg));
  });

  it('fetch repository setting', async () => {
    const spy = jest.spyOn(axios, 'get');
    const setting = {} as RepositorySetting;
    spy.mockResolvedValueOnce({
      data: setting
    } as AxiosResponse);
    const context = mock<ActionContext<RepoState, RootState>>();
    context.rootState.base = '';
    context.state.current = {
      Name: 'repo',
      NameSpace: 'org',
      SCM: 'github'
    } as Repository;
    await fetchRepositorySetting(context);
    expect(context.commit).toHaveBeenCalledWith(Mutations.START_REPOSITORY_LOADING);
    expect(context.commit).toHaveBeenCalledWith(Mutations.SET_REPOSITORY_SETTING, setting);
    expect(context.commit).toHaveBeenLastCalledWith(Mutations.STOP_REPOSITORY_LOADING);
    spy.mockRejectedValueOnce({
      response: {
        status: 404
      }
    } as AxiosError);
    spy.mockClear();
    context.commit.mockClear();
    expect(context.commit).not.toHaveBeenCalledWith(Mutations.START_REPOSITORY_LOADING);
    await fetchRepositorySetting(context);
    expect(context.commit).toHaveBeenCalledWith(Mutations.START_REPOSITORY_LOADING);
    expect(context.commit).toHaveBeenCalledWith(Mutations.SET_REPOSITORY_SETTING, undefined);
    expect(context.commit).toHaveBeenLastCalledWith(Mutations.STOP_REPOSITORY_LOADING);
  });

  it('set repository setting undefined if no current repository', async () => {
    const context = mock<ActionContext<RepoState, RootState>>();
    expect(context.state.current).toBeUndefined();
    await fetchRepositorySetting(context);
    expect(context.commit).not.toHaveBeenCalledWith(Mutations.START_REPOSITORY_LOADING);
    expect(context.commit).toHaveBeenCalledWith(Mutations.SET_REPOSITORY_SETTING, undefined);
  });
});
