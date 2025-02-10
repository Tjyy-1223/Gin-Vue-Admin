<template>
    <CommonPage :show-header="false">
        <NTabs type="line" animated>
            <NTabPane name="website" tab="修改信息">
                <div class="m-7 flex items-center">
                    <div class="mr-7 w-50">
                        <UploadOne v-model:preview="infoForm.avatar" :width="130" />
                    </div>
                    <NForm ref="infoFormRef" label-placement="left" label-align="left" label-width="100"
                        :model="infoForm" :rules="infoFormRules" class="w-80">
                        <NFormItem label="昵称" path="nickname">
                            <NInput v-model:value="infoForm.nickname" type="text" placeholder="请填写昵称" />
                        </NFormItem>
                        <NFormItem label="个人简介" path="intro">
                            <NInput v-model:value="infoForm.intro" type="text" placeholder="请填写个人简介" />
                        </NFormItem>
                        <NFormItem label="个人网站" path="website">
                            <NInput v-model:value="infoForm.website" type="text" placeholder="请填写个人网站" />
                        </NFormItem>
                        <NButton type="primary" @click="updateProfile">
                            修改
                        </NButton>
                    </NForm>
                </div>
            </NTabPane>
            <NTabPane name="contact" tab="修改密码">
                <NForm ref="passwordFormRef" label-placement="left" label-align="left" :model="passwordForm"
                    label-width="100" :rules="passwordFormRules" class="m-[30px] w-[400px]">
                    <NFormItem label="旧密码" path="old_password">
                        <NInput v-model:value="passwordForm.old_password" type="password" show-password-on="mousedown"
                            placeholder="请输入旧密码" />
                    </NFormItem>
                    <NFormItem label="新密码" path="new_password">
                        <NInput v-model:value="passwordForm.new_password" :disabled="!passwordForm.old_password"
                            type="password" show-password-on="mousedown" placeholder="请输入新密码" />
                    </NFormItem>
                    <NFormItem label="确认密码" path="confirm_password">
                        <NInput v-model:value="passwordForm.confirm_password" :disabled="!passwordForm.new_password"
                            type="password" show-password-on="mousedown" placeholder="请再次输入新密码" />
                    </NFormItem>
                    <NButton type="primary" @click="updatePassword">
                        修改
                    </NButton>
                </NForm>
            </NTabPane>
        </NTabs>
    </CommonPage>
</template>


<script setup>
import { onMounted, ref } from 'vue'
import { NButton, NForm, NFormItem, NInput, NTabPane, NTabs } from 'naive-ui'

import CommonPage from '@/components/common/CommonPage.vue'
import UploadOne from '@/components/UploadOne.vue' // 用户头像上传组件
import { useUserStore } from '@/store' // 引入用户状态管理的 store
import api from '@/api' // API 请求模块

// 获取用户状态管理 store 实例
const userStore = useUserStore()

// 个人信息表单相关的引用和初始数据
const infoFormRef = ref(null) // 用于引用表单组件
const infoForm = ref({
    avatar: userStore.avatar, // 头像
    nickname: userStore.nickname, // 昵称
    intro: userStore.intro, // 个性签名
    website: userStore.website, // 个人网站
})

// 在组件挂载时，获取用户信息并更新表单
onMounted(async () => {
    // 从用户 store 获取最新的用户信息
    await userStore.getUserInfo()
    // 更新表单的初始数据
    infoForm.value = {
        avatar: userStore.avatar,
        nickname: userStore.nickname,
        intro: userStore.intro,
        website: userStore.website,
    }
})

// 更新个人信息函数
async function updateProfile() {
    // 校验表单
    infoFormRef.value?.validate(async (err) => {
        if (!err) { // 如果表单验证没有错误
            // 调用 API 更新当前用户的个人信息
            await api.updateCurrent(infoForm.value)
            $message.success('更新成功!') // 显示成功消息
            userStore.getUserInfo() // 更新用户信息
        }
    })
}

// 个人信息表单验证规则
const infoFormRules = {
    nickname: [
        {
            required: true, // 昵称是必填字段
            message: '请输入昵称', // 错误提示消息
            trigger: ['input', 'blur', 'change'], // 触发验证的事件
        },
    ],
}

// 修改密码相关的表单引用和数据
const passwordFormRef = ref(null) // 用于引用密码修改表单组件
const passwordForm = ref({
    old_password: '', // 旧密码
    new_password: '', // 新密码
    confirm_password: '', // 确认密码
})

// 修改密码的处理函数
function updatePassword() {
    // 校验密码表单
    passwordFormRef.value?.validate(async (err) => {
        if (!err) { // 如果表单验证没有错误
            // 调用 API 更新当前用户密码
            await api.updateCurrentPassword(passwordForm.value)
            $message.success('修改成功!') // 显示成功消息
        }
    })
}

// 密码表单验证规则
const passwordFormRules = {
    old_password: [
        {
            required: true, // 旧密码是必填字段
            message: '请输入旧密码',
            trigger: ['input', 'blur', 'change'], // 触发验证的事件
        },
    ],
    new_password: [
        {
            required: true, // 新密码是必填字段
            message: '请输入新密码',
            trigger: ['input', 'blur', 'change'],
        },
    ],
    confirm_password: [
        {
            required: true, // 确认密码是必填字段
            message: '请再次输入密码',
            trigger: ['input', 'blur'], // 触发验证的事件
        },
        {
            // 校验确认密码与新密码是否一致
            validator: validatePasswordStartWith,
            message: '两次密码输入不一致',
            trigger: 'input', // 触发事件为输入时
        },
        {
            // 校验确认密码与新密码是否一致
            validator: validatePasswordSame,
            message: '两次密码输入不一致',
            trigger: ['blur', 'password-input'], // 触发事件为失焦和密码输入时
        },
    ],
}

// 校验新密码是否和确认密码一致（密码的起始部分是否匹配，保证长度一致）
function validatePasswordStartWith(rule, value) {
    return !!passwordForm.value.new_password && passwordForm.value.new_password.startsWith(value) && passwordForm.value.new_password.length >= value.length
}

// 校验确认密码与新密码是否完全一致
function validatePasswordSame(rule, value) {
    return value === passwordForm.value.new_password
}

</script>

<style lang="scss" scoped></style>