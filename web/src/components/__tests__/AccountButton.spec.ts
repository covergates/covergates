import Vue from 'vue';
import Vuetify from 'vuetify';
import Vuex, { Store } from 'vuex';
import { createLocalVue, shallowMount, Wrapper } from '@vue/test-utils';
import { cloneDeep } from 'lodash';
import { storeConfig, RootState, Mutations } from '@/store';
import AccountButton from '@/components/AccountButton.vue';

Vue.use(Vuetify);

function findActions<T extends Vue>(wrapper: Wrapper<T>): Wrapper<Vue>[] {
  const actions = wrapper.findAll('v-list-item-icon-stub');
  const icons = [];
  for (const action of actions.wrappers) {
    const w = action.find('v-icon-stub');
    icons.push(w);
  }
  return icons;
}

describe('AccountButton.vue', () => {
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

  it('show login button when user undefined', () => {
    const wrapper = shallowMount(AccountButton, {
      vuetify, store, stubs: ['router-link']
    }) as Wrapper<AccountButton & { user?: User }>;
    expect(wrapper.vm.user).toBeUndefined();
    const actions = findActions(wrapper);
    const icons = actions.map(action => { return action.text(); });
    expect(actions.length).toBeGreaterThanOrEqual(1);
    expect(icons).toContain('mdi-login');
    expect(icons).not.toContain('mdi-logout');
  });

  it('show logout button when user defined', async () => {
    const $router = {
      push: jest.fn()
    };
    const wrapper = shallowMount(AccountButton, {
      vuetify,
      store,
      stubs: ['router-line'],
      mocks: {
        $router
      }
    }) as Wrapper<AccountButton & { user?: User }>;
    globalThis.window = Object.create(window);
    Object.defineProperty(window, 'location', {
      value: {
        href: ''
      }
    });
    expect(wrapper.vm.user).toBeUndefined();
    store.commit(Mutations.UPDATE_USER, {} as User);
    expect(wrapper.vm.user).toBeDefined();
    await wrapper.vm.$nextTick();
    const actions = findActions(wrapper);
    const icons = actions.map(action => { return action.text(); });
    expect(icons).toContain('mdi-logout');
    const index = Array.from(actions.keys()).find(i => { return icons[i] === 'mdi-logout'; });
    if (index !== undefined) {
      wrapper.findAll('v-list-item-stub').at(index).vm.$emit('click');
    }
    // expect($router.push).toHaveBeenCalledWith('/logoff');
    expect(window.location.href).toMatch('/logoff');
  });
});
