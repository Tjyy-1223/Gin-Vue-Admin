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





## 9.5 Auth 相关界面搭建

