import _Vue from 'vue';
import Axios from 'axios';

// eslint-disable-next-line
export function AxiosPlugin(Vue: typeof _Vue, options?: any): void {
  Vue.prototype.$http = Axios;
}
