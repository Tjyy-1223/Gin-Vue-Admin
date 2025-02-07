import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore, usePermissionStore, useUserStore } from '@/store'
import { setupRouterGuard } from './guard'


// 基础路由, 无需权限, 总是会注册到最终路由中
export const basicRoutes = [
    {
        name: 'Login',
        path: '/login',
        component: () => import('@/views/Login.vue'),
        isHidden: true,
        meta: {
            title: '登录页',
        },
    },
    {
        name: '404',
        path: '/404',
        component: () => import('@/views/error-page/404.vue'),
        isHidden: true,
        meta: {
            title: '错误页',
        },
    },
]

export const router = createRouter({
    history: createWebHistory(import.meta.env.VITE_PUBLIC_PATH), // '/admin'
    routes: basicRoutes,
    scrollBehavior: () => ({ left: 0, top: 0 }),
})


/**
 * 初始化路由
 */
export async function setupRouter(app) {
    await addDynamicRoutes()
    setupRouterGuard(router)
    app.use(router)
}


/**
 * ! 添加动态路由: 根据配置由前端或后端生成路由
 */
export async function addDynamicRoutes() {
    const authStore = useAuthStore()

    if (!authStore.token) {
        authStore.toLogin()
        return
    }

    // 有 token 的情况
    try {
        const userStore = useUserStore()
        const permissionStore = usePermissionStore()

        // userId 不存在, 则调用接口根据 token 获取用户信息
        if (!userStore.userId) {
            await userStore.getUserInfo()
        }

        // 根据环境变量中的值决定前端生成路由还是后端路由
        const accessRoutes = JSON.parse(import.meta.env.VITE_BACK_ROUTER)
            ? await permissionStore.generateRoutesBack() // ! 后端生成路由
            : permissionStore.generateRoutesFront(['admin']) // ! 前端生成路由 (根据角色), 待完善
        console.log(accessRoutes)

        // 打印所有路由信息
        console.log(router.getRoutes());

        accessRoutes.forEach(route => {
            // 检查父路由名称是否冲突
            if (router.hasRoute(route.name)) {
                console.warn(`Route with name "${route.name}" already exists.`);
                return; // 跳过已存在的父路由
            }

            // 检查子路由名称是否冲突，并修改冲突的子路由名称
            if (route.children) {
                route.children.forEach(child => {
                    if (router.hasRoute(child.name) || child.name === route.name) {
                        // 如果子路由名称与父路由或已存在的路由冲突
                        console.warn(`Child route with name "${child.name}" conflicts with parent or existing route. Renaming...`);
                        // 修改子路由名称
                        child.name = `${route.name}-${child.name}`;
                    }
                });
            }

            // 添加父路由
            router.addRoute(route);
        });

        // 添加新的动态路由
        // accessRoutes.forEach(route => {
        //     if (!router.hasRoute(route.name)) {
        //         console.log(route.name)
        //         console.log(router.hasRoute(route.name))
        //         router.addRoute(route);
        //     }
        // });
    }
    catch (err) {
        console.error('addDynamicRoutes Error: ', err)
    }
}


/**
 * 重置路由
 */
export async function resetRouter() {
    router.getRoutes().forEach((route) => {
        const name = route.name
        if (!basicRoutes.some(e => e.name === name) && router.hasRoute(name)) {
            router.removeRoute(name)
        }
    })
}



