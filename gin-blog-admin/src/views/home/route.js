const Layout = () => import('@/layout/index.vue')

export default {
  name: 'Home',  // 父路由的名称
  path: '/',
  component: Layout,
  redirect: '/home',  // 重定向到子路由
  meta: {
    order: 0,
  },
  isCatalogue: true,
  children: [
    {
      name: 'HomePage',  // 修改子路由名称为不同的名称
      path: 'home',
      component: () => import('./index.vue'),
      meta: {
        title: '首页',
        icon: 'ic:sharp-home',
        order: 0,
      },
    },
  ],
}
