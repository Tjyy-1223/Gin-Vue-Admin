# 第二章 gin-blog-front 静态页面搭建

## 2.1 消息提示组件

构建顶部中间的消息提示以及右上方的消息组件

**src/components/ui/UToast.vue** 构建该组件并应用于 App.vue

**src/App.vue**

```vue
<template>
  <!-- 顶部中间的消息提示 -->
  <UToast ref="messageRef" position="top" align="center" :timeout="3000" closeable />
  <!-- 右上方的消息通知 -->
  <UToast ref="notifyRef" position="top" align="right" :timeout="3000" closeable />
</template>

<script setup lang="ts">
import UToast from './components/ui/UToast.vue';
</script>

<style scoped lang="scss">
div {
  height: 300px;
}
</style>
```

**src/components/ui/UToast.vue 具体代码如下：**

```vue
<template>
    <Teleport to="body">
        <div v-show="flux.events.length" class="pointer-events-none fixed w-4/5 sm:w-[400px]" :class="{
            'left-1/2 -translate-x-1/2': align === 'center',
            'left-16': align === 'left',
            'right-16': align === 'right',
            'top-4': position === 'top',
            'bottom-6 ': position === 'bottom',
        }" :style="{ zIndex }">
            <!--  
            1. 指定渲染的根元素为 <ul> 标签
            2. enter-class -> enter-to-class: 进入时的初始状态，基于 position 的值来设置不同的动画
            3. leave-class -> leave-to-class 设置里开始的动画效果
            4. move-class 列表项移动时的过渡效果，使用平滑的过渡
        -->
            <TransitionGroup tag="ul" enter-active-class="transition ease-out duration-200"
                leave-active-class="transition ease-in duration-200 absolute w-full" :enter-class="position === 'bottom'
                    ? 'transform translate-y-3 opacity-0'
                    : 'transform -translate-y-3 opacity-0'" enter-to-class="transform translate-y-0 opacity-100"
                leave-class="transform translate-y-0 opacity-100" :leave-to-class="position === 'bottom'
                    ? 'transform translate-y-1/4 opacity-0'
                    : 'transform -translate-y-1/4 opacity-0'" move-class="ease-in-out duration-200"
                class="inline-block w-full">

                <!-- li 标签 -->
                <li v-for="event in flux.events" :key="event.id" :class="{
                    'pb-2': position === 'bottom',
                    'pt-2': position === 'top',
                }">
                    <!-- 这是一个插槽（slot）元素，用于为父组件提供一个可插入的区域 -->
                    <!-- 
                        这个模板是一个通用的通知组件（例如提示框），它：
                        根据 event.type 显示不同的图标和颜色（如成功、信息、警告、错误）。
                        显示 event.content 作为通知的文本内容。
                        提供一个关闭按钮，当点击时会调用 flux.remove(event) 来移除该通知。
                     -->
                    <slot :type="event.type" :content="event.content">
                        <div
                            class="pointer-events-auto w-full overflow-hidden rounded-lg bg-white ring-1 ring-black ring-opacity-5">
                            <div class="flex justify-between px-4 py-3">
                                <div class="flex items-center">
                                    <div class="mr-6 h-6 w-6" :class="{
                                        'i-mdi:check-circle text-green': event.type === 'success',
                                        'i-mdi:information-outline text-blue': event.type === 'info',
                                        'i-mdi:alert-outline text-yellow': event.type === 'warning',
                                        'i-mdi:alert-circle-outline text-red': event.type === 'error',
                                    }" />
                                    <div class="ml-1">
                                        <slot>
                                            <div> {{ event.content }} </div>
                                        </slot>
                                    </div>

                                </div>
                                <button v-if="closeable"
                                    class="i-mdi:close h-5 w-5 flex items-center justify-center rounded-full rounded-full p-1 font-bold text-gray-400 hover:text-gray-600"
                                    @click="flux.remove(event)" />
                            </div>
                        </div>
                    </slot>

                </li>
            </TransitionGroup>

        </div>
    </Teleport>
</template>

<script setup lang="ts">
import { defineProps, reactive } from 'vue';

const props = defineProps({
    show: { type: Boolean, default: false },
    position: {
        type: String,
        default: 'top',
        validator: (value: string) => {
            const positions: string[] = ['top', 'bottom'];
            return positions.includes(value);
        },
    },
    align: {
        type: String,
        default: 'center',
        validator: (value: string) => {
            const alignments: string[] = ['left', 'center', 'right'];
            return alignments.includes(value);
        },
    },
    timeout: { type: Number, default: 2500 },
    queue: { type: Boolean, default: true },
    zIndex: { type: Number, default: 100 },
    closeable: { type: Boolean, default: false },
});


interface FluxEvent {
    id: number;
    content: string;
    type: 'success' | 'info' | 'warning' | 'error';
}

interface Flux {
    events: FluxEvent[];
    success: (content: string) => void;
    info: (content: string) => void;
    warning: (content: string) => void;
    error: (content: string) => void;
    show: (type: 'success' | 'info' | 'warning' | 'error', content: string) => void;
    add: (type: 'success' | 'info' | 'warning' | 'error', content: string) => void;
    remove: (event: FluxEvent) => void;
}

// 这段代码定义了一个响应式对象 flux，它用于管理和控制一组通知事件（例如成功、信息、警告、错误通知）。在这个对象中，定义了几个方法来添加、显示和移除这些通知。
// 假设这是一个通知管理系统，你可以通过调用 flux.success('Message') 来显示一个成功类型的通知，flux.info('Info message') 来显示一个信息类型的通知，通知会在页面上显示一段时间后自动消失。
const flux = reactive<Flux>({
    events: [],

    // success、info、warning、error：这些是 flux 对象的快捷方法，用于向 flux.events 数组中添加不同类型的通知事件。
    // flux.success('message') 会等同于 flux.add('success', 'message')
    success: (content: string) => flux.add('success', content),
    info: (content: string) => flux.add('info', content),
    warning: (content: string) => flux.add('warning', content),
    error: (content: string) => flux.add('error', content),

    /**
     * @param {'success' | 'info' | 'warning' | 'error'} type
     * @param {string} content
     */
    show: (type: 'success' | 'info' | 'warning' | 'error', content: string) => flux.add(type, content),

    // 通过 setTimeout 延迟 100 毫秒后，创建一个新的通知事件对象，并将其添加到 flux.events 中。每个通知事件有一个唯一的 id，生成方式是通过 Date.now() 来保证唯一性。
    add: (type: 'success' | 'info' | 'warning' | 'error', content: string) => {
        if (!props.queue)
            flux.events = [];  // 这里的 props.queue 假设是一个外部的响应式属性，需确保 props 已被正确传递

        setTimeout(() => {
            const event: FluxEvent = { id: Date.now(), content, type };
            flux.events.push(event);
            setTimeout(() => flux.remove(event), props.timeout); // 假设 props.timeout 是预定义的超时时间
        }, 100);
    },

    // 这个方法接受一个通知事件对象 event，然后通过过滤掉该事件的 id，从 flux.events 数组中移除该事件。
    remove: (event: FluxEvent) => {
        flux.events = flux.events.filter((e) => e.id !== event.id);
    },
});

// defineExpose 将 flux 对象的 show、success、info、warning 和 error 方法暴露给外部组件。
defineExpose({
    show: flux.show,
    success: flux.success,
    info: flux.info,
    warning: flux.warning,
    error: flux.error,
})
</script>

<style scoped></style>
```

### 2.1.1 作用分析

这段 Vue 代码是一个实现通知管理系统的组件，提供了一个基于 Vue 3 的通知组件，允许显示不同类型的通知（例如成功、信息、警告、错误），并且可以通过不同的配置来控制显示的位置、关闭按钮以及超时时间等。

该组件的作用是：

- **显示通知**：向用户展示成功、信息、警告、错误等不同类型的消息。
- **动画效果**：通过 `TransitionGroup` 和 `Transition` 添加了进入和离开的动画效果。
- **自动关闭**：通知会在一定的超时时间后自动消失。
- **可关闭**：通知提供了一个关闭按钮，可以手动关闭通知。
- **队列管理**：通知按队列展示，可以设置是否启用队列模式。

这段代码实现了一个可配置的通知组件，可以用于显示不同类型的通知，并支持以下功能：

- **自动消失**：通过 `timeout` 控制通知的显示时长。
- **动态位置**：可以设置通知的位置（顶部或底部）以及对齐方式（左、中、右）。
- **通知队列**：支持队列模式，确保多个通知不会同时



### 2.1.2 样式如何起作用的

例如下面这种使用 class 指定样式的过程，是如何起作用的？

```
<div class="h-full w-full flex flex-col">
  <header>Header</header>
  <main>Main Content</main>
  <footer>Footer</footer>
</div>
```

这里是用了tailwindcss 属于颗粒化的 css，因此更`建议熟练使用 css 的人使用`，这是`新手陷阱`(会让新手成长更慢，毕竟接触到实际的css属性变少)

另外本篇也只介绍很多常见的，有些不常见或者不常用的，可以直接上style、css，或者查询 [tailwind 文档](https://link.juejin.cn?target=https%3A%2F%2Ftailwindcss.com%2Fdocs%2Finstallation)

此外由于已经使用了 颗粒化css 的 tailwindcss，其加上 style、css 可以完全替代 sass、less 这种，基本没有必要引入后两者了，当然真需要也不冲突

第一步：vscode 搜索 **tailwind** 插件(自己看着安装适合自己就行了)，直接搜索安装即可，不然会没有相关提示，这样开发效率反而可能会变低哈

第二步：开始安装相关库，在你的项目目录下运行以下命令：

```js
pnpm install -D tailwindcss@3 postcss autoprefixer 
```

这个命令会将 `tailwindcss`、`postcss` 和 `autoprefixer` 安装为开发依赖。

```
npx tailwindcss init -p
```

这会在项目根目录下创建一个 `tailwind.config.js` 文件，你可以根据需要修改它。

```
/** @type {import('tailwindcss').Config} */
export default {
  content: ["index.html", "./src/**/*.{html,js,ts,jsx,tsx,vue}"],
  theme: {
    extend: {},
  },
  plugins: [],
}
```

根目录创建 `postcss.config.js` 配置一下

```js
module.exports = {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
}
```

这会告诉 PostCSS 使用 `tailwindcss` 和 `autoprefixer` 插件来处理 CSS。

在项目中创建一个 `src/styles/tailwind.css` 文件，并添加以下内容：

```
/* src/styles/tailwind.css */
@tailwind base;
@tailwind components;
@tailwind utilities;
```

根目录创建 `tailwind.config.js`，或者执行 `npx tailwindcss init` 配置一下

在你的 `main.ts` 文件中，使用 `import` 语句引入你刚刚创建的 `tailwind.css` 文件。这样，整个项目就可以使用 Tailwind 的样式了。

修改你的 `main.ts` 文件，确保它类似于下面的内容：

```js
// main.ts

import { createApp } from 'vue'
import App from './App.vue'

// 引入 Tailwind CSS
import './styles/tailwind.css'

createApp(App).mount('#app')

```

然后可以在 webpack.config 的 plugins 中加入 trailwindcss 吧

先做一个简易的进度条吧，看看效果(这个不可以滑动哈)

```vue
<div
    className={`flex items-center w-[220px] h-[6px] bg-[#ffffff32] rounded-[2px] ${className}`}
>
    <div
        className={`bg-gradient-to-r from-[#0058BA] to-[#5FFFFB] h-[4px] overflow-hidden`}
        style={{
            flex: progress,
        }}
    />
    <div
        className="w-[5px] h-[8px] bg-[#a8deff] rounded-[6px]"
        style={{
            boxShadow: "0px 0px 8px 1px #A8DEFF",
        }}
    />
    <div
        className={`h-[4px] bg-[#ffffff32]`}
        style={{
            flex: 100 - progress,
        }}
    />
</div>
```

通过上面也可以简单看出，`tailwindcss` 使用简单，其就像设置 `style` 一样，`颗粒度非常小`，并且其还是class 形式，再配合 style，使用起来别提多舒服了，尤其是在 react 中可圈可点，并且还不是说不能写css了，想写仍然可以编写，并不冲突

此外，其又是缩写的class，所以整体css代码量会减少不少，也避免了起名的问题，缺点是对新手不友好，可能会忘记基础的css使用，但瑕不掩瑜，开发非常方便

因此也更建议 `css 基础扎实` 的人使用

**下面介绍一下常用的 tailwindcss 基础，很多都是类推，不介绍那么细，因此重要是再重复，不适合css新手看**

https://juejin.cn/post/7441229834116939803



## 2.2 AppHeader 顶部导航栏

```vue
<div class="h-full w-full flex flex-col">
    ...
</div>
```

`flex`：使容器成为 **弹性容器**。

`flex-col`：使子元素在容器内 **垂直方向排列**。

例如:

```vue
<div class="flex flex-col h-64 bg-gray-100">
  <div class="bg-red-500 h-1/4">Item 1</div>
  <div class="bg-blue-500 h-1/4">Item 2</div>
  <div class="bg-green-500 h-1/4">Item 3</div>
  <div class="bg-yellow-500 h-1/4">Item 4</div>
</div>
```

在这个例子中：

- 外部的 `<div>` 使用了 `flex flex-col`，这意味着它是一个 **弹性容器**，并且它的子元素（`Item 1`、`Item 2` 等）会按 **垂直方向** 从上到下排列。
- 每个子元素的高度 (`h-1/4`) 都占父容器的 1/4，所以这些项将均匀分布在容器内。

### 2.2.1 全局状态 pinia

Pinia 提供了一个中心化的地方来管理整个应用的状态，避免了多组件间状态传递的混乱。在一个大型应用中，各种状态（例如用户信息、购物车、主题设置等）可能会在多个组件中被用到，使用 Pinia 可以方便地在这些组件之间共享这些状态。

首先进行相关安装：

```
pnpm add pinia
pnpm add pinia-plugin-persistedstate
```

src/store/index.js

```javascript
import { createPinia } from 'pinia'

// https://github.com/prazdevs/pinia-plugin-persistedstate
// pinia 数据持久化，解决刷新数据丢失的问题
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'

export const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)

// export * from './app'
export * from './user'
```

在 main.ts 中进行相关配置

```
import { pinia } from './store'

let app = createApp(App);
app.use(pinia);
```

然后， 我们在 src/store/user.ts 构建一个用户相关的全局信息

```typescript
import { defineStore } from 'pinia';

// 自定义 state 的类型
interface UserState {
  userInfo: {
    id: string;
    nickname: string;
    avatar: string;
    website: string;
    intro: string;
    email: string;
    articleLikeSet: any[];
    commentLikeSet: any[];
  };
  token: string | null;
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    userInfo: {
      id: '',
      nickname: '',
      avatar: 'https://example.com/default-avatar.png',
      website: '',
      intro: '',
      email: '',
      articleLikeSet: [],
      commentLikeSet: [],
    },
    token: null,
  }),
  getters: {
    userId: state => state.userInfo.id ?? '',
    nickname: state => state.userInfo.nickname ?? '',
    avatar: state => state.userInfo.avatar ?? 'https://example.com/default-avatar.png',
    website: state => state.userInfo.website ?? '',
    intro: state => state.userInfo.intro ?? '',
    email: state => state.userInfo.email ?? '',
    articleLikeSet: state => state.userInfo.articleLikeSet || [],
    commentLikeSet: state => state.userInfo.commentLikeSet || [],
  },
  actions: {
    // 如果有 actions，可以在这里添加
  },
  persist: {
    key: 'gvb_blog_user',
    paths: ['token'], // 确保 token 是直接在 state 的根路径下
  }
});
```

**同理，简单配置一个：src/store/app.ts**

```typescript
import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', {
    state: () => ({
        searchFlag: false,
        loginFlag: false,
        registerFlag: false,
        collapsed: false, // 侧边栏折叠（移动端）

        page_list: [], // 页面数据
        // TODO: 优化
        blogInfo: {
            article_count: 0,
            category_count: 0,
            tag_count: 0,
            view_count: 0,
            user_count: 0,
        },
        blog_config: {
            website_name: 'CodingHome',
            website_author: 'Tjyy',
            website_intro: 'coding is coding',
            website_avatar: '',
        },
    }),
    getters: {
        isMobile: () => !!navigator.userAgent.match(/(phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone)/i),
        articleCount: state => state.blogInfo.article_count ?? 0,
        categoryCount: state => state.blogInfo.category_count ?? 0,
        tagCount: state => state.blogInfo.tag_count ?? 0,
        viewCount: state => state.blogInfo.view_count ?? 0,
        pageList: state => state.page_list ?? [],
        blogConfig: state => state.blog_config,
    },
    actions: {
        setCollapsed(flag:any) { this.collapsed = flag },
        setLoginFlag(flag:any) { this.loginFlag = flag },
        setRegisterFlag(flag:any) { this.registerFlag = flag },
        setSearchFlag(flag:any) { this.searchFlag = flag },
    },
})

```

**至此，我们可以继续补充 src/components/layout/AppHeader.vue**



### 2.2.2 AppHeader.vue

为了使用 RouterLink，这里首先要安装并配置 vue-router

```
pnpm add vue-router
```

在项目的入口文件（通常是 `main.js` 或 `main.ts`）中，确保正确引入并注册了 Vue Router：

```javascript
import { createApp } from 'vue';
import App from './App.vue';
import router from './router'; // 确保路由配置文件已正确导入

const app = createApp(App);
app.use(router); // 注册路由
app.mount('#app');
```

确保路由配置文件（通常是 `router/index.js` 或 `router/index.ts`）中正确配置了路由规则。例如：

```javascript
import { createRouter, createWebHistory } from 'vue-router';
import Home from '../views/Home.vue';
import About from '../views/About.vue';

const routes = [
  { path: '/', component: Home },
  { path: '/about', component: About },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
```

**然后继续构建 AppHeader.vue**

```vue
<template>
    <!-- 移动端顶部导航栏 -->
    <Transition name="slide-face" appear>
        <div v-if="barShow" :class="navClass"
            class="fixed inset-x-0 top-0 z-11 h-[60px] flex items-center justify-between px-4 py-2 lg:hidden">
            <!-- 左上角标题 -->
            <!-- <RouterLink> 是 Vue Router 提供的组件，用于在 Vue 应用中进行页面导航。在这里，点击该链接会导航到首页（/）。 -->
            <RouterLink to="/" class="text-[18px] font-bold">
                {{ appStore.blogConfig.website_author }}
            </RouterLink>

            <!-- 右上角图标 -->
            <div class="flex items-center gap-2 text-2xl">
                <button @click="appStore.setSearchFlag(true)">
                    <Icon icon="ic:round-search" />
                </button>
                <button @click="appStore.setCollapsed(true)">
                    <Icon icon="ic:sharp-menu" />
                </button>
            </div>
        </div>
    </Transition>

    <!-- 电脑端顶部导航栏 -->
    <Transition name="slide-fade" appear>
        <div v-if="barShow" :class="navClass" class="fixed inset-x-0 top-0 z-11 hidden h-[60px] lg:block">
            <div class="h-full flex items-center justify-between px-9">
                <!-- 左上角标题 -->
                <RouterLink to="/" class="text-xl font-bold">
                    {{ appStore.blogConfig.website_author }}
                </RouterLink>
                <!-- 右上角菜单 -->
                <div class="flex items-center space-x-4">
                    <!-- 搜索 -->
                    <div class="menus-item">
                        <a class="menu-btn flex items-center" @click="appStore.setSearchFlag(true)">
                            <Icon icon="mdi:magnify" class="text-xl" />
                            <span class="ml-1"> 搜索 </span>
                        </a>
                    </div>

                    <!-- 多个菜单，根据 menuOptions 配置 -->
                    <div v-for="item of menuOptions" :key="item.text" class="menus-item">
                        <!-- 不包含子菜单 -->
                        <RouterLink v-if="!item.subMenu" :to="item.path" class="menu-btn flex items-center">
                            <Icon :icon="item.icon" class="text-xl" />
                            <span class="ml-1"> {{ item.text }} </span>
                        </RouterLink>
                        <!-- 包含子菜单 -->
                        <div v-else class="menu-btn">
                            <div class="flex items-center">
                                <Icon :icon="item.icon" class="text-xl" />
                                <span class="mx-1"> {{ item.text }} </span>
                                <Icon icon="ep:arrow-down-bold" class="text-xl" />
                            </div>
                            <ul class="menus-submenu">
                                <RouterLink v-for="sub of item.subMenu" :key="sub.text" :to="sub.path">
                                    <div class="flex items-center">
                                        <Icon :icon="sub.icon" class="text-xl" />
                                        <span class="ml-1"> {{ sub.text }} </span>
                                    </div>
                                </RouterLink>
                            </ul>
                        </div>
                    </div>


                    <!-- 登陆设计 -->
                    <div class="menus-item">
                        <!-- 登陆图标 -->
                        <a v-if="!userStore.userId" class="menu-btn" @click="appStore.setLoginFlag(true)">
                            <div class="flex items-center">
                                <Icon icon="ph:user-bold" class="text-xl" />
                                <span class="ml-1"> 登录 </span>
                            </div>
                        </a>
                        <!-- 如果当前是登陆状态，则展示下面的部分 -->
                        <template v-else>
                            <img :src="convertImgUrl(userStore.avatar)" class="w-8 cursor-pointer rounded-full">
                            <ul class="menus-submenu">
                                <RouterLink to="/user">
                                    <div class="flex items-center">
                                        <Icon icon="mdi:account-circle" class="mr-1 text-xl" /> 个人中心
                                    </div>
                                </RouterLink>
                                <a @click="logout">
                                    <div class="flex items-center">
                                        <Icon icon="mdi:logout" class="mr-1 text-xl" /> 退出登录
                                    </div>
                                </a>
                            </ul>
                        </template>
                    </div>
                </div>
            </div>
        </div>
    </Transition>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { ref } from 'vue';
import { Icon } from '@iconify/vue'
import { useAppStore, useUserStore } from '@/store'
import { useWindowScroll, watchThrottled } from '@vueuse/core'
import { convertImgUrl } from '@/utils'
const navClass = ref('nav')
const barShow = ref(true)
const appStore = useAppStore()
const userStore = useUserStore()
const router = useRouter()
const route = useRoute()

const menuOptions = [
    { text: '首页', icon: 'mdi:home', path: '/' },
    {
        text: '发现',
        icon: 'mdi:apple-safari',
        subMenu: [
            { text: '归档', icon: 'mdi:archive', path: '/archives' },
            { text: '分类', icon: 'mdi:menu', path: '/categories' },
            { text: '标签', icon: 'mdi:tag', path: '/tags' },
        ],
    },
    {
        text: '娱乐',
        icon: 'mdi:gamepad-circle',
        subMenu: [
            { text: '相册', icon: 'mdi:view-gallery', path: '/albums' },
        ],
    },
    { text: '友链', icon: 'mdi:vector-link', path: '/links' },
    { text: '关于', icon: 'mdi:information-outline', path: '/about' },
    { text: '留言', icon: 'mdi:forum', path: '/message' },
]


// * 监听 y 效果比添加 scroll 监听器效果更好
// * 这段代码的作用是实现一个基于滚动位置的响应式导航栏效果，通常用于网页或应用的顶部导航栏。它通过监听窗口的滚动位置（y 值），并根据滚动方向和滚动距离动态调整导航栏的显示状态和样式。
const { y } = useWindowScroll()
const preY = ref(0) // 记录上一次的 y 滚动距离
watchThrottled(y, () => {
    if (Math.abs(preY.value - y.value) >= 50) { // 小幅度滚动不进行操作
        barShow.value = (y.value < preY.value)
        navClass.value = (y.value > 60) ? 'nav-fixed' : 'nav'
        preY.value = y.value
    }
}, { throttle: 100 })

// 退出登录
async function logout() {
    //   await userStore.logout()
    if (route.name === 'User') {
        router.push('/') // 退出后跳转到首页
    }
    window.$notify.success('退出登录成功!')
}
</script>

<style scoped lang="scss">
.nav {
    transition: all 0.8s;
    //   color: #fff;
    background: rgba(0, 0, 0, 0) !important;
}

.nav-fixed {
    transition: all 0.8s;
    color: #000;
    background: rgba(255, 255, 255, 0.8) !important;
    box-shadow: 0 5px 6px -5px rgba(133, 133, 133, 0.6);

    & .menu-btn:hover {
        color: #49b1f5 !important;
    }
}

.menus-item {
    position: relative;
    display: inline-block;

    // margin: 0 0 0 1rem;
    a {
        transition: all 0.2s;
    }

    a::after {
        position: absolute;
        bottom: -5px;
        left: 0;
        z-index: -1;
        width: 0;
        height: 3px;
        background-color: #80c8f8;
        content: "";
        transition: all 0.3s ease-in-out;
    }

    .menu-btn {
        cursor: pointer;

        &:hover::after {
            width: 100%;
        }
    }
}

// 当鼠标悬停在 .menus-item 元素上时，会改变其内部的 .menus-submenu 元素的显示状态
.menus-item:hover .menus-submenu {
    display: block;
}


.menus-submenu {
    position: absolute;
    display: none; // 默认情况下，.menus-submenu 是隐藏的，不会显示在页面上。
    right: 0;
    width: max-content;
    margin-top: 8px;
    box-shadow: 0 5px 20px -4px rgba(0, 0, 0, 0.5);
    background-color: #fff;
    animation: submenu 0.3s 0.1s ease both;

    &::before {
        position: absolute;
        top: -8px;
        left: 0;
        width: 100%;
        height: 20px;
        content: "";
    }

    a {
        line-height: 2;
        color: #4c4948 !important;
        text-shadow: none;
        display: block;
        padding: 6px 14px;
    }

    a:hover {
        background: #4ab1f4;
    }
}

@keyframes submenu {
    0% {
        opacity: 0;
        filter: alpha(opacity=0);
        transform: translateY(10px);
    }

    100% {
        opacity: 1;
        filter: none;
        transform: translateY(0);
    }
}
</style>
```



### 2.2.3 效果展示

![image-20250201141646625](./assets/image-20250201141646625.png)



