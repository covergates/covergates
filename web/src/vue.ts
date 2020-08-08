import Vue, { VueConstructor } from 'vue';
import { Store } from 'vuex';
import { State } from '@/store';

export abstract class RootVueClass extends Vue {
  public $store!: Store<State>;
}

const RootVue = Vue as VueConstructor<RootVueClass>;
export default RootVue;
