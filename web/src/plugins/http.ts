import _Vue from 'vue';
import Axios, { AxiosStatic, AxiosError } from 'axios';

export function reasonToError(reason: Error | AxiosError): Error {
  if ((reason as AxiosError).response) {
    const response = (reason as AxiosError).response;
    if (response) {
      return response.data
        ? new Error(response.data)
        : new Error(`Server return status code ${response.status}`);
    }
  } else if (reason.message) {
    return new Error(reason.message);
  }
  return new Error('Unknown Error');
}

// eslint-disable-next-line
export function AxiosPlugin(Vue: typeof _Vue, options?: any): void {
  Vue.prototype.$http = Axios;
  Vue.prototype.$httpError = reasonToError;
}

declare module 'vue/types/vue' {
  interface Vue {
    readonly $http: AxiosStatic;
    /**
     * convert http reject reason to error object
     */
    $httpError: typeof reasonToError;
  }
}
