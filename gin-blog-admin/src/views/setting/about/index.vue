<template>
    <CommonPage :show-header="false">
        <div class="mb-4 flex items-center justify-between">
            <div class="mx-1 text-2xl font-bold">
                关于我
            </div>
            <NButton type="primary" :loading="btnLoading" @click="handleSave">
                <template #icon>
                    <p v-if="!btnLoading" class="i-line-md:confirm-circle" />
                </template>
                保存
            </NButton>
        </div>
        <!-- TODO: 文件上传封装 -->
        <MdEditor v-model="aboutContent" style="height: calc(100vh - 245px)" />
    </CommonPage>
</template>


<script setup>
import { onMounted, ref } from 'vue'  // 导入 Vue 3 的响应式 API
import { NButton } from 'naive-ui'   // 导入 Naive UI 的按钮组件

// 导入 Markdown 编辑器并引入相关样式
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'

// 导入自定义的公共页面组件
import CommonPage from '@/components/common/CommonPage.vue'

// 导入 API 请求模块
import api from '@/api'

// 设置组件名称为 "关于我"
defineOptions({ name: '关于我' })

// 定义响应式数据
const aboutContent = ref('')  // 用于存储 "关于我" 页面内容（Markdown 格式）
const btnLoading = ref(false)  // 用于控制保存按钮的加载状态（防止多次点击）

// 页面加载时获取关于我的内容
onMounted(async () => {
    // 调用 API 获取 "关于我" 页面内容
    const resp = await api.getAbout()
    // 将返回的数据赋值给 aboutContent（页面内容）
    aboutContent.value = resp.data
})

// 保存编辑内容的处理函数
async function handleSave() {
    try {
        // 设置按钮为加载状态，防止重复点击
        btnLoading.value = true
        // 调用 API 更新关于我页面的内容，传入编辑后的 content
        await api.updateAbout({ content: aboutContent.value })
        // 弹出成功消息
        window.$message.success('更新成功')
    }
    finally {
        // 不管请求成功或失败，都将按钮的加载状态设置为 false
        btnLoading.value = false
    }
}

</script>

<style lang="scss" scoped>
.md-preview {

    ul,
    ol {
        list-style: revert;
    }
}
</style>