import { Route, NavigationGuardNext } from 'vue-router';
import { Store } from 'vuex';
import { RootState, Actions, State } from '@/store';
import { FetchReportOption } from '@/store/modules/report/actions';

type RouteHandler = (to: Route, from: Route, next: NavigationGuardNext) => void;

export function fetchCurrentRepository(store: Store<RootState>): RouteHandler {
  return (to, from, next) => {
    store.dispatch(Actions.CHANGE_CURRENT_REPOSITORY, to.params)
      .then(() => {
        if ((store.state as State).repository.current) {
          store.dispatch(
            Actions.FETCH_REPORT_CURRENT,
            {
              ReportID: (store.state as State).repository.current?.ReportID,
              Ref: to.query.ref
            } as FetchReportOption
          )
            .then(() => {
              store.dispatch(Actions.FETCH_REPOSITORY_OWNER);
            });
          store.dispatch(Actions.FETCH_REPOSITORY_SETTING);
        }
      });
    next();
  };
}

export function fetchReportSource(store: Store<RootState>): RouteHandler {
  return (to, from, next) => {
    store.dispatch(Actions.FETCH_REPORT_SOURCE, to);
    next();
  };
}

export function fetchReportSetting(store: Store<RootState>): RouteHandler {
  return (to, from, next) => {
    store.dispatch(Actions.FETCH_REPOSITORY_SETTING);
    next();
  };
}

export function fetchUserSCM(store: Store<RootState>): RouteHandler {
  return (to, from, next) => {
    store.dispatch(Actions.FETCH_USER_SCM);
    next();
  };
}
