import Vue from 'vue';
import { createLocalVue, shallowMount } from '@vue/test-utils';
import RepoListItem from '@/components/RepoListItem.vue';
import Vuetify from 'vuetify';

Vue.use(Vuetify);

describe('RepoListItem.vue', () => {
  const localVue = createLocalVue();
  let vuetify: typeof Vuetify;
  beforeEach(() => {
    vuetify = new Vuetify();
  });
  localVue.use(Vuetify);
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
  it('route to root if repository undefined', () => {
    const wrapper = shallowMount(RepoListItem, {
      localVue,
      vuetify
    });
    const button = wrapper.find('v-btn-stub');
    expect(button.attributes('to')).toBe('/');
  });
  it('route to root if repository has to name', () => {
    const wrapper = shallowMount(RepoListItem, {
      localVue,
      vuetify,
      propsData: {
        repo: {
          SCM: 'gitea'
        }
      }
    });
    const button = wrapper.find('v-btn-stub');
    expect(button.attributes('to')).toBe('/');
  });
});
