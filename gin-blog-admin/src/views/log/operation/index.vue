<template>
    <CommonPage title="操作日志">
        <template #action>
            <NButton type="error" :disabled="!$table?.selections.length" @click="handleDelete($table?.selections)">
                <template #icon>
                    <span class="i-material-symbols:playlist-remove" />
                </template>
                批量删除
            </NButton>
        </template>

        <CrudTable ref="$table" v-model:query-items="queryItems" :columns="columns" :get-data="api.getOperationLogs">
            <template #queryBar>
                <QueryItem label="模块名" :label-width="50">
                    <NInput v-model:value="queryItems.keyword" clearable type="text" placeholder="请输入模块名或描述"
                        @keydown.enter="$table?.handleSearch()" />
                </QueryItem>
            </template>
        </CrudTable>

        <CrudModal v-model:visible="modalVisible" title="日志详情" :show-footer="false" :loading="modalLoading"
            width="full">
            <NForm ref="modalFormRef" label-placement="left" label-align="left" :label-width="90" :model="modalForm">
                <NFormItem label="操作模块: " path="opt_module">
                    {{ modalForm.opt_module }}
                </NFormItem>
                <NFormItem label="请求地址: " path="opt_url">
                    {{ modalForm.opt_url }}
                </NFormItem>
                <NFormItem label="请求方法: " path="request_method">
                    <NTag :type="tagType(modalForm.request_method)">
                        {{ modalForm.request_method }}
                    </NTag>
                </NFormItem>
                <NFormItem label="操作类型: " path="opt_type">
                    {{ modalForm.opt_type }}
                </NFormItem>
                <NFormItem label="操作方法: " path="opt_method">
                    <NCode :code="modalForm.opt_method" code-wrap language="json" />
                </NFormItem>
                <NFormItem label="操作人员: " path="nickname">
                    {{ modalForm.nickname }}
                </NFormItem>
                <NFormItem label="请求参数: " path="request_param">
                    <NCode class="word-wrap cursor-pointer p-7"
                        :code="JSON.stringify(JSON.parse(modalForm.request_param), null, 2)" language="json"
                        @click="copyFormatCode(modalForm.request_param)" />
                </NFormItem>
                <NFormItem label="返回数据: " path="response_data">
                    <NCode class="cursor-pointer p-7"
                        :code="JSON.stringify(JSON.parse(modalForm.response_data), null, 2)" language="json"
                        @click="copyFormatCode(modalForm.response_data)" />
                </NFormItem>
            </NForm>
        </CrudModal>
    </CommonPage>
</template>


<script setup>
import { h, onMounted, ref } from 'vue'
import { NButton, NCode, NForm, NFormItem, NInput, NPopconfirm, NTag } from 'naive-ui'
import { useClipboard } from '@vueuse/core'

import CommonPage from '@/components/common/CommonPage.vue'
import QueryItem from '@/components/crud/QueryItem.vue'
import CrudModal from '@/components/crud/CrudModal.vue'
import CrudTable from '@/components/crud/CrudTable.vue'

import { formatDate } from '@/utils'  // 导入格式化日期的工具函数
import { useCRUD } from '@/composables'  // 自定义的 CRUD 组合函数
import api from '@/api'  // 导入 API 请求模块

defineOptions({ name: '操作日志' })  // 定义组件名称

// 根据请求方法返回不同类型的标签 (动态计算属性)
function tagType(type) {
    switch (type) {
        case 'GET':
            return 'info'  // GET 请求显示 info 类型
        case 'POST':
            return 'success'  // POST 请求显示 success 类型
        case 'PUT':
            return 'warning'  // PUT 请求显示 warning 类型
        case 'DELETE':
            return 'error'  // DELETE 请求显示 error 类型
        default:
            return 'info'  // 默认返回 info 类型
    }
}

const $table = ref(null)  // 定义表格的引用
const queryItems = ref({
    keyword: '',  // 定义搜索关键词
})

const {
    modalVisible,
    modalLoading,
    handleDelete,
    modalForm,
    modalFormRef,
    handleView,
} = useCRUD({
    name: '日志',  // 表示操作日志的管理
    doDelete: api.deleteOperationLogs,  // 删除操作日志的 API 请求
    refresh: () => $table.value?.handleSearch(),  // 删除后刷新表格数据
})

onMounted(() => {
    // 页面加载后自动搜索操作日志
    $table.value?.handleSearch()
})

const columns = [
    { type: 'selection', width: 20, fixed: 'left' },  // 列选择框，固定在左侧
    { title: '系统模块', key: 'opt_module', width: 70, align: 'center', ellipsis: { tooltip: true } },  // 系统模块列
    { title: '操作类型', key: 'opt_type', width: 70, align: 'center', ellipsis: { tooltip: true } },  // 操作类型列
    {
        title: '请求方法',  // 请求方法列
        key: 'request_method',
        width: 80,
        align: 'center',
        ellipsis: { tooltip: true },
        render(row) {
            // 动态渲染请求方法标签，使用 `tagType` 函数设置标签类型
            return h(
                NTag,
                { type: tagType(row.request_method) },  // 使用计算属性来决定标签类型
                { default: () => row.request_method },
            )
        },
    },
    { title: '操作人员', key: 'nickname', width: 80, align: 'center', ellipsis: { tooltip: true } },  // 操作人员列
    { title: '登录IP', key: 'ip_address', width: 80, align: 'center', ellipsis: { tooltip: true } },  // 登录 IP 列
    { title: '登录地址', key: 'ip_source', width: 80, align: 'center', ellipsis: { tooltip: true } },  // 登录地址列
    {
        title: '发布时间',  // 发布时间列
        key: 'created_at',
        align: 'center',
        width: 80,
        render(row) {
            // 格式化时间并渲染为按钮
            return h(
                NButton,
                { size: 'small', type: 'text', ghost: true },
                {
                    default: () => formatDate(row.created_at),  // 格式化时间
                    icon: () => h('i', { class: 'i-mdi:update' }),  // 显示更新时间图标
                },
            )
        },
    },
    {
        title: '操作',  // 操作列，包含查看和删除按钮
        key: 'actions',
        width: 120,
        align: 'center',
        fixed: 'right',  // 固定在右侧
        render(row) {
            // 渲染查看和删除按钮
            return [
                h(
                    NButton,
                    {
                        size: 'small',
                        quaternary: true,
                        type: 'info',
                        onClick: () => handleView(row),  // 点击查看按钮时，调用 handleView
                    },
                    {
                        default: () => '查看',
                        icon: () => h('i', { class: 'i-ic:outline-remove-red-eye' }),  // 查看图标
                    },
                ),
                h(
                    NPopconfirm,
                    { onPositiveClick: () => handleDelete([row.id], false) },  // 点击确认删除时调用 handleDelete
                    {
                        trigger: () =>
                            h(
                                NButton,
                                {
                                    size: 'small',
                                    quaternary: true,
                                    type: 'error',
                                    style: 'margin-left: 15px;',  // 删除按钮样式
                                },
                                {
                                    default: () => '删除',
                                    icon: () => h('i', { class: 'i-material-symbols:delete-outline' }),  // 删除图标
                                },
                            ),
                        default: () => h('div', {}, '确定删除该日志吗?'),  // 删除确认的提示框
                    },
                ),
            ]
        },
    },
]

function copyFormatCode(code) {
    // 使用 VueUse 的 useClipboard 来复制格式化后的 JSON 代码
    const { copy } = useClipboard()
    copy(JSON.stringify(JSON.parse(code), null, 2))  // 格式化 JSON 代码后复制
    window.$message.success('内容已复制到剪切板!')  // 提示复制成功
}

</script>

<style lang="scss" scoped></style>