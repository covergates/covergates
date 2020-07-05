import { mock } from 'jest-mock-extended';
import { ActionContext } from 'vuex';
import axios from 'axios';
import { RepoState, Mutations } from '..';
import { fetchRepositoryList, changeCurrentRepository } from '@/store/modules/repository/actions';
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
});
