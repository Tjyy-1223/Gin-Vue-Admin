<template>
  <!-- FIXME: 使用 style="background-image: url(/image/login_bg.webp);" 不生效, 需要写到 style 里的 class 中 -->
  <AppPage class="backgroundImg bg-cover">
    <div style="transform: translateY(25px)" class="m-auto max-w-[700px] min-w-[345px] flex items-center justify-center rounded-2 bg-white bg-opacity-60 p-4 shadow">
      <div class="hidden w-[380px] px-5 py-9 md:block">
        <img src="/image/login_banner.webp" class="w-full" alt="login_banner">
      </div>

      <div class="w-[320px] flex flex-col px-4 py-9 space-y-5.5">
        <h5 class="flex items-center justify-center text-2xl text-gray font-normal">
          <img src="/image/logo.svg" alt="logo" class="mr-2 h-[50px] w-[50px]">
          <span> {{ title }} </span>
        </h5>
        <NInput
          v-model:value="loginForm.username"
          class="h-[50px] items-center pl-2"
          autofocus
          placeholder="test@qq.com"
          :maxlength="20"
        />
        <NInput
          v-model:value="loginForm.password"
          class="h-[50px] items-center pl-2"
          type="password"
          show-password-on="mousedown"
          placeholder="11111"
          :maxlength="20"
          @keydown.enter="handleLogin"
        />
        <NCheckbox
          :checked="isRemember"
          label="记住我"
          :on-update:checked="(val) => (isRemember = val)"
        />
        <NButton
          class="h-[50px] w-full rounded-5"
          type="primary"
          :loading="loading"
          @click="handleLogin"
        >
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

import AppPage from '@/components/common/AddPage.vue'

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