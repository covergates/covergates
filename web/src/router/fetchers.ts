import { Route, NavigationGuardNext } from 'vue-router';
import { Store } from 'vuex';
import { RootState, Actions, State, Mutations } from '@/store';
import { FetchReportOption } from '@/store/modules/report/actions';

type RouteHandler = (to: Route, from: Route, next: NavigationGuardNext) => void;

function fetchReportHistory(store: Store<RootState>) {
  const report = (store.state as State).report.current;
  const repo = (store.state as State).repository.current;
  if (report && repo) {
    store.dispatch(Actions.FETCH_REPORT_HISTORY, {
      ReportID: report.reportID,
      Ref: repo.Branch
    } as FetchReportOption);
  } else {
    store.commit(Mutations.SET_REPORT_HISTORY, []);
  }
}

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
          ).then(() => {
            fetchReportHistory(store);
            store.dispatch(Actions.FETCH_REPOSITORY_COMMITS);
            store.dispatch(Actions.FETCH_REPOSITORY_BRANCHES);
            store.dispatch(Actions.FETCH_REPOSITORY_OWNER);
          }).finally(() => {
            next();
          });
          store.dispatch(Actions.FETCH_REPOSITORY_SETTING);
        } else {
          next();
        }
      }).catch(reason => {
        console.warn(reason);
        next();
      });
  };
}

export function fetchNewRepository(store: Store<RootState>): RouteHandler {
  return (to, from, next) => {
    if (from.name !== null && to.query.ref !== from.query.ref && to.meta.checkRenew) {
      fetchCurrentRepository(store)(to, from, next);
    } else {
      next();
    }
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

export function fetchUserSettings(store: Store<RootState>): RouteHandler {
  return (to, from, next) => {
    store.dispatch(Actions.FETCH_USER_TOKENS);
    next();
  };
}
