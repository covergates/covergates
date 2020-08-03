import Vue from 'vue';
import { createLocalVue, shallowMount } from '@vue/test-utils';
import Vuetify from 'vuetify';
import flushPromises from 'flush-promises';
import RepoListItem from '@/components/RepoListItem.vue';
import '@testing-library/jest-dom';
import axios, { AxiosError } from 'axios';
import Vuex, { Store } from 'vuex';
import { RootState, storeConfig } from '@/store';
import { cloneDeep } from 'lodash';
import { AxiosPlugin } from '@/plugins/http';

Vue.use(Vuetify);

describe('RepoListItem.vue', () => {
  console.warn = jest.fn();
  const localVue = createLocalVue();
  let vuetify: typeof Vuetify;
  let store: Store<RootState>;
  beforeEach(() => {
    vuetify = new Vuetify();
    store = new Vuex.Store(cloneDeep(storeConfig));
  });
  localVue.use(Vuetify);
  localVue.use(Vuex);
  localVue.use(AxiosPlugin);

  it('render source repository with unknown SCM', () => {
    const wrapper = shallowMount(RepoListItem, {
      localVue,
      vuetify,
      propsData: {
        repo: {
          SCM: 'unknown'
        }
      },
      stubs: ['router-link']
    });
    const icons = wrapper.findAll('v-icon-stub');
    expect(icons.length).toBeGreaterThan(0);
    expect(icons.at(0).text()).toBe('mdi-source-repository');
  });

  it('route to root if repository undefined', async () => {
    const wrapper = shallowMount(RepoListItem, {
      localVue,
      vuetify,
      propsData: {
        repo: {
          ReportID: 'abc'
        } as Repository
      }
    });
    await wrapper.vm.$nextTick();
    const button = wrapper.findComponent({
      ref: 'goto'
    });
    expect(button.attributes('to')).toBe('/');
  });

  it('route to root if repository has no name', async () => {
    const wrapper = shallowMount(RepoListItem, {
      localVue,
      vuetify,
      propsData: {
        repo: {
          SCM: 'gitea',
          ReportID: 'abc'
        } as Repository
      }
    });
    await wrapper.vm.$nextTick();
    const button = wrapper.findComponent({
      ref: 'goto'
    });
    expect(button.attributes('to')).toBe('/');
  });

  it('show activate button according to repository state', async () => {
    const wrapper = shallowMount(RepoListItem, {
      localVue,
      vuetify
    });
    let activateBtn = wrapper.findComponent({
      ref: 'activate'
    });
    expect(activateBtn.element).toBeVisible();
    await wrapper.setProps({
      repo: {
        ReportID: '1234'
      } as Repository
    });
    await wrapper.vm.$nextTick();
    activateBtn = wrapper.findComponent({
      ref: 'activate'
    });
    expect(activateBtn.element).not.toBeDefined();
  });

  it('handle repository activation fails', async () => {
    const wrapper = shallowMount(RepoListItem, {
      localVue,
      vuetify,
      store,
      propsData: {
        repo: {
          SCM: 'github',
          NameSpace: 'org',
          Name: 'repo'
        } as Repository
      }
    });
    const activateBtn = wrapper.findComponent({ ref: 'activate' });
    const snackBar = wrapper.find('v-snackbar-stub');
    const mockGet = jest.spyOn(axios, 'get');
    const mockPost = jest.spyOn(axios, 'post');
    mockGet.mockImplementationOnce(() => {
      expect(wrapper.vm.$data.loading).toBeTruthy();
      const err: Error = {
        response: {
          status: 404
        }
      } as AxiosError;
      return Promise.reject(err);
    });
    mockPost.mockRejectedValueOnce({
      response: {
        status: 500
      }
    } as AxiosError);
    expect(activateBtn.element).toBeVisible();
    activateBtn.vm.$emit('click');
    await flushPromises();
    expect(mockGet).toHaveBeenCalled();
    expect(mockPost).toHaveBeenCalled();
    expect(activateBtn.element).toBeVisible();
    expect(wrapper.vm.$data.loading).toBeFalsy();
    expect(snackBar.element).toHaveTextContent(/500/);
  });

  it('show route link button when activate successfully', async () => {
    const wrapper = shallowMount(RepoListItem, {
      localVue,
      vuetify,
      store,
      propsData: {
        repo: {
          SCM: 'github',
          NameSpace: 'org',
          Name: 'repo'
        } as Repository
      }
    });
    await wrapper.vm.$nextTick();
    const mockGet = jest.spyOn(axios, 'get');
    const mockPatch = jest.spyOn(axios, 'patch');
    mockGet.mockResolvedValueOnce({});
    mockPatch.mockResolvedValueOnce({});

    let activateBtn = wrapper.findComponent({ ref: 'activate' });
    let routeBtn = wrapper.findComponent({ ref: 'goto' });
    expect(routeBtn.element).not.toBeDefined();
    activateBtn.vm.$emit('click');
    return flushPromises().then(() => {
      routeBtn = wrapper.findComponent({ ref: 'goto' });
      activateBtn = wrapper.findComponent({ ref: 'activate' });
      expect(mockGet).toHaveBeenCalled();
      expect(mockPatch).toHaveBeenCalled();
      expect(routeBtn.element).toBeVisible();
      expect(activateBtn.element).not.toBeDefined();
    });
  });
});
