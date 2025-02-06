## 第七章 gin-blog-admin 首页搭建

数据库中的默认用户：

- 管理员 admin 123456
- 普通用户 user 123456
- 测试用户 test 123456

## 7.0 初始化项目

**进入终端并使用命令行：** 

1. pnpm create vite -> gin-blog-admin -> vue -> ts -> 初始化完毕

2. 安装 vue 核心依赖：cd gin-blog-admin -> pnpm i
3. 运行项目 : pnpm run dev
4. 项目起步完成

------

**之后，将下面的内容放在 package.json 之后， 执行 pnpm i，完成相关配置：**

```
{
  "name": "gin-blog-front",
  "private": true,
  "version": "0.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite --open",
    "build": "vue-tsc -b && vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "@iconify/vue": "^4.3.0",
    "@unocss/reset": "^65.4.3",
    "@vueuse/core": "^12.5.0",
    "axios": "^1.7.9",
    "dayjs": "^1.11.13",
    "easy-typer-js": "^2.1.0",
    "highlight.js": "^11.11.1",
    "marked": "^15.0.6",
    "nprogress": "^0.2.0",
    "pinia": "^2.3.1",
    "pinia-plugin-persistedstate": "^4.2.0",
    "rollup-plugin-visualizer": "^5.14.0",
    "unocss": "^65.4.3",
    "v3-infinite-loading": "^1.3.2",
    "vite-plugin-compression": "^0.5.1",
    "vue": "^3.5.13",
    "vue-router": "^4.5.0",
    "vue3-danmaku": "^1.6.1"
  },
  "devDependencies": {
    "@egoist/tailwindcss-icons": "^1.8.2",
    "@iconify-json/mdi-light": "^1.2.2",
    "@iconify/json": "^2.2.302",
    "@iconify/tailwind": "^1.2.0",
    "@mdi/font": "^7.4.47",
    "@types/node": "^22.12.0",
    "@vitejs/plugin-vue": "^5.2.1",
    "@vue/tsconfig": "^0.7.0",
    "autoprefixer": "^10.4.20",
    "postcss": "^8.5.1",
    "sass-embedded": "^1.83.4",
    "tailwindcss": "^3.4.17",
    "typescript": "~5.6.2",
    "vite": "^6.0.5",
    "vue-tsc": "^2.2.0"
  }
}
```

**postcss.config.js**

```
export default {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
}
```

**tsconfig.app.json**

```json
{
  "extends": "@vue/tsconfig/tsconfig.dom.json",
  "compilerOptions": {
    "tsBuildInfoFile": "./node_modules/.tmp/tsconfig.app.tsbuildinfo",

    /* Linting */
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    "noUncheckedSideEffectImports": true,
    "lib": ["es6", "dom", "es2017"],
    "target": "es5",
    "module": "ESNext",
    "moduleResolution": "node",
    "jsx": "preserve",
    "allowJs": true,
    "baseUrl": "./",
    "paths": {
      "@/*": ["src/*"]  // 将 '@' 映射到 src 目录
    },
    "types": ["node","pinia"],
    "esModuleInterop": true,
    "skipLibCheck": true,
  },
  "include": ["src/**/*.ts", "src/**/*.tsx", "src/**/*.vue"]
}
```

**tailwind.config.js 中使用**

```javascript
/** @type {import('tailwindcss').Config} */

const { iconsPlugin, getIconCollections } = require('@egoist/tailwindcss-icons');
export default {
  content: ["index.html", "./src/**/*.{html,js,ts,jsx,tsx,vue}"],
  theme: {
    extend: {
      transitionDuration: {
        '500': '500ms',
      },
    },
  },
  plugins: [
    iconsPlugin({
      collections: getIconCollections(['mdi', 'lucide']), // 或使用 "all" 来使用全部图标
    }),
  ],
}
```

**修改 vite.config.ts 如下：**

```typescript
import path from 'path'
import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import unocss from 'unocss/vite'
import viteCompression from 'vite-plugin-compression'
import { visualizer } from 'rollup-plugin-visualizer'

export default defineConfig((configEnv) => {
  const env = loadEnv(configEnv.mode, process.cwd())

  return {
    base: env.VITE_PUBLIC_PATH || '/',
    resolve: {
      alias: {
        '@': path.resolve(__dirname, 'src'),
        '~': path.resolve(process.cwd()),
      },
    },
    plugins: [
      vue(),
      unocss(),
      viteCompression({ algorithm: 'gzip' }),
      visualizer({ open: false, gzipSize: true, brotliSize: true }),
    ],
    server: {
      host: '0.0.0.0',
      port: 3000,
      open: false,
      proxy: {
        '/api': {
          target: env.VITE_SERVER_URL,
          changeOrigin: true,
        },
      },
    },
    // https://cn.vitejs.dev/guide/api-javascript.html#build
    build: {
      chunkSizeWarningLimit: 1024, // chunk 大小警告的限制（单位kb）
    },
    esbuild: {
      drop: ['debugger'], // console
    },
  }
})
```

根目录下创建 uno.config.js

```javascript
import {
  defineConfig,
  presetIcons,
  presetTypography,
  presetUno,
  transformerDirectives,
  transformerVariantGroup,
} from 'unocss'

export default defineConfig({
  shortcuts: [
    ['f-c-c', 'flex justify-center items-center'],
  ],
  presets: [
    presetUno(),
    presetIcons({ warn: true }),
    presetTypography(),
  ],
  transformers: [
    transformerDirectives(),
    transformerVariantGroup(),
  ],
})

```

**然后，可以将 main.ts 编写如下：**

```typescript
import { createApp } from 'vue'
import App from './App.vue'

// unocss
import 'uno.css'
import '@unocss/reset/tailwind.css'

const app = createApp(App);
app.mount('#app')
```

之后，就可以开始编写对应的 vue 代码了

使用 pnpm run dev 查看目前的效果



## 7.1 配置路由 router

首先编写一个基础版本的 AuthStore, 提供一个 token 相关的状态，完整版本会在 7.2 中补全

src/store/modules/auth.js

```javascript
import { unref } from "vue";
import { defineStore } from "pinia";

export const useAuthStore = defineStore('auth', {
    persist: {
        key: 'gvb_admin_auth',
        paths: ['token'],
    },
    state: () => ({
        token: null,
    }),
    actions: {
        setToken(token) {
            this.token = token
        },
    },
    getters: {

    }
})
```

这段代码是一个使用 Vue 和 Pinia 的组合，用于创建一个全局的状态管理模块（Store），专门用于处理用户认证（Authentication）相关的数据。以下是代码的详细解析和作用说明：

### 1. **引入依赖**

JavaScript复制

```javascript
import { unref } from "vue";
import { defineStore } from "pinia";
```

- **`unref`**：这是 Vue 的一个工具函数，用于获取响应式引用（ref）的实际值。虽然在这段代码中没有直接使用 `unref`，但它可能是为了后续扩展或其他地方的代码复用而引入的。
- **`defineStore`**：这是 Pinia 的核心函数，用于定义一个 Pinia Store。Pinia 是 Vue 的状态管理库，类似于 Vuex，但更简洁、更灵活。

### 2. **定义 Store**

```javascript
export const useAuthStore = defineStore('auth', {
    // 配置持久化选项
    persist: {
        key: 'gvb_admin_auth', // 持久化存储的键名
        paths: ['token'],      // 只持久化 `token` 属性
    },
    // 定义状态
    state: () => ({
        token: null, // 存储用户认证的 token，初始值为 null
    }),
    // 定义操作方法
    actions: {
        setToken(token) {
            this.token = token; // 设置 token 的值
        },
    },
    // 定义计算属性（getters）
    getters: {
        // 这里可以定义计算属性，但目前为空
    }
});
```

这段代码通过 Pinia 定义了一个简单的认证状态管理模块，用于存储和操作用户的认证令牌（Token），并支持持久化存储，确保用户登录状态在页面刷新后仍然有效。

- **`state`**：定义了 Store 的状态，这里只有一个 `token` 属性，用于存储用户的认证令牌（Token）。初始值为 `null`，表示用户未登录时的状态。
- **持久化**：通过 `persist` 配置，将 `token` 持久化到本地存储（通常是 `localStorage` 或 `sessionStorage`）中，键名为 `gvb_admin_auth`。这样即使页面刷新，用户的登录状态（Token）也不会丢失。

- **`actions`**：定义了操作状态的方法。这里有一个 `setToken` 方法，用于设置 `token` 的值。当用户登录成功时，可以调用 `setToken` 方法将 Token 存储到 Store 中。



#### 7.1.1 guard.js

这段代码主要用于在前端路由中设置三个导航守卫，以改善用户体验和页面安全性。具体来说：

1. **页面加载进度条守卫（createPageLoadingGuard）**
   - 在路由切换开始时，通过 `router.beforeEach` 调用 `window.$loadingBar?.start()` 开启顶部加载条，告知用户正在加载新页面。
   - 在路由切换完成后，通过 `router.afterEach` 延时 200 毫秒调用 `window.$loadingBar?.finish()` 关闭加载条。
   - 如果在导航过程中发生未捕获的错误，使用 `router.onError` 调用 `window.$loadingBar?.error()`，将加载条状态设为错误状态。
2. **权限校验守卫（createPermissionGuard）**
   - 在路由跳转前，通过 `router.beforeEach` 检查用户是否存在 Token（通过 `useAuthStore()` 获取 token）。
   - 没有 Token 的情况：
     - 如果目标页面是 `/login` 或 `/404`，允许访问。
     - 否则提示“没有 Token, 请先登录！”，并重定向到登录页面，同时携带目标页面路径作为重定向参数（方便登录后自动跳转到原本想去的页面）。
   - 有 Token 的情况：
     - 如果目标页面为登录页，则提示“已登录，无需重复登录！”并重定向到首页。
     - 如果目标路由在路由配置中存在，则允许访问。
     - 如果目标路由不存在，则重定向到 404 页面，并传递当前完整路径作为参数（用于显示未找到的路径信息）。
   - （注：代码中还预留了刷新 Token 的 TODO 注释，可由后端来实现更完善的权限控制。）
3. **页面标题设置守卫（createPageTitleGuard）**
   - 在路由切换完成后，通过 `router.afterEach` 获取目标路由的 `meta.title`，并将 `document.title` 更新为 `页面标题 | 基础标题` 的格式。其中基础标题从环境变量 `VITE_TITLE` 中读取。如果目标页面没有设定 `meta.title`，则直接显示基础标题。

总体来说，这段代码通过在路由切换过程中动态处理加载动画、权限验证和页面标题设置，提高了用户体验和页面安全性。

```javascript
import { useAuthStore } from "../store/modules/auth";

export function setupRouterGuard(router) {
    createPageLoadingGuard(router)
    createPermissionGuard(router)
    createPageTitleGuard(router)
}

/**
 * 根据导航设置顶部加载条的状态
 */
function createPageLoadingGuard(router) {
    router.beforeEach(() => window.$loadingBar?.start())
    router.afterEach(() => setTimeout(() => window.$loadingBar?.finish(), 200))
    // 在导航期间每次发生未捕获的错误时都会调用该处理程序
    router.onError(() => window.$loadingBar?.error())
}

/**
 * 根据有无 Token 判断能否访问页面
 */
function createPermissionGuard(router) {
    // const base = import.meta.env.VITE_BASE_URL
    // 路由前置守卫: 根据有没有 Token 判断前往哪个页面
    router.beforeEach(async (to) => {
        const { token } = useAuthStore()

        // 没有 Token
        if (!token) {
            // login 和 404 不需要 token 即可访问
            if (['/login', '/404'].includes(to.path)) {
                return true
            }

            window.$message.error('没有 Token, 请先登录!')
            // 重定向到登录页, 并且携带 redirect 参数, 登录后自动重定向到原本的目标页面
            return { name: 'Login', query: { ...to.query, redirect: to.path } }
        }


        // 有 Token 的时候无需访问登录页面
        if (to.name === 'Login') {
            window.$message.success('已登录，无需重复登录！')
            return { path: '/' }
        }

        // 能在路由中找到, 则正常访问
        if (router.getRoutes().find(e => e.name === to.name)) {
            return true
        }

        // TODO: 刷新 Token - 可以交给后端去做
        // await refreshAccessToken()

        // TODO: 判断是无权限还是 404
        return { name: '404', query: { path: to.fullPath } }
    })

}

/**
 * 根据路由元信息设置页面标题
 */
function createPageTitleGuard(router) {
    const baseTitle = import.meta.env.VITE_TITLE
    router.afterEach((to) => {
        const pageTitle = to.meta?.title
        document.title = pageTitle ? `${pageTitle} | ${baseTitle}` : baseTitle
    })
}
```



#### 7.1.2 routes.js

这段代码的作用是：

1. **定义基础路由**：定义了项目中始终需要的路由，如登录页和404错误页。
2. **动态加载路由模块**：通过 `import.meta.glob` 动态加载 `@/views` 目录下每个模块的 `route.js` 文件，并将这些路由模块存储到 `asyncRoutes` 数组中。
3. **动态加载组件**：通过 `import.meta.glob` 动态加载 `@/views` 目录下每个模块的 `index.vue` 文件，供按需加载组件使用。
4. **模块化和可扩展性**：通过动态加载的方式，使得项目路由和组件的管理更加灵活，便于扩展和维护。

```javascript
import { useAuthStore } from '@/store/modules/auth'

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

// 前端控制路由: 加载 views 下每个模块的 routes.js 文件
const routeModules = import.meta.glob('@/views/**/route.js', { eager: true })
const asyncRoutes = []
Object.keys(routeModules).forEach((key) => {
  asyncRoutes.push(routeModules[key].default)
})

// 加载 views 下每个模块的 index.vue 文件
const vueModules = import.meta.glob('@/views/**/index.vue')

export {
  asyncRoutes,
  vueModules,
}
```

**`basicRoutes`**：定义了基础路由，这些路由是项目中始终需要注册的路由，通常包括登录页和404错误页。

- **`name`**：路由的名称，用于在代码中引用路由。
- **`path`**：路由的路径。
- **`component`**：路由对应的组件，这里使用了动态加载的方式（`import()`），可以按需加载组件，减少初始加载时间。
- **`isHidden`**：自定义字段，可能用于标识该路由是否在侧边栏或其他导航中隐藏。
- **`meta`**：路由元信息，可以存储额外的数据，例如页面标题。

------

**`import.meta.glob`**：这是 Vite 提供的一个方法，用于动态加载匹配特定模式的模块。这里加载了 `@/views` 目录下所有模块的 `route.js` 文件。

- **`@/views/\**/route.js`**：匹配 `@/views` 目录及其子目录下的所有 `route.js` 文件。
- **`{ eager: true }`**：表示立即加载所有匹配的模块，而不是按需加载。

**`Object.keys(routeModules).forEach`**：遍历 `routeModules` 的键，将每个模块的默认导出（`default`）推入 `asyncRoutes` 数组。

------

**`vueModules`**：使用 `import.meta.glob` 动态加载 `@/views` 目录下所有模块的 `index.vue` 文件。

- 这些文件通常是一个模块的主组件文件，可能用于按需加载组件，优化应用性能。

将 `asyncRoutes` 和 `vueModules` 导出，供其他模块使用。例如：

- **`asyncRoutes`**：可以在路由配置中动态注册这些路由。
- **`vueModules`**：可以在需要的地方按需加载组件。



#### 7.1.3 index.js

首先编写两个基础版本的 store，完整版本会在 7.2 中补全

src/store/modules/permission.js

```javascript
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
```

src/store/modules/user.js

```javascript
import { defineStore } from 'pinia'

// 用户全局变量
export const useUserStore = defineStore('user', {
    state: () => ({
        userInfo: {
            id: null,
            nickname: '',
            avatar: '',
            intro: '',
            website: '',
            // roles: [], // TODO: 后端返回 roles
        },
    }),
    getters: {
        userId: state => state.userInfo.id,
        nickname: state => state.userInfo.nickname,
        intro: state => state.userInfo.intro,
        website: state => state.userInfo.website,
        // avatar: state => convertImgUrl(state.userInfo.avatar),
        // roles: state => state.userInfo.roles,
    },
    actions: {
    },
})
```

------

**然后，我们完成 index.js 的书写，如下：**

- **动态路由：** 根据用户的登录状态和权限，动态加载前端或后端生成的路由，确保用户只能访问有权限的页面。
- **路由守卫：** 配置路由守卫来处理权限、页面加载等逻辑。
- **路由重置：** 提供了重置路由的功能，允许在需要时清空动态添加的路由。

**这种结构能够灵活地根据用户状态和权限动态地生成路由配置，同时提供了管理和控制路由的能力。**

```javascript
import { createRouter, createWebHistory } from 'vue-router'

import { basicRoutes } from './routes'
import { setupRouterGuard } from './guard'

import { useAuthStore, usePermissionStore, useUserStore } from '@/store'

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

        // 将当前没有的路由添加进去
        accessRoutes.forEach(route => !router.hasRoute(route.name) && router.addRoute(route))
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
```

**`setupRouter` 是一个初始化路由的方法。它会在应用启动时调用，主要执行以下几项任务：**

1. 调用 `addDynamicRoutes()` 方法，动态加载路由。
2. 调用 `setupRouterGuard(router)` 方法，为路由设置守卫（例如：权限验证、页面标题、加载条等）。
3. 最后，使用 `app.use(router)` 将路由实例挂载到 Vue 应用中。

------

**`addDynamicRoutes`：这个方法用于根据用户权限动态添加路由。逻辑如下：**

1. 首先获取 `authStore` 中的 `token`，如果没有 Token，调用 `authStore.toLogin()` 方法，跳转到登录页。
2. 如果有 Token，继续处理用户信息和权限：
   - 如果用户没有 `userId`（通常意味着未加载用户信息），调用 `userStore.getUserInfo()` 方法从后端获取用户信息。
   - 根据环境变量 VITE_BACK_ROUTER  来决定路由是由后端生成还是前端生成：
     - 如果环境变量为 `true`，调用 `permissionStore.generateRoutesBack()` 从后端获取路由。
     - 如果环境变量为 `false`，则根据角色（如 `'admin'`）调用 `permissionStore.generateRoutesFront()` 来在前端生成路由。
3. 使用 `router.addRoute()` 添加动态生成的路由，确保没有重复的路由。

------

**`resetRouter` 用于重置路由，通常用于用户登出或权限变更时刷新路由配置。**

1. 获取所有现有路由，通过 `router.getRoutes()` 获取路由列表。
2. 遍历每个路由，如果该路由名称不在基本路由 `basicRoutes` 中，且路由已经存在（通过 `router.hasRoute()` 检查），则通过 `router.removeRoute()` 删除它。
3. 这样可以移除动态生成的路由，恢复到初始的基本路由配置。

### 



## 7.2 主页搭建 - Store 仓库

首先复制三个配置文件到根目录下，启动时会设置好相关的后端请求路径配置:

+ .env
+ .env.development
+ .env.production

之后，建立如下的 store 仓库

#### 7.2.1 src/store/modules/premission.js

```javascript
import { defineStore } from 'pinia'
import { shallowRef } from 'vue'
import { asyncRoutes, basicRoutes, vueModules } from '@/router/routes'
import Layout from '@/layout/index.vue'
import api from '@/api'

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
    // ! 后端生成路由: 后端返回的就是最终路由, 处理成前端格式
    async generateRoutesBack() {
      const resp = await api.getUserMenus() // 调用接口获取后端传来的路由
      this.accessRoutes = buildRoutes(resp.data) // 处理成前端路由格式
      return this.accessRoutes
    },
    // ! 前端控制路由权限: 根据角色过滤路由
    generateRoutesFront(role = []) {
      this.accessRoutes = filterAsyncRoutes(asyncRoutes, role)
      return this.accessRoutes
    },
    resetPermission() {
      this.$reset()
    },
  },
})

// ! 前端路由相关函数
/**
 * 前端过滤出有权限访问的路由
 */
function filterAsyncRoutes(routes = [], role) {
  const result = []
  routes.forEach((route) => {
    if (hasPermission(route, role)) {
      const curRoute = { ...route, children: [] }
      if (route.children?.length) {
        curRoute.children = filterAsyncRoutes(route.children, role)
      }
      else {
        Reflect.deleteProperty(curRoute, 'children')
      }
      result.push(curRoute)
    }
  })
  return result
}

/**
 * 前端判断用户角色是否有权限访问路由
 */
function hasPermission(route, role) {
  // 如果该路由不需要权限直接返回 true
  if (!route.meta?.requireAuth) {
    return true
  }
  // 路由需要的角色
  const routeRole = route.meta?.role ?? []
  // 登录用户没有角色 或者 路由没有设置角色判断, 为没有权限
  if (!role.length || !routeRole.length) {
    return false
  }
  // 路由指定的角色包含任一登录用户角色则判定有权限
  return role.some(item => routeRole.includes(item))
}

// ! 后端路由相关函数
// 根据后端传来数据构建出前端路由
function buildRoutes(routes = []) {
  const result = []

  for (const e of routes) {
    if (e.is_catalogue) {
      result.push({
        name: e.name,
        path: '/', // *
        component: shallowRef(Layout),
        isHidden: e.is_hidden,
        isCatalogue: true, // *
        redirect: e.redirect,
        meta: {
          order: e.order_num,
        },
        children: [{
          name: e.name,
          path: e.path,
          component: vueModules[`/src/views${e.component}/index.vue`],
          meta: {
            title: e.name,
            icon: e.icon,
            keepAlive: e.keep_alive, // TODO:
            order: 0,
          },
        }],
      })
    }
    else {
      result.push({
        name: e.name,
        path: e.path,
        component: shallowRef(Layout),
        isHidden: e.is_hidden,
        redirect: e.redirect,
        meta: {
          title: e.name,
          icon: e.icon,
          keepAlive: e.keep_alive, // TODO:
          order: e.order_num,
        },
        children: e.children?.map(ee => ({
          name: ee.name,
          path: ee.path, // 父路径 + 当前菜单路径
          component: vueModules[`/src/views${ee.component}/index.vue`],
          isHidden: ee.is_hidden,
          meta: {
            title: ee.name,
            icon: ee.icon,
            order: ee.order_num,
            keepAlive: ee.keep_alive,
          },
        })),
      })
    }
  }

  return result
}
```

**我们逐个分析其中的每个函数：**

**function filterAsyncRoutes(routes = [], role)**

作用总结

- **过滤路由**: 主要用于根据用户角色权限过滤出用户可以访问的路由。
- **支持嵌套路由**: 递归遍历路由和子路由，确保所有子路由也被权限过滤。
- **动态路由控制**: 通过递归和权限判断，动态决定哪些路由应当出现在应用中，哪些应当被隐藏或屏蔽。

```javascript
// ! 前端路由相关函数
/**
 * 前端过滤出有权限访问的路由
 */
function filterAsyncRoutes(routes = [], role) {
  const result = []
  routes.forEach((route) => {
    if (hasPermission(route, role)) {
      const curRoute = { ...route, children: [] }
      if (route.children?.length) {
        curRoute.children = filterAsyncRoutes(route.children, role)
      }
      else {
        Reflect.deleteProperty(curRoute, 'children')
      }
      result.push(curRoute)
    }
  })
  return result
}
```

1. **输入参数**:
   - `routes`: 传入的路由列表，可以是嵌套的多层路由结构。
   - `role`: 当前用户的角色数组，用来判断用户是否有权限访问某个路由。
2. **函数流程**:
   - 创建一个空数组 `result` 用来保存最终符合权限的路由。
   - 遍历传入的 routes 数组中的每个路由对象，执行以下操作：
     - **权限判断**: 调用 `hasPermission(route, role)` 检查当前路由是否符合用户权限。`hasPermission` 函数会基于路由的 `meta.role` 配置和当前用户的 `role` 数组来判断是否有访问权限。
     - **满足权限条件时处理当前路由**:
       - 创建一个 `curRoute` 对象，使用 `...route` 创建路由对象的浅拷贝，并清空 `children` 属性（默认为空数组）。
       - 递归处理子路由:
         - 如果该路由有子路由（`route.children` 存在且长度大于 0），递归调用 `filterAsyncRoutes` 来处理子路由。
         - 如果该路由没有子路由，则删除 `children` 属性（`Reflect.deleteProperty(curRoute, 'children')`），因为不需要存储空的 `children`。
       - 将符合权限的路由（`curRoute`）推入 `result` 数组。
3. **返回结果**:
   - 最终，返回 `result` 数组，其中包含了所有符合权限的路由（包括所有符合权限的子路由）。



------

**function hasPermission(route, role)** 

**权限检查**: 这段代码根据路由的 `meta` 信息和用户的角色来判断是否允许访问某个路由。

- 如果路由不需要权限，用户可以直接访问。
- 如果路由设置了权限（通过 `requireAuth` 和 `role`），则检查用户的角色是否包含在路由要求的角色列表中。

**灵活性**: 该方法支持多角色配置，即一个路由可以设置多个角色，用户只需要拥有其中一个角色即可访问该路由。

```javascript
/**
 * 前端判断用户角色是否有权限访问路由
 */
function hasPermission(route, role) {
  // 如果该路由不需要权限直接返回 true
  if (!route.meta?.requireAuth) {
    return true
  }
  // 路由需要的角色
  const routeRole = route.meta?.role ?? []
  // 登录用户没有角色 或者 路由没有设置角色判断, 为没有权限
  if (!role.length || !routeRole.length) {
    return false
  }
  // 路由指定的角色包含任一登录用户角色则判定有权限
  return role.some(item => routeRole.includes(item))
}
```

`route`: 当前路由对象，包含该路由的所有信息，例如路由路径、组件、元信息（`meta`）等。

`role`: 当前用户的角色数组，包含用户所拥有的角色。

------

**function buildRoutes(routes = [])**

这段代码主要是将后端返回的路由数据转换成前端 Vue Router 可以使用的格式，使得前端路由配置能够根据后端的权限和配置动态生成。

- **目录类型路由**（`is_catalogue` 为 `true`）会作为父路由，包含至少一个子路由。子路由的组件路径通过 `vueModules` 动态导入。
- **普通路由**（`is_catalogue` 为 `false`）直接映射为页面路由，可以有多个子路由，子路由的路径和组件也会动态设置。
- **动态组件加载**: 通过 `vueModules` 动态引入子组件，使得路由的组件路径是灵活的，便于根据后端配置动态生成路由。
- **路由元信息**（`meta`）包含了每个路由的标题、图标、排序、是否需要缓存等信息，这些可以用于路由管理、菜单生成等场景。

```javascript
// ! 后端路由相关函数
// 根据后端传来数据构建出前端路由
function buildRoutes(routes = []) {
  const result = []

  for (const e of routes) {
    if (e.is_catalogue) {
      result.push({
        name: e.name,
        path: '/', // *
        component: shallowRef(Layout),
        isHidden: e.is_hidden,
        isCatalogue: true, // *
        redirect: e.redirect,
        meta: {
          order: e.order_num,
        },
        children: [{
          name: e.name,
          path: e.path,
          component: vueModules[`/src/views${e.component}/index.vue`],
          meta: {
            title: e.name,
            icon: e.icon,
            keepAlive: e.keep_alive, // TODO:
            order: 0,
          },
        }],
      })
    }
    else {
      result.push({
        name: e.name,
        path: e.path,
        component: shallowRef(Layout),
        isHidden: e.is_hidden,
        redirect: e.redirect,
        meta: {
          title: e.name,
          icon: e.icon,
          keepAlive: e.keep_alive, // TODO:
          order: e.order_num,
        },
        children: e.children?.map(ee => ({
          name: ee.name,
          path: ee.path, // 父路径 + 当前菜单路径
          component: vueModules[`/src/views${ee.component}/index.vue`],
          isHidden: ee.is_hidden,
          meta: {
            title: ee.name,
            icon: ee.icon,
            order: ee.order_num,
            keepAlive: ee.keep_alive,
          },
        })),
      })
    }
  }

  return result
}

```

**下面是对这段代码的具体分析：**

1. **函数作用**:

```javascript
function buildRoutes(routes = []) {
  const result = []
  // ... (逻辑处理)
  return result
}
```

- **`buildRoutes`** 函数接收一个包含路由信息的数组 `routes`（通常是从后端接口返回的路由数据）。
- 它通过处理这些路由数据，将每一条路由转换成符合 Vue Router 格式的路由配置，最终返回一个新的数组 `result`，其中包含所有构建好的路由。

2. **处理每个路由对象**:

```javascript
for (const e of routes) {
  if (e.is_catalogue) { 
    // 处理目录类型路由
  } else {
    // 处理普通路由
  }
}
```

- 遍历传入的路由数据 `routes`，对每个路由对象 `e` 进行处理。
- **`e.is_catalogue`**: 如果该路由是一个目录类型路由（例如：菜单或导航栏中的父项），会执行不同的处理逻辑，生成具有子路由的父级路由。
- 否则，将其处理为普通路由。

3. **目录类型路由（`e.is_catalogue` 为 `true`）**:

```javascript
result.push({
  name: e.name,
  path: '/',
  component: shallowRef(Layout),
  isHidden: e.is_hidden,
  isCatalogue: true,
  redirect: e.redirect,
  meta: { order: e.order_num },
  children: [{
    name: e.name,
    path: e.path,
    component: vueModules[`/src/views${e.component}/index.vue`],
    meta: {
      title: e.name,
      icon: e.icon,
      keepAlive: e.keep_alive,
      order: 0,
    },
  }],
})
```

**目录路由**: 目录类型的路由通常是用于组织或分组其他路由的父级路由。对于这种类型的路由：

- `name`: 设置路由的名称，通常用于引用路由。
- `path: '/'`: 路径设置为 `'/'`，通常表示这是一个父级路由，具有子路由。
- `component: shallowRef(Layout)`: 使用 `Layout` 作为父级路由的组件，`shallowRef` 是 Vue 3 中的 `ref` 用法，表示该组件引用是浅层的（不包含其子组件的依赖）。
- `children`: 这个目录路由会有一个子路由（一般只有一个），这个子路由指向具体的页面组件。

子路由配置：

- `name`: 子路由的名称，通常是 `e.name`。
- `path`: 子路由的路径，直接使用 `e.path`。
- `component`: 动态导入子路由的组件，路径由 `vueModules` 动态决定。路径形式为 `vueModules['/src/views{e.component}/index.vue']`。
- `meta`: 包含该路由的元数据，如 `title`（路由标题）、`icon`（路由图标）、`keepAlive`（是否缓存此页面）、`order`（排序）。

4. **普通路由（`e.is_catalogue` 为 `false`）**:

```javascript
result.push({
  name: e.name,
  path: e.path,
  component: shallowRef(Layout),
  isHidden: e.is_hidden,
  redirect: e.redirect,
  meta: {
    title: e.name,
    icon: e.icon,
    keepAlive: e.keep_alive,
    order: e.order_num,
  },
  children: e.children?.map(ee => ({
    name: ee.name,
    path: ee.path,
    component: vueModules[`/src/views${ee.component}/index.vue`],
    isHidden: ee.is_hidden,
    meta: {
      title: ee.name,
      icon: ee.icon,
      order: ee.order_num,
      keepAlive: ee.keep_alive,
    },
  })),
})
```

普通路由: 普通路由不属于目录类型，直接用于页面路由。对于这些路由，处理逻辑如下：

- `name`: 路由名称。
- `path`: 路由的路径。
- `component: shallowRef(Layout)`: 使用 `Layout` 作为父组件（可能是一个布局组件）。
- children: 该路由可以有子路由，使用 map遍历 e.children（如果有子路由）。对于每个子路由 ee，会动态构建子路由配置：
  - `name`: 子路由名称。
  - `path`: 子路由路径。
  - `component`: 动态加载子路由组件，路径为 `vueModules['/src/views{ee.component}/index.vue']`。
  - `meta`: 子路由的元信息（例如标题、图标、缓存配置等）。

5. **返回路由结果**:

```javascript
return result
```

最终，函数返回处理好的路由数组 `result`，这些路由已经被转换为 Vue Router 所需的格式，可以直接用来配置路由系统。





#### 7.2.2 src/store/modules/tag.js

这个 store 适用于需要标签页管理的应用，特别是那些支持多标签、多页面的后台管理系统。

+ **标签管理**: 该 store 实现了一个标签栏系统，允许用户动态添加、删除、关闭标签，并管理标签的状态。
+ **标签刷新**: 提供了刷新标签的功能，模拟标签的刷新效果。
+ **动态路由与标签**: 与路由（通过 `router.push`）和标签（通过 `tags`）密切集成，使得标签和路由的跳转同步更新。
+ **标签持久化**: 通过 `sessionStorage` 实现标签栏数据的持久化，使得页面刷新后标签栏能够恢复

```javascript
import { nextTick } from 'vue'
import { defineStore } from 'pinia'
import { router } from '@/router'

export const useTagStore = defineStore('tag', {
  persist: {
    key: 'gvb_admin_tag',
    paths: ['tags'],
    storage: window.sessionStorage,
  },
  state: () => ({
    tags: [], // 标签栏的所有标签
    activeTag: '', // 当前激活的标签 path
    reloading: true, // 是否正在刷新
    /**
     * ! keepAlive 路由的 key, 重新赋值可重置 keepAlive
     * key 是 route name
     */
    aliveKeys: {},
  }),
  getters: {
    // 获取当前激活的标签的索引
    activeIndex: state => state.tags.findIndex(tag => tag.path === state.activeTag),
  },
  actions: {
    /**
     * 更新 keepAlive 路由, 让其重新渲染
     * @param {string} name route name
     */
    updateAliveKey(name) {
      this.aliveKeys[name] = (+new Date())
    },
    /**
     * 设置当前激活的标签
     * @param {string} path 标签对应的路由路径
     */
    async setActiveTag(path) {
      await nextTick() // 将回调延迟到下次 DOM 更新循环之后执行
      this.activeTag = path
    },
    /**
     * 设置当前显示的所有标签
     * @param {string[]} tags 数组
     */
    setTags(tags) {
      this.tags = tags
    },
    /**
     * 添加标签 (不添加白名单中 和 已存在的)
     * @param {{ name, path, title, icon }} tag 标签对象
     * 添加新标签: 如果标签已存在（根据 path 判断），则更新该标签；否则，将新标签添加到标签数组中，并激活该标签。
     */
    addTag(tag = {}) {
      const index = this.tags.findIndex(item => item.path === tag.path)
      if (index !== -1) {
        this.tags.splice(index, 1, tag)
      }
      else {
        this.setTags([...this.tags, tag])
      }
      this.setActiveTag(tag.path)
    },
    /**
     * 移除标签 , 如果只有一个标签, 无法移除
     * @param {string} path 标签对应的路由路径
     */
    removeTag(path) {
      // 如果关闭的是当前标签
      if (path === this.activeTag) {
        if (this.activeIndex === 0) { // 如果是第一个标签, 则选中第二个标签
          router.push(this.tags[1].path)
        }
        else { // 否则选中左边的标签
          router.push(this.tags[this.activeIndex - 1].path)
        }
      }
      this.setTags(this.tags.filter(tag => tag.path !== path))
    },
    /**
     * 关闭其他标签
     * @param {string} path
     */
    removeOther(path = this.activeTag) {
      this.setTags(this.tags.filter(tag => tag.path === path))
      // 如果点击的不是当前标签, 会将当前标签关闭, 那么跳转到第一个标签
      if (path !== this.activeTag) {
        router.push(this.tags[0].path) // 关闭其他后只剩一个标签
      }
    },
    /**
     * 关闭左侧标签
     * @param {string} path
     */
    removeLeft(path) {
      const curIndex = this.tags.findIndex(item => item.path === path)
      // 过滤出右边的标签
      const filterTags = this.tags.filter((item, index) => index >= curIndex)
      this.setTags(filterTags)
      // 如果当前浏览的标签被关闭, 打开一个新标签
      if (!filterTags.find(item => item.path === this.activeTag)) {
        router.push(filterTags[filterTags.length - 1].path)
      }
    },
    /**
     * 关闭左侧标签
     * @param {string} path
     */
    removeRight(path) {
      const curIndex = this.tags.findIndex(item => item.path === path)
      // 过滤出左边的标签
      const filterTags = this.tags.filter((item, index) => index <= curIndex)
      this.setTags(filterTags)
      // 如果当前浏览的标签被关闭, 打开一个新标签
      if (!filterTags.find(item => item.path === this.activeTag)) {
        router.push(filterTags[filterTags.length - 1].path)
      }
    },
    /**
     * 重置标签
     */
    resetTags() {
      this.$reset()
    },
    /**
     * 刷新页面
     * @description 效果并非按 F5 刷新整个网页, 而是模拟刷新 (nextTick + 滚动到顶部)
     */
    async reloadTag() {
      window.$loadingBar.start()

      // 配合 v-if="reloadFlag" 实现白屏效果
      this.reloadFlag = false
      await nextTick() // 将回调延迟到下次 DOM 更新循环之后执行
      this.reloadFlag = true

      // 滚动到顶部, 模拟刷新
      setTimeout(() => {
        document.documentElement.scrollTo({ left: 0, top: 0 })
        window.$loadingBar.finish()
      }, 100)
    },
  },
})

```



#### 7.2.3 src/store/modules/theme.js

这段代码定义了一个用于管理主题相关状态的 Pinia Store，包括侧边栏折叠状态、水印显示状态和暗色模式状态。通过持久化配置，`collapsed` 和 `watermarked` 状态会在页面刷新后仍然保持。同时，通过 `useDark`，可以动态检测和切换暗色模式。

- **`useDark`**：获取当前主题是否为暗色模式，并将其值赋给 `isDark`。
- **`defineStore`**：定义了一个名为 `theme-store` 的 Store。
  - **`persist`**：配置了状态持久化，将 `collapsed` 和 `watermarked` 状态存储到本地存储中，键名为 `gvb_admin_theme`。
  - **`state`**：定义了三个状态：
    - **`collapsed`**：布尔值，表示侧边栏是否折叠，初始值为 `false`。
    - **`watermarked`**：布尔值，表示是否显示水印，初始值为 `false`。
    - **`darkMode`**：布尔值，表示是否为暗色模式，初始值为 `isDark` 的值。
  - **`actions`**：定义了三个操作方法：
    - **`switchWatermark`**：切换水印状态。
    - **`switchCollapsed`**：切换侧边栏折叠状态。
    - **`switchDarkMode`**：切换暗色模式状态。

```javascript
import { defineStore } from 'pinia'
import { useDark } from '@vueuse/core'

const isDark = useDark()
export const useThemeStore = defineStore('theme-store', {
  persist: {
    key: 'gvb_admin_theme',
    paths: ['collapsed', 'watermarked'],
  },
  state: () => ({
    collapsed: false, // 侧边栏折叠
    watermarked: false, // 水印
    darkMode: isDark, // 黑暗模式
  }),
  actions: {
    switchWatermark() {
      this.watermarked = !this.watermarked
    },
    switchCollapsed() {
      this.collapsed = !this.collapsed
    },
    switchDarkMode() {
      this.darkMode = !this.darkMode
    },
  },
})

```





#### 7.2.4 src/store/modules/auth.js

这段代码定义了一个用于管理用户认证状态的 Pinia Store，主要功能包括：

1. **存储认证令牌（Token）**：通过 `token` 状态存储用户的认证信息，并支持持久化。
2. **设置 Token**：通过 `setToken` 方法设置用户的认证令牌。
3. **退出登录**：通过 `logout` 方法主动退出登录，调用后端接口并重置状态。
4. **强制下线**：通过 `forceOffline` 方法处理用户被强制下线的场景，重置状态。
5. **重定向到登录页**：通过 `toLogin` 方法将用户重定向到登录页面，并保留当前路由的查询参数。
6. **重置状态**：通过 `resetLoginState` 方法重置所有与登录相关的状态，包括其他 Store 的状态和路由配置。

这段代码的设计考虑了用户认证的完整流程，包括登录、登出和状态管理，同时通过与其他 Store 和路由的协同操作，确保了应用状态的一致性。

```javascript
import { unref } from 'vue'
import { defineStore } from 'pinia'
import { usePermissionStore, useTagStore, useUserStore } from '@/store'
import { resetRouter, router } from '@/router'
import api from '@/api'

export const useAuthStore = defineStore('auth', {
  persist: {
    key: 'gvb_admin_auth',
    paths: ['token'],
  },
  state: () => ({
    token: null,
  }),
  actions: {
    setToken(token) {
      this.token = token
    },
    toLogin() {
      const currentRoute = unref(router.currentRoute)
      router.replace({
        path: '/login',
        query: currentRoute.query,
      })
    },
    resetLoginState() {
      useUserStore().$reset()
      usePermissionStore().$reset()
      useTagStore().$reset()
      resetRouter()
      this.$reset()
    },
    /**
     * 主动退出登录
     */
    async logout() {
      await api.logout()
      this.resetLoginState()
      this.toLogin()
      window.$message.success('您已经退出登录！')
    },
    /**
     * TODO: 被强制退出
     */
    async forceOffline() {
      this.resetLoginState()
      this.toLogin()
      window.$message.error('您已经被强制下线！')
    },
  },
})
```



#### 7.2.5 src/store/modules/user.js

这段代码定义了一个用于管理用户信息的 Pinia Store，主要功能包括：

1. **存储用户信息**：通过 `userInfo` 状态存储用户的详细信息，包括 ID、昵称、头像、简介和个人网站。
2. **获取用户信息**：通过 `getUserInfo` 方法从后端接口获取用户信息，并更新状态。
3. **计算属性**：通过 `getters` 提供了方便访问用户信息的计算属性，例如 `userId`、`nickname`、`avatar` 等。
4. **图片路径处理**：通过 `convertImgUrl` 函数处理头像 URL，确保图片能够正确加载。

5. **扩展性**

- **`roles`**：虽然当前代码中 `roles` 字段被注释掉了，但可以通过后端接口返回用户角色信息，并在需要时启用该字段。
- **错误处理**：在 `getUserInfo` 方法中，通过 `try...catch` 捕获错误，确保在获取用户信息失败时能够正确处理。

这段代码的设计简洁明了，易于扩展和维护，能够很好地满足用户信息管理的需求。

```javascript
import { defineStore } from 'pinia'
import { convertImgUrl } from '@/utils'
import api from '@/api'

// 用户全局变量
export const useUserStore = defineStore('user', {
  state: () => ({
    userInfo: {
      id: null,
      nickname: '',
      avatar: '',
      intro: '',
      website: '',
      // roles: [], // TODO: 后端返回 roles
    },
  }),
  getters: {
    userId: state => state.userInfo.id,
    nickname: state => state.userInfo.nickname,
    intro: state => state.userInfo.intro,
    website: state => state.userInfo.website,
    avatar: state => convertImgUrl(state.userInfo.avatar),
    // roles: state => state.userInfo.roles,
  },
  actions: {
    async getUserInfo() {
      try {
        const resp = await api.getUserInfo()
        this.userInfo = resp.data
        return Promise.resolve(resp.data)
      }
      catch (err) {
        return Promise.reject(err)
      }
    },
  },
})

```



#### 7.2.6 src/store/index.js

这段代码的作用是：

1. **初始化 Pinia**：创建一个 Pinia 实例并将其挂载到 Vue 应用中。
2. **启用数据持久化**：通过 `pinia-plugin-persistedstate` 插件，使 Pinia Store 的状态能够在页面刷新后仍然保持，解决了数据丢失的问题。
3. **模块化 Store**：通过导入和导出各个模块化的 Store 文件，实现了状态管理的模块化，便于管理和维护。

```javascript
import { createPinia } from 'pinia'

// https://github.com/prazdevs/pinia-plugin-persistedstate
// pinia 数据持久化，解决刷新数据丢失的问题
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'

export function setupStore(app) {
  const pinia = createPinia()
  pinia.use(piniaPluginPersistedstate)
  app.use(pinia)
}

export * from './modules/permission'
export * from './modules/tag'
export * from './modules/theme'
export * from './modules/user'
export * from './modules/auth'

```

同时，及时更新 main.ts 中的方法：

```typescript
import { createApp } from 'vue'
import App from './App.vue'
import { setupRouter } from './router'
import { setupStore } from './store'

// unocss
import 'uno.css'
import '@unocss/reset/tailwind.css'

const app = createApp(App);
setupStore(app); // 优先级最高
await setupRouter(app);
app.mount('#app')
```



## 7.3 主页搭建 - utils / api  - 网络请求

#### 7.3.1 utils - http.js

首先先根据  7.2.4 src/store/modules/auth.js 配置好对应的仓库， 因为要用到其中的 token 来判断当前是否处于登陆状态（鉴权）





#### 7.3.1 utils - local.js



#### 7.3.1 utils - naiveTool.js



#### 7.3.1 utils - index.js



#### 7.3.5 api.js

建立 api 提供对应的后端方法接口， 如下：