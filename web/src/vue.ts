import Vue, { VueConstructor } from 'vue';
import { Store } from 'vuex';
import { State } from '@/store';

abstract class VueClass extends Vue {
  public $store!: Store<State>;
}

const RootVue = Vue as VueConstructor<VueClass>;
export default RootVue;
