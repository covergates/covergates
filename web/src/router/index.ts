import Vue from 'vue';
import VueRouter, { RouteConfig } from 'vue-router';

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
      }
    ]
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue')
  }
];

const router = new VueRouter({
  mode: 'history',
  base: process.env.NODE_ENV === 'production' ? VUE_BASE : process.env.BASE_URL,
  routes
});

export default router;
