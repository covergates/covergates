import { mock } from 'jest-mock-extended';
import { ActionContext } from 'vuex';
import { UserState, Mutations } from '..';
import { fetchSCM } from '../actions';
import { RootState } from '@/store';
import axios, { AxiosResponse } from 'axios';

describe('store.module.user', () => {
  let context: ActionContext<UserState, RootState>;
  beforeEach(() => {
    context = mock<ActionContext<UserState, RootState>>();
    context.rootState.base = '';
  });

  it('commit user scm after fetch', async () => {
    const state = {
      gitea: true,
      github: false
    };
    const mockGet = jest.spyOn(axios, 'get');
    mockGet.mockResolvedValueOnce({
      status: 200,
      data: state
    } as AxiosResponse);
    await fetchSCM(context);
    expect(mockGet).toHaveBeenCalled();
    expect(context.commit).toHaveBeenCalledWith(Mutations.START_USER_LOADING);
    expect(context.commit).toHaveBeenCalledWith(Mutations.UPDATE_USER_SCM, state);
    expect(context.commit).toHaveBeenLastCalledWith(Mutations.STOP_USER_LOADING);
  });
});
