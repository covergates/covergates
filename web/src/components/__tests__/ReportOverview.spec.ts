import Vue from 'vue';
import { createLocalVue, shallowMount, Wrapper } from '@vue/test-utils';
import Vuex, { Store } from 'vuex';
import { cloneDeep } from 'lodash';
import Vuetify from 'vuetify';
import flushPromises from 'flush-promises';
import ReportOverview from '@/components/ReportOverview.vue';
import { AxiosPlugin } from '@/plugins/http';
import { Mutations, storeConfig, RootState } from '@/store';

jest.mock('axios');
Vue.use(Vuetify);
Vue.use(AxiosPlugin);

describe('ReportOverview.vue', () => {
  console.warn = jest.fn();
  const localVue = createLocalVue();
  let vuetify: typeof Vuetify;
  let store: Store<RootState>;
  beforeEach(() => {
    vuetify = new Vuetify();
    store = new Vuex.Store(cloneDeep(storeConfig));
  });

  localVue.use(Vuetify);
  it('update file count when current repository is set', async () => {
    const wrapper = shallowMount(ReportOverview, {
      localVue,
      vuetify,
      store
    }) as Wrapper<ReportOverview & { filesCount: number }>;
    expect(wrapper.vm.$store.state.repository.current).toBeUndefined();
    expect(wrapper.vm.filesCount).toEqual(0);
    wrapper.vm.$store.commit(Mutations.SET_REPOSITORY_CURRENT, {
      Name: 'repo',
      NameSpace: 'org',
      SCM: 'github',
      Files: ['a', 'b', 'c']
    } as Repository);
    await flushPromises();
    expect(wrapper.vm.filesCount).toEqual(3);
  });
});
