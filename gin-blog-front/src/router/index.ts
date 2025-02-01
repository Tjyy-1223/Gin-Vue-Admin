import { createRouter, createWebHistory } from 'vue-router'

const basicRoutes = [
    {
      name: 'Home',
      path: '/',
      component: () => import('@/App.vue')
    },
]

export const router = createRouter({
    history: createWebHistory('/'),
    routes: basicRoutes,
    scrollBehavior: () => ({ left: 0, top: 0 }),
  })

