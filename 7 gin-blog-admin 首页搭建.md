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





#### 7.2.2 src/store/modules/tag.js



#### 7.2.3 src/store/modules/theme.js



#### 7.2.4 src/store/modules/auth.js





#### 7.2.5 src/store/modules/theme.js



## 7.3 主页搭建 - utils / api  - 网络请求

#### 7.3.1 utils - http.js

首先先根据  7.2.4 src/store/modules/auth.js 配置好对应的仓库， 因为要用到其中的 token 来判断当前是否处于登陆状态（鉴权）





#### 7.3.1 utils - local.js



#### 7.3.1 utils - naiveTool.js



#### 7.3.1 utils - index.js



#### 7.3.5 api.js

建立 api 提供对应的后端方法接口， 如下：