import { makeServer, MockServer } from '@/server';
import { fetchRepositoryList } from '@/store/modules/repository/actions';
import { mock } from 'jest-mock-extended';
import { ActionContext } from 'vuex';
import { RepoState, Mutations } from '..';
import { RootState } from '@/store';

describe('store.module.repository', () => {
  let server: MockServer;
  beforeEach(() => {
    server = makeServer();
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
});
