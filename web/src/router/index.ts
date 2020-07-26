import Vue from 'vue';
import VueRouter, { RouteConfig } from 'vue-router';
import {
  fetchReportSource,
  fetchCurrentRepository,
  fetchUserSCM
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
            path: '',
            name: 'report-overview',
            component: () => import('@/components/ReportOverview.vue')
          },
          {
            path: 'code',
            name: 'report-code',
            meta: { requiresAuth: true },
            component: () => import('@/components/ReportFiles.vue')
          },
          {
            path: 'code/:path+',
            name: 'report-source',
            meta: { requiresAuth: true },
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

export default router;
