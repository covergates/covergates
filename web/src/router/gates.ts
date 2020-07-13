import { Store } from 'vuex';
import { NavigationGuard, Location } from 'vue-router';
import { RootState, State } from '@/store';

export function authorize(store: Store<RootState>): NavigationGuard {
  return (to, from, next) => {
    if (to.meta && to.meta.requiresAuth && (!(store.state as State).user.current || (store.state as State).user.current.error)) {
      next({
        name: 'Login'
      } as Location);
    } else {
      next();
    }
  };
}
