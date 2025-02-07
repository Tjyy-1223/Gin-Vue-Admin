# 第九章 gin-blog-admin login-home-article-auth

## 9.1 App.vue

这段代码展示了一个 Vue 3 组件，它通过 `Naive UI` 配置全局主题和本地化，支持暗黑模式和代码高亮显示，并且使用 Vue Router 动态渲染路由组件。代码简洁、功能清晰，适合用于需要根据用户主题偏好调整界面和高亮显示代码的应用。

**全局配置**：

- 通过 `NConfigProvider` 组件全局配置了 Naive UI 的主题、语言、日期格式以及高亮库，确保整个应用的主题和本地化设置一致。

**动态路由渲染**：

- 使用 `RouterView` 和 Vue 动态组件机制，能够根据当前路由渲染对应的页面组件。`RouterView` 和 `v-slot` 的使用使得动态渲染更加灵活。

**主题切换**：

- 根据 `themeStore.darkMode` 的状态来动态切换 Naive UI 的主题，支持深色模式（dark mode）。

**代码高亮**：

- 配置了 `highlight.js` 用于代码高亮显示，特别是 JSON 格式的代码。

```vue
<template>
  <NConfigProvider class="h-full w-full" :theme="themeStore.darkMode ? darkTheme : undefined"
    :theme-overrides="themes.naiveThemeOverrides" :locale="zhCN" :date-locale="dateZhCN" :hljs="hljs">
    <RouterView v-slot="{ Component }">
      <component :is="Component" />
    </RouterView>
  </NConfigProvider>
</template>


<script setup>
import { onMounted } from 'vue'
import { NConfigProvider, darkTheme, dateZhCN, zhCN } from 'naive-ui'
import hljs from 'highlight.js/lib/core'
import json from 'highlight.js/lib/languages/json'

import { useAuthStore, useThemeStore } from '@/store'
import themes from '@/assets/themes'
import api from '@/api'

hljs.registerLanguage('json', json)
const themeStore = useThemeStore()
</script>

<style lang="scss" scoped></style>
```

可以首先生成一个src/views/error-page/404.vue

```vue
<template>
    <div>

    </div>
</template>

<script setup>

</script>

<style lang="scss" scoped>

</style>
```

可以访问 http://localhost:3333/404 暂时看到一个空白页面





## 9.2 Login 登陆界面

### 9.2.1 index.html

这段代码是一个 Vue 应用的 `index.html` 文件，它是 Vue 项目的入口文件之一，主要作用是设置页面的基本结构、引入必要的资源，并启动 Vue 应用。下面我们逐步分析每个部分的作用。

1. **`<head>` 部分**

```html
<head>
  <meta charset="UTF-8" />
  <meta http-equiv="Expires" content="0" />
  <meta http-equiv="Pragma" content="no-cache" />
  <meta http-equiv="Cache-control" content="no-cache" />
  <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <link rel="icon" type="image/svg+xml" href="/favicon.svg" />
  <link rel="stylesheet" href="/resource/loading.css" />
  <title> Gin Vue Blog </title>
</head>
```

- **`<meta charset="UTF-8" />`**：设置文档字符编码为 UTF-8，支持多语言字符集。
- **`<meta http-equiv="Expires" content="0" />`**, **`<meta http-equiv="Pragma" content="no-cache" />`**, **`<meta http-equiv="Cache-control" content="no-cache" />`**：这些标签禁用了页面的缓存，确保每次加载页面时都从服务器获取最新的内容。
- **`<meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1" />`**：确保浏览器使用最新版本的渲染引擎，兼容老旧版本的 IE 浏览器。
- **`<meta name="viewport" content="width=device-width, initial-scale=1.0" />`**：设置视口的大小，确保页面在不同设备上适配，并且支持移动端响应式设计。
- **`<link rel="icon" type="image/svg+xml" href="/favicon.svg" />`**：设置页面的 favicon 图标，通常用于浏览器标签栏显示。
- **`<link rel="stylesheet" href="/resource/loading.css" />`**：引入了外部 CSS 文件，用于加载页面时的动画效果。
- **`<title> Gin Vue Blog </title>`**：设置页面的标题。

2. **`<style>` 部分**

```html
<style>
  html,
  body {
    width: 100%;
    height: 100%;
    overflow: hidden;
  }
</style>
```

- **`width: 100%; height: 100%;`**：使 `<html>` 和 `<body>` 元素的宽高占满整个屏幕。
- **`overflow: hidden;`**：禁用页面的滚动条，确保页面不会滚动，通常在加载动画期间使用。

3. **`<body>` 部分**

```html
<body class="dark:text-#e9e9e9 auto-bg">
  <div id="app" class="w-full h-full">
    <!-- loading -->
    <div class="loading-container">
      <img class="loading-logo" src="/resource/logo.svg" alt="logo" />
      <div class="loading-spin__container">
        <div class="loading-spin">
          <div class="left-0 top-0 loading-spin-item"></div>
          <div class="left-0 bottom-0 loading-spin-item loading-delay-500"></div>
          <div class="right-0 top-0 loading-spin-item loading-delay-1000"></div>
          <div class="right-0 bottom-0 loading-spin-item loading-delay-1500"></div>
        </div>
      </div>
      <div class="loading-title"> Gin Vue Blog </div>
    </div>
    <script src="/resource/loading.js"></script>
  </div>
  <script type="module" src="/src/main.ts"></script>
</body>
```

- **`<body class="dark:text-#e9e9e9 auto-bg">`**：`class="dark:text-#e9e9e9 auto-bg"` 可能是用于支持暗黑模式的样式，`text-#e9e9e9` 设置文字颜色，`auto-bg` 可能是设置背景样式的类名（具体的效果依赖于样式表）。

- `id="app"`：Vue 应用的挂载点，Vue 实例将会挂载到此 DOM 元素上。
- `class="w-full h-full"`：确保 `#app` 元素的宽度和高度占满整个页面。

- **`loading-container`**：包含了整个加载动画的容器，显示加载界面，直到 Vue 应用加载完毕。
- **`loading-logo`**：显示应用的 logo。
- **`loading-spin__container`** 和 **`loading-spin`**：这些元素包含了旋转的加载动画，具体样式和动画效果在 `loading.css` 中定义。
- **`loading-title`**：显示应用标题文本 `Gin Vue Blog`。

4. **引入 Vue 应用脚本**

```
<script type="module" src="/src/main.ts"></script>
```

- **`<script type="module" src="/src/main.ts"></script>`**：引入 Vue 应用的入口文件 `main.ts`，这是 Vue 项目的启动脚本。在这个文件中，会创建 Vue 实例、配置路由、全局状态管理等。`type="module"` 表示该脚本是一个模块化的 JavaScript 文件，采用 ES6 模块导入方式。

------

这个 `index.html` 文件为 Vue 应用提供了一个初始化界面，具体的流程如下：

1. 页面加载时显示一个加载动画（包括 logo 和旋转的加载图标）。
2. 加载动画和静态内容通过 `loading.css` 和 `loading.js` 控制。
3. Vue 应用在 `main.ts` 文件中初始化，并挂载到 `#app` 元素上。
4. 一旦 Vue 应用加载完成，加载动画会被隐藏，实际的应用界面会呈现给用户。

**这段代码为单页面应用（SPA）提供了一个良好的用户体验，确保在 Vue 应用初始化和渲染过程中，用户不会看到空白页面，而是显示一个友好的加载界面。**

```html
<!DOCTYPE html>
<html lang="cn">

<head>
  <meta charset="UTF-8" />
  <meta http-equiv="Expires" content="0" />
  <meta http-equiv="Pragma" content="no-cache" />
  <meta http-equiv="Cache-control" content="no-cache" />
  <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />

  <!-- <meta http-equiv="Content-Security-Policy" content="upgrade-insecure-requests"> -->

  <link rel="icon" type="image/svg+xml" href="/favicon.svg" />
  <link rel="stylesheet" href="/resource/loading.css" />
  <title> Gin Vue Blog </title>
  <!-- 腾讯验证码 -->
  <!-- <script src="https://ssl.captcha.qq.com/TCaptcha.js"></script> -->
  <style>
    html,
    body {
      width: 100%;
      height: 100%;
      overflow: hidden;
    }
  </style>
</head>

<body class="dark:text-#e9e9e9 auto-bg">
  <div id="app" class="w-full h-full">
    <!-- loading -->
    <div class="loading-container">
      <img class="loading-logo" src="/resource/logo.svg" alt="logo" />
      <div class="loading-spin__container">
        <div class="loading-spin">
          <div class="left-0 top-0 loading-spin-item"></div>
          <div class="left-0 bottom-0 loading-spin-item loading-delay-500"></div>
          <div class="right-0 top-0 loading-spin-item loading-delay-1000"></div>
          <div class="right-0 bottom-0 loading-spin-item loading-delay-1500"></div>
        </div>
      </div>
      <div class="loading-title"> Gin Vue Blog </div>
    </div>
    <script src="/resource/loading.js"></script>
  </div>
  <script type="module" src="/src/main.ts"></script>
</body>

</html>
```



### 9.2.2 AddPage.vue

**src/components/common/TheFooter.vue**

这段代码的目的是渲染一个简单的页脚组件，包含版权信息和指向 GitHub 的个人主页链接。组件使用了 Tailwind CSS 来进行样式布局，并动态显示当前年份。主要功能和设计如下：

1. **版权信息**：显示当前年份。
2. **GitHub 链接**：提供一个指向 GitHub 个人主页的链接，点击后会在新标签页中打开。
3. **使用了 Flexbox 布局**：实现垂直和水平居中。
4. **没有动态功能**：目前没有 JavaScript 逻辑，仅渲染静态内容。
5. **样式**：使用了 Tailwind CSS 来快速实现页面布局，但 `hover` 样式的写法存在一个小错误，需要调整为 `hover:underline` 和 `hover:text-primary`。

```vue
<template>
    <footer class="flex flex-col items-center justify-center text-gray-700">
        <p>
            Copyright©{{ new Date().getFullYear() }}
            <a href="https://github.com/Tjyy-1223" target="__blank" hover="decoration-underline color-primary">
                tjyy
            </a>
        </p>
    </footer>
</template>

<script setup>

</script>

<style lang="scss" scoped></style>
```

**src/components/common/AppPage.vue**

这个 Vue 3 组件提供了以下功能：

1. **过渡动画**：组件在切换时使用了 `fade-slide` 过渡动画，并且动画模式设置为 `out-in`，即移除旧元素后再添加新元素。
2. **页面布局**：
   - 使用 Flexbox 布局将内容垂直排列。
   - 支持暗黑模式背景色。
3. **插槽支持**：通过 `<slot />` 插入父组件传入的内容，提供灵活性。
4. **页脚**：通过 `showFooter` 属性控制是否显示 `TheFooter` 页脚组件，默认不显示。
5. **返回顶部按钮**：使用 `NBackTop` 组件来提供返回顶部功能，`bottom="20"` 设置按钮距离底部 20px。

```vue
<template>
    <Transition name="fade-slide" mode="out-in" appear>
        <section class="cus-scroll-y h-full w-full flex flex-col bg-[#f5f6fb] p-4 dark:bg-[#121212]">
            <slot />
            <TheFooter v-if="showFooter" class="mt-5" />
            <NBackTop :bottom="20" />
        </section>
    </Transition>
</template>


<script setup>
import { NBackTop } from 'naive-ui'
import TheFooter from '@/components/common/TheFooter.vue'

defineProps({
    showFooter: { type: Boolean, default: false },
})
</script>

<style lang="scss" scoped>
/* transition fade-slide */
.fade-slide-leave-active,
.fade-slide-enter-active {
    transition: all 0.3s;
}

.fade-slide-enter-from {
    opacity: 0;
    transform: translateX(-30px);
}

.fade-slide-leave-to {
    opacity: 0;
    transform: translateX(30px);
}

/* 自定义滚动条样式 */
.cus-scroll {
    overflow: auto;

    &::-webkit-scrollbar {
        width: 8px;
        height: 8px;
    }
}

.cus-scroll-x {
    overflow-x: auto;

    &::-webkit-scrollbar {
        width: 0;
        height: 8px;
    }
}

.cus-scroll-y {
    overflow-y: auto;

    &::-webkit-scrollbar {
        width: 8px;
        height: 0;
    }
}

.cus-scroll,
.cus-scroll-x,
.cus-scroll-y {
    &::-webkit-scrollbar-thumb {
        background-color: transparent;
        border-radius: 4px;
    }

    &:hover {
        &::-webkit-scrollbar-thumb {
            background: #bfbfbf;
        }

        &::-webkit-scrollbar-thumb:hover {
            background: var(--primary-color);
        }
    }
}
</style>
```



### 9.2.3 Login.vue

<img src="./assets/image-20250207141312354.png" alt="image-20250207141312354" style="zoom:80%;" />

**src/views/Login.vue**

**页面结构**：整个页面布局简单清晰，左侧为登录框的 `Banner` 图，右侧为表单，用户可以输入用户名、密码并选择是否记住登录信息。

**逻辑实现**：组件通过 `vuex` 的 `store` 和 `localStorage` 来管理用户的登录状态，实现了动态路由和用户信息的获取。

**样式**：通过 Tailwind CSS 实现响应式布局和现代化的用户界面，背景图和登录框的设计简洁美观。

```vue
<template>
    <!-- FIXME: 使用 style="background-image: url(/image/login_bg.webp);" 不生效, 需要写到 style 里的 class 中 -->
    <!-- AppPage 是一个容器，可能是用于整个页面的布局，使用了背景图片和背景覆盖的类 -->
    <AppPage class="backgroundImg bg-cover">
        <!-- 这里设置了一个白色半透明的登录框，并将其居中显示 -->
        <div style="transform: translateY(25px)"
            class="m-auto max-w-[700px] min-w-[345px] flex items-center justify-center rounded-2 bg-white bg-opacity-60 p-4 shadow">
            <!-- 登录框左侧的 Banner 图，仅在 md 屏幕以上显示 -->
            <div class="hidden w-[380px] px-5 py-9 md:block">
                <img src="/image/login_banner.webp" class="w-full" alt="login_banner">
            </div>

            <!-- 登录框右侧的内容区域，用于显示登录表单 -->
            <div class="w-[320px] flex flex-col px-4 py-9 space-y-5.5">
                <!-- 登录框顶部的 logo 和标题 -->
                <h5 class="flex items-center justify-center text-2xl text-gray font-normal">
                    <!-- Logo 图标 -->
                    <img src="/image/logo.svg" alt="logo" class="mr-2 h-[50px] w-[50px]">
                    <!-- 页面标题，绑定了 title 数据 -->
                    <span> {{ title }} </span>
                </h5>

                <!-- 用户名输入框，使用 Naive UI 的 NInput 组件 -->
                <!-- v-model:value 双向绑定用户名，设置最大长度为 20 -->
                <NInput v-model:value="loginForm.username" class="h-[50px] items-center pl-2" autofocus
                    placeholder="test@qq.com" :maxlength="20" />

                <!-- 密码输入框，使用 Naive UI 的 NInput 组件 -->
                <!-- 显示密码的切换功能通过 show-password-on 属性实现，设置最大长度为 20 -->
                <!-- 按下回车键时触发 handleLogin 方法 -->
                <NInput v-model:value="loginForm.password" class="h-[50px] items-center pl-2" type="password"
                    show-password-on="mousedown" placeholder="11111" :maxlength="20" @keydown.enter="handleLogin" />

                <!-- 记住我复选框，双向绑定 isRemember 状态 -->
                <NCheckbox :checked="isRemember" label="记住我" :on-update:checked="(val) => (isRemember = val)" />

                <!-- 登录按钮，绑定 loading 状态，点击时触发 handleLogin 方法 -->
                <NButton class="h-[50px] w-full rounded-5" type="primary" :loading="loading" @click="handleLogin">
                    登录
                </NButton>
            </div>
        </div>
    </AppPage>
</template>



<script setup>
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useStorage } from '@vueuse/core'
import { NButton, NCheckbox, NInput } from 'naive-ui'

import AppPage from '@/components/common/AppPage.vue'

import { addDynamicRoutes } from '@/router' // 用于添加动态路由
import { getLocal, removeLocal, setLocal } from '@/utils' // 本地存储操作工具
import { useAuthStore, useUserStore } from '@/store' // 状态管理中的用户和认证 store
import api from '@/api' // API 请求封装

// 获取环境变量中的网站标题
const title = import.meta.env.VITE_TITLE // 环境变量中读取

// 从 store 中获取用户认证和信息
const userStore = useUserStore()
const authStore = useAuthStore()

// 创建 router 和 route 实例，用于页面跳转和获取 URL 查询参数
const router = useRouter()
const { query } = useRoute()

// 初始化登录表单的用户名和密码
const loginForm = reactive({
    username: 'guest', // 默认用户名为 'guest'
    password: '123456', // 默认密码为 '123456'
})

// 在组件初始化时，尝试从本地存储中获取保存的登录信息
initLoginInfo()

// 从 localStorage 中获取记住的用户名和密码
function initLoginInfo() {
    const localLoginInfo = getLocal('loginInfo') // 获取存储在 localStorage 中的登录信息
    if (localLoginInfo) {
        // 如果存在，则将用户名和密码填充到登录表单中
        loginForm.username = localLoginInfo.username
        loginForm.password = localLoginInfo.password
    }
}

// 使用 vueuse 提供的 useStorage 创建响应式的本地存储变量
const isRemember = useStorage('isRemember', false) // 用于判断是否记住密码
const loading = ref(false) // 控制登录按钮的加载状态

// 处理登录的函数
async function handleLogin() {
    const { username, password } = loginForm // 获取表单中的用户名和密码
    if (!username || !password) {
        // 如果用户名或密码为空，弹出警告提示
        $message.warning('请输入用户名和密码')
        return
    }

    // 登录操作的具体实现
    const doLogin = async (username, password) => {
        loading.value = true // 设置加载状态为 true，表示正在执行登录请求

        // 调用登录接口进行用户身份验证
        try {
            const resp = await api.login({ username, password }) // 向后端发送登录请求
            authStore.setToken(resp.data.token) // 登录成功后，将获取到的 token 存储到 authStore 中

            // 获取用户信息，并添加动态路由
            await userStore.getUserInfo()
            await addDynamicRoutes()

            // 根据是否勾选“记住我”来保存或删除用户名和密码
            isRemember ? setLocal('loginInfo', { username, password }) : removeLocal('loginInfo')

            // 弹出登录成功的提示消息
            $message.success('登录成功')

            // 根据 URL 中的 redirect 查询参数进行跳转
            if (query.redirect) {
                const path = query.redirect
                Reflect.deleteProperty(query, 'redirect') // 删除 query 对象中的 redirect 属性
                router.push({ path, query }) // 跳转到原本的目标路径
            }
            else {
                // 如果没有 redirect 参数，则跳转到首页
                router.push('/')
            }
        }
        finally {
            loading.value = false // 无论登录成功或失败，都结束加载状态
        }
    }

    // 执行登录操作
    doLogin(username, password)
}
</script>


<style lang="scss" scoped>
.backgroundImg {
    background-image: url(/image/login_bg.webp);
}
</style>
```







## 9.3 Home 界面搭建

### 9.3.1 home/index.vue

![image-20250207153318162](./assets/image-20250207153318162.png)

**主要代码位于 src/views/home**

这个 Vue 3 组件展示了用户的个人信息、统计数据以及项目卡片。其主要功能包括：

1. **用户信息**：显示头像、昵称和一句话。通过 `useUserStore` 获取用户信息，动态加载一句话。
2. **统计信息**：通过 API 获取并显示访问量、用户量、文章量和留言量。
3. **项目卡片**：显示一组项目卡片，其中包含项目标题和描述。

该组件结合了 `naive-ui` 组件库的使用，展示了一种现代的网页布局和交互方式。

```vue
<template>
    <!-- 外层容器是一个 AppPage 组件，包裹整个页面 -->
    <AppPage>
        <div class="flex-1">

            <!-- 用户信息卡片 -->
            <NCard>
                <div class="flex items-center">
                    <!-- 用户头像，圆形，大小为 60px -->
                    <NAvatar round :size="60" :src="avatar" />
                    <div class="ml-5">
                        <!-- 显示用户昵称 -->
                        <p> Hello, {{ nickname }} </p>
                        <!-- 使用线性渐变的文本显示一句话 -->
                        <NGradientText class="mt-1 op-60"
                            gradient="linear-gradient(90deg, red 0%, green 50%, blue 100%)">
                            {{ sentence }}
                        </NGradientText>
                    </div>
                    <div class="ml-auto flex items-center">
                        <!-- 显示 GitHub 项目的 Stars 数量 -->
                        <NStatistic label="Stars" class="w-[80px]">
                            <a href="https://github.com/szluyu99/gin-vue-blog" target="_blank">
                                <img alt="stars" src="https://badgen.net/github/stars/szluyu99/gin-vue-blog">
                            </a>
                        </NStatistic>
                        <!-- 显示 GitHub 项目的 Forks 数量 -->
                        <NStatistic label="Forks" class="ml-10 w-[100px]">
                            <a href="https://github.com/szluyu99/gin-vue-blog" target="_blank">
                                <img alt="forks" src="https://badgen.net/github/forks/szluyu99/gin-vue-blog">
                            </a>
                        </NStatistic>
                    </div>
                </div>
            </NCard>

            <!-- 首页统计信息展示，使用 NGrid 布局 -->
            <NGrid class="mt-4" x-gap="12" :cols="4">
                <!-- v-for 遍历数组渲染每个统计项 -->
                <template v-for="item of [
                    { icon: 'i-fa6-solid:users', color: 'text-[#40C9C6]', label: '访问量', key: 'view_count' },
                    { icon: 'i-heroicons:users-solid', color: 'text-[#34BFA3]', label: '用户量', key: 'user_count' },
                    { icon: 'i-material-symbols:article', color: 'text-[#F4516C]', label: '文章量', key: 'article_count' },
                    { icon: 'i-bxs:comment-dots', color: 'text-[#36A3F7]', label: '留言量', key: 'message_count' },
                ]" :key="item.key">
                    <!-- 每个统计项的布局 -->
                    <NGi>
                        <NCard>
                            <!-- 图标部分，大小为 60px，并根据数据动态设置颜色和图标 -->
                            <span class="text-[60px]" :class="[item.icon, item.color]" />
                            <!-- 显示统计数据 -->
                            <NStatistic class="float-right" :label="item.label">
                                {{ homeInfo[item.key] ?? 'unknown' }} <!-- 如果数据不存在，显示 'unknown' -->
                            </NStatistic>
                        </NCard>
                    </NGi>
                </template>
            </NGrid>

            <!-- 项目展示卡片，待完善首页设计 -->
            <NCard title="项目" size="small" class="mt-4">
                <!-- 卡片的右上角有一个额外的按钮 -->
                <template #header-extra>
                    <NButton text type="primary">
                        更多
                    </NButton>
                </template>
                <!-- v-for 渲染 5 个项目卡片 -->
                <NCard v-for="i in 5" :key="i" class="my-2 w-[300px] flex-shrink-0 cursor-pointer hover:shadow-lg"
                    title="Gin Blog Admin" size="small">
                    <!-- 项目描述 -->
                    <p class="op-60">
                        这是个基于 gin 开发的博客管理后台
                    </p>
                </NCard>
            </NCard>
        </div>
    </AppPage>
</template>


<script setup>
import { onMounted, ref } from 'vue'
import { NAvatar, NButton, NCard, NGi, NGradientText, NGrid, NStatistic } from 'naive-ui'

import AppPage from '@/components/common/AppPage.vue'
import { useUserStore } from '@/store'
import api from '@/api'

const { nickname, avatar } = useUserStore()

const homeInfo = ref({
    view_count: 0,
    user_count: 0,
    article_count: 0,
    message_count: 0,
})

onMounted(async () => {
    getOneSentence()
    const res = await api.getHomeInfo()
    homeInfo.value = res.data
})

// 一言
const sentence = ref('')
async function getOneSentence() {
    fetch('https://v1.hitokoto.cn?c=i')
        .then(resp => resp.json())
        .then(data => sentence.value = data.hitokoto)
        .catch(() => sentence.value = '宠辱不惊，看庭前花开花落；去留无意，望天上云卷云舒。')
}
</script>

<style lang="scss" scoped></style>
```

**src/views/home/route.js**

```javascript
const Layout = () => import('@/layout/index.vue')

export default {
    name: 'Home',
    path: '/',
    component: Layout,
    redirect: '/home',
    meta: {
        order: 0,
    },
    isCatalogue: true,
    children: [
        {
            name: 'Home',
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
```



## 9.4 Article 相关界面搭建

### 9.4.1 route.js



### 9.4.2 category/index.js

#### 1 CommonPage.vue

src/components/common/CommonPage.vue

这个组件的目的是封装一个具有可选页头和页脚的页面结构。通过插槽和属性灵活控制头部内容、操作按钮和是否显示页脚。它使用了 Vue 3 的 `script setup` 和 `naive-ui` 组件库来构建简洁、可复用的页面结构。

如果你需要进一步的优化或自定义，可能会根据需求调整插槽的内容或样式。

```vue
<template>
    <AppPage :show-footer="showFooter">
        <header v-if="showHeader" class="mb-3.5 min-h-[45px] flex items-center justify-between px-1">
            <slot v-if="$slots.header" name="header" />
            <template v-else>
                <h2 class="text-2xl text-[#333] font-normal dark:text-[#ccc]">
                    {{ title || $route.meta?.title }}
                </h2>
                <div class="space-x-5">
                    <slot name="action" />
                </div>
            </template>
        </header>
        <NCard class="flex-1">
            <slot />
        </NCard>
    </AppPage>
</template>

<script setup>
import { NCard } from 'naive-ui'
import AppPage from './AppPage.vue'

defineProps({
    showFooter: { type: Boolean, default: false },
    showHeader: { type: Boolean, default: true },
    title: { type: String, default: undefined },
})
</script>

<style lang="scss" scoped></style>
```



#### 2 CrudModal.vue

**src/components/crud/CrudModal.vue**

**这段代码定义了一个模态框组件，具有以下特点：**

- 支持自定义宽度、标题、按钮文本以及显示/隐藏控制。
- 使用了 Naive UI 的 `NModal` 和 `NButton` 组件来实现模态框和按钮。
- 支持插槽，允许父组件自定义模态框内容或底部按钮。
- 通过 `v-model` 和事件 (`update:visible`, `save`) 处理与父组件的交互。

你可以通过传递不同的属性来控制模态框的显示、内容、按钮文本等内容，并能通过事件与父组件进行交互。

```vue
<template>
    <NModal v-model:show="show" :style="{ width }" preset="card" :title="title" size="huge" :bordered="false">
        <slot />
        <template v-if="showFooter" #footer>
            <footer class="flex justify-end space-x-5">
                <slot name="footer">
                    <NButton @click="show = false">
                        {{ cancelText }}
                    </NButton>
                    <NButton :loading="loading" type="primary" @click="emit('save')">
                        {{ okText }}
                    </NButton>
                </slot>
            </footer>
        </template>
    </NModal>
</template>


<script setup>
import { computed } from 'vue'
import { NButton, NModal } from 'naive-ui'

const props = defineProps({
    visible: { type: Boolean, required: true },
    width: { type: String, default: '600px' },
    title: { type: String, default: '' },
    showFooter: { type: Boolean, default: true },
    loading: { type: Boolean, default: false },
    cancelText: { type: String, default: '取消' },
    okText: { type: String, default: '确定' },
})

const emit = defineEmits(['update:visible', 'save'])

const show = computed({
    get: () => props.visible,
    set: v => emit('update:visible', v),
})
</script>

<style lang="scss" scoped></style>
```



#### 3 CrudTable.vue

src/components/crud/CrudTable.vue

这段代码构建了一个带有查询栏、表格展示、分页、重置、搜索、和导出功能的 Vue 3 组件。通过与父组件的交互，支持动态查询、数据更新、行选择和导出操作。该组件非常适用于需要展示大量数据并提供灵活操作的场景。

1. **分页和远程请求**：
   - 支持前端分页和后端分页，根据 `remote` 和 `isPagination` 控制分页行为。
2. **重置和搜索**：
   - 提供了搜索和重置功能，能够根据查询条件刷新表格数据。
3. **导出功能**：
   - 支持导出表格数据为 Excel 文件，使用 `xlsx` 库实现。
4. **事件交互**：
   - 通过 `emit` 向父组件传递更新的数据或事件（例如更新查询条件、选中行、排序等）。

```vue
<template>
    <!-- 如果插槽 queryBar 存在，则渲染查询栏 -->
    <div v-if="$slots.queryBar"
        class="mb-7 min-h-[60px] flex items-start justify-between border border-gray-200 border-gray-400 rounded-2 border-solid bg-gray-50 p-3.5 dark:bg-black dark:bg-opacity-5">

        <!-- 使用 Naive UI 的 NSpace 来布局查询栏，设置间距 -->
        <NSpace wrap :size="[35, 15]">
            <!-- 渲染插槽 queryBar，允许父组件自定义查询条件部分 -->
            <slot name="queryBar" />
        </NSpace>

        <!-- 按钮区域，包含重置和搜索按钮 -->
        <div class="flex-shrink-0 space-x-4">
            <!-- 重置按钮，点击时触发 handleReset 方法 -->
            <NButton ghost type="primary" @click="handleReset">
                <template #icon>
                    <!-- 重置按钮的图标 -->
                    <i class="i-lucide:rotate-ccw" />
                </template>
                重置
            </NButton>

            <!-- 搜索按钮，点击时触发 handleSearch 方法 -->
            <NButton type="primary" @click="handleSearch">
                <template #icon>
                    <!-- 搜索按钮的图标 -->
                    <i class="i-fe:search" />
                </template>
                搜索
            </NButton>

            <!-- TODO: 未来可以添加额外的插槽，允许用户自定义其他按钮 -->
        </div>
    </div>

    <NDataTable :remote="remote" :loading="loading" :scroll-x="scrollX" :columns="columns" :data="tableData"
        :row-key="(row) => row[rowKey]" :single-line="singleLine" :pagination="isPagination ? pagination : false"
        :checked-row-keys="selections" @update:checked-row-keys="onChecked" @update:page="onPageChange"
        @update:sorter="onSorterChange" />

</template>


<script setup>
// 导入 Vue 和 Naive UI 的相关模块
import { nextTick, reactive, ref } from 'vue'
import { NButton, NDataTable, NSpace } from 'naive-ui'
import { utils, writeFile } from 'xlsx'

// 定义接收的 props 参数
const props = defineProps({
    /** 是否不设定列的分割线 */
    singleLine: { type: Boolean, default: false },
    /** true: 后端分页 false: 前端分页 */
    remote: { type: Boolean, default: true },
    /** 是否分页 */
    isPagination: { type: Boolean, default: true },
    /** 表格内容的横向宽度 */
    scrollX: { type: Number, default: 1200 },
    /** 主键 name */
    rowKey: { type: String, default: 'id' },
    /** 需要展示的列 */
    columns: { type: Array, required: true },
    /** queryBar 中的参数 */
    queryItems: {
        type: Object,
        default() { return {} },
    },
    /** 补充参数（可选） */
    extraParams: {
        type: Object,
        default() { return {} },
    },
    /** 获取数据的请求 API */
    getData: {
        type: Function,
        required: true,
    },
})

// 定义组件需要触发的事件
const emit = defineEmits(['update:queryItems', 'checked', 'dataChange', 'sorterChange'])

// 定义表格的加载状态、选中行、表格数据等响应式变量
const loading = ref(false) // 表示是否正在加载数据
const selections = ref([]) // 存储当前选中的行的 rowKey（主键）
const tableData = ref([]) // 存储表格的数据
const initQuery = { ...props.queryItems } // 初始化查询条件

// 分页配置，控制分页行为
const pagination = reactive({
    page: 1, // 当前页
    pageSize: 10, // 每页条数
    showSizePicker: true, // 是否显示选择每页多少条的选项
    pageSizes: [5, 10, 20], // 每页条数的选择范围
    // 分页变化时触发
    onChange: (page) => {
        pagination.page = page
        handleQuery()
    },
    // 页面大小变化时触发
    onUpdatePageSize: (pageSize) => {
        pagination.page = 1 // 重置为第一页
        pagination.pageSize = pageSize
        handleQuery()
    },
    // 显示分页信息的前缀
    prefix({ itemCount }) {
        return `共 ${itemCount} 条`
    },
})

// 请求数据的核心方法
async function handleQuery() {
    selections.value = [] // 重置选中的行

    try {
        loading.value = true // 开始加载数据
        let paginationParams = {}
        // 如果启用了分页并且是远程分页，则加入分页参数
        if (props.isPagination && props.remote) {
            paginationParams = {
                page_num: pagination.page,
                page_size: pagination.pageSize,
            }
        }
        // 调用父组件传递的 getData 函数请求数据
        const { data } = await props.getData({
            ...props.queryItems, // 当前的查询条件
            ...props.extraParams, // 补充的额外参数
            ...paginationParams, // 分页参数
        })
        // 更新表格数据
        tableData.value = data?.page_data || data
        pagination.itemCount = data?.total ?? data.length // 更新总数据条数
    }
    catch (error) {
        tableData.value = [] // 请求失败时清空表格数据
        pagination.itemCount = 0 // 重置数据总数
    }
    finally {
        emit('dataChange', tableData.value) // 通知父组件表格数据发生变化
        loading.value = false // 结束加载状态
    }
}

// 搜索按钮点击时的处理函数
function handleSearch() {
    pagination.page = 1 // 搜索时回到第一页
    handleQuery() // 重新请求数据
}

// 重置按钮点击时的处理函数
async function handleReset() {
    const queryItems = { ...props.queryItems } // 拷贝查询条件
    // 重置查询条件中的所有字段为 null
    for (const key in queryItems) {
        queryItems[key] = null // 注意类型问题，可能需要根据实际类型来重置
    }
    // 更新查询条件
    emit('update:queryItems', { ...queryItems, ...initQuery })
    await nextTick() // 等待 DOM 更新
    pagination.page = 1 // 回到第一页
    handleQuery() // 重新请求数据
}

// 分页变化时的处理函数
function onPageChange(currentPage) {
    pagination.page = currentPage // 更新当前页
    props.remote && handleQuery() // 如果是远程分页，则重新请求数据
}

// 表格行选择变化时的处理函数
function onChecked(rowKeys) {
    selections.value = rowKeys // 更新选中的行
    // 如果表格有选择列，则触发父组件的 'checked' 事件
    if (props.columns.some(item => item.type === 'selection')) {
        emit('checked', rowKeys)
    }
}

// 排序变化时的处理函数
function onSorterChange(sorter) {
    emit('sorterChange', sorter) // 通知父组件排序变化
}

// 导出功能，导出当前表格数据为 Excel 文件
function handleExport(columns = props.columns, data = tableData.value) {
    if (!data?.length) {
        return window.$message.warning('没有数据') // 如果没有数据，则提示
    }
    // 过滤掉没有标题或者设置了隐藏的列
    const columnsData = columns.filter(item => !!item.title && !item.hideInExcel)
    const thKeys = columnsData.map(item => item.key) // 获取列的 key
    const thData = columnsData.map(item => item.title) // 获取列的标题
    const trData = data.map(item => thKeys.map(key => item[key])) // 获取每一行的数据

    // 使用 xlsx 库创建工作表
    const sheet = utils.aoa_to_sheet([thData, ...trData])
    const workBook = utils.book_new() // 创建新的工作簿
    utils.book_append_sheet(workBook, sheet, '数据报表') // 将工作表添加到工作簿
    writeFile(workBook, '数据报表.xlsx') // 导出为 Excel 文件
}

// 暴露给父组件的 API 方法
defineExpose({
    handleQuery,
    handleSearch,
    handleReset,
    handleExport,
    selections,
    tableData,
})
</script>


<style lang="scss" scoped></style>
```



#### 4 QueryItem.vue

src/components/crud/QueryItem.vue

该组件适合用于表单或其他类似的布局场景，其中每一行显示一个标签和值，且希望能够控制标签和内容区域的宽度。例如，用于显示表单字段的名称和对应的输入框或其他内容。

1. **标签和内容布局**：这个组件主要用于显示一组标签和值对的布局。标签和内容部分通过 Flexbox 水平排列，且可以灵活控制标签和内容的宽度。
   - `labelWidth` 控制标签的宽度。
   - `contentWidth` 控制内容部分的宽度。
2. **插槽**：通过 `<slot />`，该组件可以插入自定义内容，允许父组件动态传入各种内容（例如输入框、文本等）。
3. **动态样式**：通过 `:style` 绑定，可以根据父组件传入的 `labelWidth` 和 `contentWidth` 控制元素的宽度，提升了组件的灵活性和可复用性。

```vue
<template>
    <div class="flex items-center">
        <label v-if="label" class="flex-shrink-0" :style="{ width: `${labelWidth}px` }">
            {{ label }}
        </label>
        <div class="flex-shrink-0" :style="{ width: `${contentWidth}px` }">
            <slot />
        </div>
    </div>
</template>

<script setup>
defineProps({
    label: { type: String, default: '' },
    labelWidth: { type: Number, default: 80 },
    contentWidth: { type: Number, default: 220 },
})
</script>

<style lang="scss" scoped></style>
```



#### 5 composables

**src/compostables/useCRUD.js**

1. **`useCRUD` 函数**：提供一个封装的 CRUD 操作，通过传入配置对象（如表单名称、初始化数据、创建/删除/更新函数等），返回一组与表单和弹框交互的逻辑。
2. **`modalAction`**：控制当前弹窗的操作类型（新增、编辑、查看），通过该值动态更新弹窗标题。
3. **`handleSave`**：根据当前操作类型，调用相应的 API（新增或更新），并且提供了成功后回调的功能（例如成功提示、刷新数据等）。
4. **`handleDelete`**：支持单条和批量删除，可以在删除前弹出确认框。

注释将帮助你或其他开发者更好地理解代码的业务逻辑和结构。

```javascript
import { computed, ref } from 'vue'
import { useForm } from './useForm'

const ACTIONS = {
  view: '查看',   // 查看操作
  edit: '编辑',   // 编辑操作
  add: '新增',    // 新增操作
}

/**
 * @typedef {object} FormObject
 * @property {string} name - 名称，表单模块的名称
 * @property {object} initForm - 初始表单数据，作为表单默认值
 * @property {Function} doCreate - 执行创建操作的函数
 * @property {Function} doDelete - 执行删除操作的函数
 * @property {Function} doUpdate - 执行更新操作的函数
 * @property {Function} refresh - 刷新操作，用来更新界面
 */

/**
 * 可复用的 CRUD 操作
 * @param {FormObject} options - 包含了CRUD操作所需的配置项，如表单名称、初始化表单数据、创建、删除、更新函数等
 */
export function useCRUD({ name, initForm = {}, doCreate, doDelete, doUpdate, refresh }) {
  const modalVisible = ref(false) // 弹框是否显示的状态
  /** @type {'add' | 'edit' | 'view'} 弹窗操作类型 */
  const modalAction = ref('') // 当前操作类型，用于控制是新增、编辑还是查看
  /** 弹窗加载状态 */
  const modalLoading = ref(false) // 弹窗是否在加载中
  /** 弹窗标题 */
  const modalTitle = computed(() => ACTIONS[modalAction.value] + name) // 根据当前操作类型动态生成弹窗标题

  // 表单模型和表单引用
  const { formModel: modalForm, formRef: modalFormRef, validation } = useForm(initForm)

  /** 新增操作 */
  function handleAdd() {
    modalAction.value = 'add' // 设置操作类型为“新增”
    modalVisible.value = true // 显示弹框
    modalForm.value = { ...initForm } // 重置表单为初始化的默认值
  }

  /** 修改操作 */
  function handleEdit(row) {
    modalAction.value = 'edit' // 设置操作类型为“编辑”
    modalVisible.value = true // 显示弹框
    modalForm.value = { ...row } // 设置表单数据为选中的行数据
  }

  /** 查看操作 */
  function handleView(row) {
    modalAction.value = 'view' // 设置操作类型为“查看”
    modalVisible.value = true // 显示弹框
    modalForm.value = { ...row } // 设置表单数据为选中的行数据
  }

  /** 保存操作，处理新增或编辑 */
  async function handleSave() {
    // 只有在“新增”或“编辑”时才进行保存操作
    if (!['edit', 'add'].includes(modalAction.value)) {
      modalVisible.value = false // 关闭弹框
      return
    }

    // 校验表单是否合法
    if (!(await validation())) {
      return false
    }

    // 根据操作类型选择对应的 API 函数和回调
    const actions = {
      add: {
        api: () => doCreate(modalForm.value), // 调用新增接口
        cb: () => window.$message.success('新增成功'), // 新增成功后的回调
      },
      edit: {
        api: () => doUpdate(modalForm.value), // 调用更新接口
        cb: () => window.$message.success('编辑成功'), // 编辑成功后的回调
      },
    }
    const action = actions[modalAction.value]

    try {
      modalLoading.value = true // 开始加载
      const data = await action.api() // 调用对应的 API 函数
      action.cb() // 执行操作成功后的回调
      modalLoading.value = modalVisible.value = false // 关闭加载状态和弹框
      data && refresh(data) // 刷新数据
    }
    catch (error) {
      console.error(error) // 错误处理
      modalLoading.value = false // 关闭加载状态
    }
  }

  /**
   * 删除操作，支持单条删除和批量删除
   * @param {Array} ids - 要删除的主键数组，单条删除传入单个 id，批量删除传入 id 数组
   * @param {boolean} needConfirm - 是否需要确认窗口
   */
  async function handleDelete(ids, needConfirm = true) {
    // 如果没有选中任何数据，则提示用户选择数据
    if (!ids || (Array.isArray(ids) && !ids.length)) {
      window.$message.info('请选择要删除的数据')
      return
    }

    // 调用删除接口
    const callDeleteAPI = async () => {
      try {
        modalLoading.value = true // 显示加载状态

        // 判断是否是批量删除或单条删除
        let data
        if (typeof ids === 'number' || typeof ids === 'string') {
          data = await doDelete(ids) // 单条删除
        }
        else {
          data = await doDelete(JSON.stringify(ids)) // 批量删除
        }

        // 针对软删除的情况做判断
        if (data?.code === 0) {
          window.$message.success('删除成功') // 删除成功后提示
        }
        modalLoading.value = false // 关闭加载状态
        refresh(data) // 刷新数据
      }
      catch (error) {
        console.error(error) // 错误处理
        modalLoading.value = false // 关闭加载状态
      }
    }

    // 如果需要确认窗口，则弹出确认框
    if (needConfirm) {
      window.$dialog.confirm({
        content: '确定删除？', // 确认删除提示文本
        confirm: () => callDeleteAPI(), // 点击确认后的操作
      })
    }
    else {
      callDeleteAPI() // 直接调用删除操作
    }
  }

  // 返回可供外部使用的函数和状态
  return {
    modalVisible,
    modalAction,
    modalTitle,
    modalLoading,
    handleAdd,
    handleDelete,
    handleEdit,
    handleView,
    handleSave,
    modalForm,
    modalFormRef,
  }
}
```

**src/compostables/useForm.js**

1. **`formRef`**：用于引用表单实例，通过 `ref` 获取表单的 DOM 实例，主要用来调用表单方法，如验证（`validate`）。
2. **`formModel`**：存储表单的实际数据。它是一个响应式对象，通过修改该对象的属性来控制表单字段的值。
3. **`validation`**：异步函数，通过 `formRef` 获取的表单实例来执行验证。返回 `true` 表示验证成功，返回 `false` 表示验证失败。
4. **`rules`**：定义了一个常见的验证规则（必填），可以根据需要扩展更多的规则，如格式验证、最大长度验证等。

这些注释应该帮助你和其他开发者更好地理解该代码的功能和用途。

```javascript
import { ref } from 'vue'

/**
 * 可复用的表单对象
 * 该函数用于封装表单的状态管理、验证逻辑以及表单规则的配置，可以在多个组件中复用。
 * @param {any} initForm 表单初始值，用来初始化表单数据
 * @returns {object} 返回一个包含表单引用、表单模型、验证函数和表单规则的对象
 */
export function useForm(initForm = {}) {
  const formRef = ref(null) // 表单的引用，用来获取表单实例进行操作（如验证）
  const formModel = ref({ ...initForm }) // 表单模型，保存表单数据，初始值为传入的 `initForm`

  /**
   * 表单验证函数
   * 该函数会触发表单验证，验证通过返回 `true`，否则返回 `false`
   * @returns {boolean} 验证是否通过
   */
  const validation = async () => {
    try {
      // 调用表单实例的 `validate` 方法进行表单验证
      await formRef.value?.validate()
      return true // 验证成功，返回 true
    }
    catch (error) {
      return false // 验证失败，返回 false
    }
  }

  // 表单字段的验证规则
  const rules = {
    required: {
      required: true,  // 必填项
      message: '此为必填项',  // 错误提示信息
      trigger: ['blur', 'change'],  // 在 `blur`（失去焦点）或 `change`（值改变）时触发验证
    },
  }

  // 返回包含表单引用、表单模型、验证方法和表单规则的对象
  return { formRef, formModel, validation, rules }
}
```

src/compostables/index.js

```javascript
export * from './useCRUD'
export * from './useForm'
```



#### 6 index.js

**CrudTable**

![image-20250207174807312](./assets/image-20250207174807312.png)



**CrudModal:**

![image-20250207180849852](./assets/image-20250207180849852.png)

**src/views/article/category/index.vue**

- **查询功能**：用户可以通过分类名进行查询，查询结果会根据输入的关键字进行过滤。
- **增、改、删功能**：通过 `handleAdd`、`handleEdit` 和 `handleDelete` 方法实现增、改、删操作，具体的业务逻辑在 API 层（`api.saveOrUpdateCategory`, `api.deleteCategory`）处理。
- **表格导出**：点击导出按钮时，调用表格的 `handleExport` 方法实现导出功能。

这段代码展示了一个 Vue 3 + Naive UI 实现的“分类管理”页面，包含了常见的增删改查操作和数据展示功能。它将业务逻辑分为多个组件，便于复用和维护。以下是对代码的详细分析：

1. **整体结构**

- **模板部分** (`<template>`): 定义了页面的视图结构，包括顶部操作按钮、查询框、数据表格以及新增/编辑分类的弹窗。
- **脚本部分** (`<script setup>`): 通过 Vue 3 的 `<script setup>` 语法定义了组件的逻辑，包括数据管理、方法绑定、表格列配置等。
- **样式部分** (`<style scoped>`): 使用了 scoped 样式来确保样式仅作用于当前组件。

2. **页面内容**

2.1 `CommonPage` 组件

- 该组件作为页面的容器，接受 `title` 属性来设置页面标题。
- action 插槽内放置了三个按钮：导出、添加分类和批量删除。
  - **导出**按钮：点击时调用 `$table?.handleExport()` 导出表格数据。
  - **新建分类**按钮：点击时弹出新建分类表单。
  - **批量删除**按钮：点击时删除当前选中的分类（通过 `$table?.selections` 获取选中的行）。

2.2 `CrudTable` 组件

- 这是一个复用的表格组件，接受：
  - `query-items`: 用于表格查询条件（如分类名）。
  - `columns`: 表格的列配置。
  - `get-data`: 用于获取表格数据的 API（这里是 `api.getCategorys`）。
- **查询栏**(`queryBar` 插槽)：提供了一个分类名称的输入框，用于按分类名称进行搜索。

2.3 `CrudModal` 组件

- 这是一个复用的弹窗组件，展示表单，接受以下属性：
  - `visible`: 控制弹窗的显示与隐藏。
  - `title`: 弹窗标题。
  - `loading`: 控制弹窗保存操作的加载状态。
- 弹窗中包含了一个分类名称的输入框，并且提供了表单验证（必填）。

2.4 表格列配置 (`columns`)

- 表格中包含以下列：
  - **分类名**：显示分类名称，宽度 100，居中显示。
  - **文章量**：显示分类下的文章数量，宽度 30，居中显示。
  - **创建日期**：通过 `NButton` 显示创建日期，并用图标展示。日期格式通过 `formatDate` 函数格式化。
  - **更新日期**：同样通过 `NButton` 显示更新日期，格式化后显示。
  - **操作列**：包含编辑和删除按钮，编辑按钮调用 `handleEdit(row)`，删除按钮调用 `handleDelete([row.id], false)`。

2.5 `handleAdd`, `handleEdit`, `handleDelete` 方法

- **`handleAdd`**：触发弹窗，初始化表单数据以新增分类。
- **`handleEdit`**：触发弹窗并填充表单数据，编辑已有的分类。
- **`handleDelete`**：批量删除或单项删除分类。调用 API `doDelete` 执行删除操作。

3. **逻辑部分**

3.1 `useCRUD` 自定义钩子通过 useCRUD自定义钩子封装了常见的增删改查操作。传入的配置包括：

- `name`: 表单名称（分类）。
- `doCreate`, `doDelete`, `doUpdate`: 处理新建、删除、更新的 API 函数。
- `refresh`: 操作成功后刷新表格数据。

- `useCRUD` 钩子返回了多个响应式数据和方法，供模板中调用，如：`modalVisible`, `modalTitle`, `modalLoading`, `handleSave`, `modalForm`, `modalFormRef` 等。

3.2 表格数据的获取与刷新

- `onMounted` 钩子：在组件挂载时调用 `$table.value?.handleSearch()`，触发表格数据的查询。
- 每次增删改操作成功后，调用 `refresh` 函数来刷新表格数据。

```vue
<template>
    <CommonPage title="分类管理">
        <template #action>
            <NButton type="primary" secondary @click="$table?.handleExport()">
                <template #icon>
                    <p class="i-mdi:download" />
                </template>
                导出
            </NButton>
            <NButton type="primary" @click="handleAdd">
                <template #icon>
                    <p class="i-material-symbols:add" />
                </template>
                新建分类
            </NButton>
            <NButton type="error" :disabled="!$table?.selections.length" @click="handleDelete($table?.selections)">
                <template #icon>
                    <p class="i-material-symbols:playlist-remove" />
                </template>
                批量删除
            </NButton>
        </template>
        <CrudTable ref="$table" v-model:query-items="queryItems" :columns="columns" :get-data="api.getCategorys">
            <template #queryBar>
                <QueryItem label="分类名" :label-width="50">
                    <NInput v-model:value="queryItems.keyword" clearable type="text" placeholder="请输入分类名"
                        @keydown.enter="$table?.handleSearch()" />
                </QueryItem>
            </template>
        </CrudTable>

        <CrudModal v-model:visible="modalVisible" :title="modalTitle" :loading="modalLoading" @save="handleSave">
            <NForm ref="modalFormRef" label-placement="left" label-align="left" :label-width="80" :model="modalForm">
                <NFormItem label="文章分类" path="name"
                    :rule="{ required: true, message: '请输入分类名称', trigger: ['input', 'blur'] }">
                    <NInput v-model:value="modalForm.name" placeholder="请输入分类名称" clearable />
                </NFormItem>
            </NForm>
        </CrudModal>
    </CommonPage>
</template>


<script setup>
import { h, onMounted, ref } from 'vue'
import { NButton, NForm, NFormItem, NInput, NPopconfirm } from 'naive-ui'

import CommonPage from '@/components/common/CommonPage.vue'
import QueryItem from '@/components/crud/QueryItem.vue'
import CrudModal from '@/components/crud/CrudModal.vue'
import CrudTable from '@/components/crud/CrudTable.vue'

import { formatDate } from '@/utils'
import { useCRUD } from '@/composables'
import api from '@/api'

defineOptions({ name: '分类管理' })

const $table = ref(null)
const queryItems = ref({
    keyword: '',
})

onMounted(() => {
    $table.value?.handleSearch()
})

const {
    modalVisible,
    modalTitle,
    modalLoading,
    handleAdd,
    handleDelete,
    handleEdit,
    handleSave,
    modalForm,
    modalFormRef,
} = useCRUD({
    name: '分类',
    initForm: {},
    doCreate: api.saveOrUpdateCategory,
    doDelete: api.deleteCategory,
    doUpdate: api.saveOrUpdateCategory,
    refresh: () => $table.value?.handleSearch(),
})

const columns = [
    { type: 'selection', width: 15, fixed: 'left' },
    { title: '分类名', key: 'name', width: 100, align: 'center', ellipsis: { tooltip: true } },
    { title: '文章量', key: 'article_count', width: 30, align: 'center' },
    {
        title: '创建日期',
        key: 'created_at',
        width: 80,
        align: 'center',
        render(row) {
            return h(
                NButton,
                { size: 'small', type: 'text', ghost: true },
                {
                    default: () => formatDate(row.created_at),
                    icon: () => h('i', { class: 'i-mdi:clock-time-three-outline' }),
                },
            )
        },
    },
    {
        title: '更新日期',
        key: 'updated_at',
        width: 80,
        align: 'center',
        render(row) {
            return h(
                NButton,
                { size: 'small', type: 'text', ghost: true },
                {
                    default: () => formatDate(row.updated_at),
                    icon: () => h('i', { class: 'i-mdi:update' }),
                },
            )
        },
    },
    {
        title: '操作',
        key: 'actions',
        width: 100,
        align: 'center',
        fixed: 'right',
        render(row) {
            return [
                h(
                    NButton,
                    { size: 'small', type: 'primary', onClick: () => handleEdit(row) },
                    { default: () => '编辑', icon: () => h('i', { class: 'i-material-symbols:edit-outline' }) },
                ),
                h(
                    NPopconfirm,
                    { onPositiveClick: () => handleDelete([row.id], false) },
                    {
                        trigger: () => h(
                            NButton,
                            { size: 'small', type: 'error', style: 'margin-left: 15px;' },
                            { default: () => '删除', icon: () => h('i', { class: 'i-material-symbols:delete-outline' }) },
                        ),
                        default: () => h('div', {}, '确定删除该分类吗?'),
                    },
                ),
            ]
        },
    },
]
</script>

<style lang="scss" scoped></style>
```





### 9.4.3 list/index.js



### 9.4.4 tag/index.js



### 9.4.5 write/index.js



## 9.5 Auth 相关界面搭建

