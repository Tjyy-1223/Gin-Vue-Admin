import { defineStore } from 'pinia'
import { shallowRef } from 'vue'
import { asyncRoutes, basicRoutes, vueModules } from '@/router/routes'

export const usePermissionStore = defineStore('permission', {
    persist: {
        key: 'gvb_admin_permission',
    },
    state: () => ({
        accessRoutes: [], // 可访问的路由
    }),
    getters: {
        // 最终可访问路由 = 基础路由 + 可访问的路由
        routes: state => basicRoutes.concat(state.accessRoutes),
        // 过滤掉 hidden 的路由作为左侧菜单显示
        menus: state => state.routes.filter(route => route.name && !route.isHidden),
    },
    actions: {

    },
})