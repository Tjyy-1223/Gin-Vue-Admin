<template>
    <UModal v-model="loginFlag" :width="480">
        <div class="mx-2 my-1">
            <div class="mb-4 text-xl font-bold">
                登录
            </div>
            <div class="my-7 space-y-4">
                <div class="flex items-center">
                    <span class="mr-4 inline-block w-16 text-right"> 用户名 </span>
                    <input v-model="form.username" required placeholder="用户名"
                        class="block w-full border-0 rounded-md p-2 text-gray-900 shadow-sm outline-none ring-1 ring-gray-300 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-emerald">
                </div>
                <div class="flex items-center">
                    <span class="mr-4 inline-block w-16 text-right"> 密码 </span>
                    <input v-model="form.password" type="password" placeholder="密码"
                        class="block w-full border-0 rounded-md p-2 text-gray-900 shadow-sm outline-none ring-1 ring-gray-300 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-emerald">
                </div>
            </div>
            <div class="my-2 text-center">
                <button class="w-full rounded-lg bg-blue-200 py-2 text-black hover:bg-light-blue" @click="handleLogin">
                    登录
                </button>
                <div class="mt-4 flex justify-between">
                    <button @click="openRegister">
                        立即注册
                    </button>
                    <button @click="openForget">
                        忘记密码？
                    </button>
                </div>
            </div>
        </div>
    </UModal>
</template>

<script setup lang="js">
import { computed, ref } from 'vue'

import UModal from '@/components/ui/UModal.vue'
import { useAppStore, useUserStore } from '@/store'
import api from '@/api'

const userStore = useUserStore()
const appStore = useAppStore()

const registerFlag = computed({
    get: () => appStore.registerFlag,
    set: val => appStore.setRegisterFlag(val),
})

const loginFlag = computed({
    get: () => appStore.loginFlag,
    set: val => appStore.setLoginFlag(val),
})

const form = ref({
    username: 'test@qq.com',
    password: '11111',
})


// 登陆操作
const doLogin = async (username, password) => {
    const resp = await api.login({username, password})
    window.$notify?.success('登录成功!')
    // 设置 token
    userStore.setToken(resp.data.token)
    // 加载用户信息, 更新 pinia 中信息, 刷新页面
    await userStore.getUserInfo()
    // 清空表单
    form.value = { username: 'test@qq.com', password: '11111' }
    loginFlag.value = false
}

// 登陆处理
async function handleLogin() {
    const { username, password } = form.value
    if (!username || !password) {
        window.$message?.warning('请输入用户名和密码')
        return
    }
    doLogin(username, password)
}


// 立即注册: 弹出注册窗口、关闭登陆窗口
function openRegister() {
    registerFlag.value = true
    loginFlag.value = false
}

// TODO:忘记密码
function openForget() {
    window?.$message?.info('暂时不支持找回密码!')
}

</script>

<style scoped></style>