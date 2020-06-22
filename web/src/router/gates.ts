import { Store } from 'vuex';
import { State } from '@/store';
import { NavigationGuard } from 'vue-router';

export function authorize(store: Store<State>, window: Window): NavigationGuard {
    return (to, from, next) => {

    }
}
