import Vue from 'vue';
import VueRouter, { RouteConfig } from 'vue-router';
import {
  fetchReportSource,
  fetchCurrentRepository,
  fetchUserSCM,
  fetchNewRepository
} from './fetchers';
import store from '@/store';

Vue.use(VueRouter);

const routes: Array<RouteConfig> = [
  {
    path: '/',
    name: 'Main',
    component: () => import('@/views/Main.vue'),
    children: [
      {
        path: '/',
        name: 'Home',
        component: () => import('@/views/Home.vue')
      },
      {
        path: '/user',
        name: 'User',
        meta: { requiresAuth: true },
        component: () => import('@/views/User.vue'),
        beforeEnter: fetchUserSCM(store)
      },
      {
        path: '/repo',
        name: 'Repo',
        meta: { requiresAuth: true },
        component: () => import('@/views/Repo.vue')
      },
      {
        path: '/report/:scm/:namespace/:name',
        name: 'Report',
        meta: { requiresAuth: true },
        component: () => import('@/views/Report.vue'),
        beforeEnter: fetchCurrentRepository(store),
        children: [
          {
            path: '/report/:scm/:namespace/:name',
            name: 'report-overview',
            meta: { checkRenew: true },
            component: () => import('@/components/ReportOverview.vue')
          },
          {
            path: 'history',
            name: 'report-history',
            meta: { checkRenew: true },
            component: () => import('@/components/ReportHistory.vue')
          },
          {
            path: 'code',
            name: 'report-code',
            meta: { requiresAuth: false, checkRenew: true },
            component: () => import('@/components/ReportFiles.vue')
          },
          {
            path: 'code/:path+',
            name: 'report-source',
            meta: { requiresAuth: true, checkRenew: true },
            component: () => import('@/components/ReportSource.vue'),
            beforeEnter: fetchReportSource(store)
          },
          {
            path: 'setting',
            name: 'report-setting',
            meta: { requiresAuth: true },
            component: () => import('@/components/ReportSetting.vue')
          }
        ]
      }
    ]
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue')
  },
  {
    path: '*',
    name: 'NoRoute',
    redirect: '/'
  }
];

const router = new VueRouter({
  mode: 'history',
  base: process.env.NODE_ENV === 'production' ? VUE_BASE : process.env.BASE_URL,
  routes
});
router.beforeEach(fetchNewRepository(store));

export default router;
