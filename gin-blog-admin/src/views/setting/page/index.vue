<template>
    <CommonPage title="页面管理">
        <template #action>
            <NButton type="primary" @click="handleAdd">
                <template #icon>
                    <p class="i-material-symbols:add" />
                </template>
                新建页面
            </NButton>
        </template>
        <div class="flex flex-wrap justify-between">
            <div v-for="page of pageList" :key="page.id" class="relative my-2 w-[300px] cursor-pointer text-center">
                <div class="absolute right-2 top-1 text-white">
                    <NDropdown :options="options" @select="handleSelect($event, page)">
                        <p class="i-ion:ellipsis-horizontal h-5 w-5 text-white hover:text-blue" />
                    </NDropdown>
                </div>
                <NImage :src="convertImgUrl(page.cover)" height="170" width="300"
                    :img-props="{ style: { 'border-radius': '5px' } }" />
                <p class="text-base">
                    {{ page.name }}
                </p>
            </div>
            <div class="h-0 w-[300px]" />
            <div class="h-0 w-[300px]" />
            <div class="h-0 w-[300px]" />
        </div>

        <CrudModal v-model:visible="modalVisible" width="550px" :title="modalTitle" :loading="modalLoading"
            @save="handleSave">
            <NForm ref="modalFormRef" label-placement="left" label-align="left" :label-width="80" :model="modalForm">
                <NFormItem label="页面名称" path="name"
                    :rule="{ required: true, message: '请输入页面名称', trigger: ['input', 'blur'] }">
                    <NInput v-model:value="modalForm.name" placeholder="页面名称" />
                </NFormItem>
                <NFormItem label="页面标签" path="label"
                    :rule="{ required: true, message: '请输入页面标签', trigger: ['input', 'blur'] }">
                    <NInput v-model:value="modalForm.label" placeholder="页面标签" />
                </NFormItem>
                <NFormItem label="页面封面" path="cover"
                    :rule="{ required: true, message: '请上传封面图片', trigger: ['input', 'blur'] }">
                    <div class="w-full flex items-center justify-between">
                        <UploadOne ref="uploadOneRef" v-model:preview="modalForm.cover" :width="300"
                            @finish="val => (modalForm.cover = val)" />

                        <span class="i-uiw:reload h-5 w-5 cursor-pointer" :class="reloadFlag ? 'animate-spin' : ''"
                            @click="refreshImg(modalForm.cover)" />
                    </div>
                </NFormItem>
                <NFormItem label="封面链接" path="cover"
                    :rule="{ required: true, message: '请输入封面链接', trigger: ['input', 'blur'] }">
                    <NInput v-model:value="modalForm.cover" type="textarea" placeholder="图片上传成功自动生成，或者直接复制外链" />
                </NFormItem>
            </NForm>
        </CrudModal>
    </CommonPage>
</template>


<script setup>
import { h, onMounted, ref } from 'vue'  // 导入 Vue 相关的 API，创建响应式变量等
import { NButton, NDropdown, NForm, NFormItem, NImage, NInput } from 'naive-ui'  // 导入 Naive UI 组件库的相关组件

import CrudModal from '@/components/crud/CrudModal.vue'  // 导入一个通用的 CRUD 操作模态框组件
import UploadOne from '@/components//UploadOne.vue'  // 导入一个文件上传组件（用于上传单个文件）
import CommonPage from '@/components/common/CommonPage.vue'  // 导入一个通用页面布局组件

import { convertImgUrl } from '@/utils'  // 导入工具函数，用于处理图片 URL
import { useCRUD } from '@/composables'  // 导入自定义的 useCRUD 组合式 API，用于处理 CRUD 操作
import api from '@/api'  // 导入 API 请求模块，用于与后台进行交互

// FIXME: 只有这个页面的 KeepAlive 为什么没有生效？
// 说明该页面的 KeepAlive 缓存没有生效，可能存在某些问题

// 使用 useCRUD 组合式 API，封装通用的增删改查操作
const {
    modalVisible,    // 控制模态框显示的状态
    modalTitle,      // 模态框的标题
    modalLoading,    // 控制模态框加载状态
    handleAdd,       // 处理添加数据的操作
    handleDelete,    // 处理删除数据的操作
    handleEdit,      // 处理编辑数据的操作
    handleSave,      // 处理保存数据的操作
    modalForm,       // 模态框中的表单数据
    modalFormRef,    // 引用模态框表单的 Ref 对象
} = useCRUD({
    name: '页面',  // 当前操作的实体名称（这里是 "页面"）
    initForm: {},  // 初始化表单数据（这里为空对象）
    doCreate: api.saveOrUpdatePage,  // 创建或更新页面的 API 请求函数
    doDelete: api.deletePage,  // 删除页面的 API 请求函数
    doUpdate: api.saveOrUpdatePage,  // 更新页面的 API 请求函数
    refresh: fetchData,  // 数据刷新函数
})

// 定义响应式数据
const pageList = ref([])  // 存储页面列表数据
const reloadFlag = ref(false)  // 控制图片刷新标志
const uploadOneRef = ref(null)  // 图片上传组件的引用

// 在组件挂载后，获取页面数据
onMounted(async () => {
    fetchData()  // 调用函数获取页面数据
})

// 获取页面列表数据
async function fetchData() {
    const resp = await api.getPages()  // 调用 API 获取页面列表
    pageList.value = resp.data  // 将获取到的数据赋值给 pageList
}

// 刷新图片预览，当用户输入新的图片 URL 时
function refreshImg(img) {
    reloadFlag.value = true  // 设置刷新标志为 true，触发图片重新加载
    uploadOneRef.value.previewImg = img  // 更新上传组件的预览图片
    // 设置一个 600ms 后重置刷新标志
    setTimeout(() => reloadFlag.value = false, 600)
}

// 处理下拉菜单选择的操作（编辑或删除）
function handleSelect(key, page) {
    if (key === 'edit') {
        handleEdit(page)  // 如果选择了 "编辑"，则调用编辑操作
    }
    else if (key === 'delete') {
        handleDelete([page.id])  // 如果选择了 "删除"，则调用删除操作，并传递页面的 ID
    }
}

// 定义下拉菜单的选项（编辑和删除）
const options = [
    {
        label: '编辑',  // 下拉菜单项的标签为 "编辑"
        key: 'edit',    // 该项的 key 为 "edit"
        icon: () => h('i', { class: 'i-mingcute:edit-2-line' }),  // 编辑项的图标
    },
    {
        label: '删除',  // 下拉菜单项的标签为 "删除"
        key: 'delete',  // 该项的 key 为 "delete"
        icon: () => h('i', { class: 'i-mingcute:delete-back-line' }),  // 删除项的图标
    },
]

</script>

<style lang="scss" scoped></style>