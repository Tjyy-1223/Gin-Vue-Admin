<template>
    <CommonPage title="留言管理">
        <template #action>
            <NButton type="error" :disabled="!$table?.selections.length" @click="handleDelete($table?.selections)">
                <template #icon>
                    <p class="i-material-symbols:recycling-rounded" />
                </template>
                批量删除
            </NButton>
            <NButton type="success" :disabled="!$table?.selections.length"
                @click="handleUpdateReview($table.selections, true)">
                <template #icon>
                    <p class="i-ic:outline-approval" />
                </template>
                批量通过
            </NButton>
        </template>
        <NTabs type="line" animated @update:value="handleChangeTab">
            <template #prefix>
                状态
            </template>
            <NTabPane name="all" tab="全部" />
            <NTabPane name="has_review" tab="通过" />
            <NTabPane name="not_review" tab="审核中" />
        </NTabs>
        <CrudTable ref="$table" v-model:query-items="queryItems" :extra-params="extraParams" :columns="columns"
            :get-data="api.getMessages">
            <template #queryBar>
                <QueryItem label="用户" :label-width="40" :content-width="180">
                    <NInput v-model:value="queryItems.nickname" clearable type="text" placeholder="请输入用户昵称"
                        @keydown.enter=" $table?.handleSearch()" />
                </QueryItem>
            </template>
        </CrudTable>
    </CommonPage>
</template>


<script setup>
import { h, onMounted, ref } from 'vue'
import { NButton, NImage, NInput, NPopconfirm, NTabPane, NTabs, NTag } from 'naive-ui'

import CommonPage from '@/components/common/CommonPage.vue'
import QueryItem from '@/components/crud/QueryItem.vue'
import CrudTable from '@/components/crud/CrudTable.vue'

import { convertImgUrl, formatDate } from '@/utils'
import { useCRUD } from '@/composables'
import api from '@/api'

// 设置组件的名称
defineOptions({ name: '留言管理' })

// 组件挂载完成后，默认显示所有留言
onMounted(() => {
    handleChangeTab('all') // 默认查看全部
})

// 定义响应式引用
const $table = ref(null) // 用来引用 CrudTable 组件实例
const queryItems = ref({
    nickname: '', // 搜索条件：留言人昵称
})
const extraParams = ref({
    is_review: null, // 评论状态：审核中 | 通过
})

// 使用自定义的 useCRUD hook 处理删除操作
const { handleDelete } = useCRUD({
    name: '留言', // 资源名称
    doDelete: api.deleteMessages, // 删除 API
    refresh: () => $table.value?.handleSearch(), // 删除后刷新表格
})

// 表格列的配置
const columns = [
    { type: 'selection', width: 15, fixed: 'left' }, // 选择框列
    {
        title: '头像', // 列标题：头像
        key: 'avatar', // 列字段：头像
        width: 40, // 列宽
        align: 'center', // 内容居中对齐
        render(row) {
            return h(NImage, {
                'height': 40,
                'imgProps': { style: { 'border-radius': '3px' } }, // 设置图片圆角
                'src': convertImgUrl(row.avatar), // 获取头像 URL
                'fallback-src': 'http://dummyimage.com/400x400', // 如果头像加载失败，显示占位图
                'show-toolbar-tooltip': true,
            })
        },
    },
    {
        title: '留言人',
        key: 'nickname',
        width: 60,
        align: 'center',
        ellipsis: { tooltip: true }, // 超过宽度时显示 tooltip
    },
    {
        title: '留言内容',
        key: 'content',
        width: 120,
        align: 'center',
    },
    {
        title: 'IP 地址',
        key: 'ip_address',
        width: 70,
        align: 'center',
        ellipsis: { tooltip: true },
    },
    {
        title: 'IP 来源',
        key: 'ip_source',
        width: 70,
        align: 'center',
        ellipsis: { tooltip: true },
        render(row) {
            return h('span', row.ip_source || '未知') // 显示 IP 来源，若为空则显示 "未知"
        },
    },
    {
        title: '留言时间',
        key: 'created_at',
        align: 'center',
        width: 80,
        render(row) {
            return h(
                NButton,
                { size: 'small', type: 'text', ghost: true },
                {
                    default: () => formatDate(row.created_at), // 格式化时间
                    icon: () => h('i', { class: 'i-mdi:update' }), // 设置图标
                },
            )
        },
    },
    {
        title: '状态',
        key: 'is_review',
        width: 50,
        align: 'center',
        render(row) {
            return h(
                NTag,
                { type: row.is_review ? 'success' : 'error' }, // 根据状态显示不同颜色
                { default: () => (row.is_review ? '通过' : '审核中') },
            )
        },
    },
    {
        title: '操作',
        key: 'actions',
        width: 100,
        align: 'center',
        fixed: 'right', // 操作列固定在右侧
        render(row) {
            return [
                // 根据评论状态显示不同的按钮：通过 / 撤下
                row.is_review
                    ? h(
                        NButton,
                        {
                            size: 'small',
                            type: 'warning',
                            onClick: () => handleUpdateReview([row.id], false), // 点击撤下按钮
                        },
                        {
                            default: () => '撤下', // 按钮文本
                            icon: () => h('i', { class: 'i-mi:circle-error' }), // 图标
                        },
                    )
                    : h(
                        NButton,
                        {
                            size: 'small',
                            type: 'success',
                            style: 'margin-left: 15px;', // 按钮间隔
                            onClick: () => handleUpdateReview([row.id], true), // 点击通过按钮
                        },
                        {
                            default: () => '通过',
                            icon: () => h('i', { class: 'i-mi:circle-check' }),
                        },
                    ),
                // 删除操作，点击确认删除
                h(
                    NPopconfirm,
                    { onPositiveClick: () => handleDelete([row.id], false) }, // 点击删除确认按钮时触发删除
                    {
                        trigger: () =>
                            h(
                                NButton,
                                { size: 'small', type: 'error', style: 'margin-left: 15px;' },
                                { default: () => '删除', icon: () => h('i', { class: 'i-material-symbols:delete-outline' }) },
                            ),
                        default: () => h('div', {}, '确定删除该条留言吗?'), // 删除确认提示文本
                    },
                ),
            ]
        },
    },
]

// 修改留言审核状态
async function handleUpdateReview(ids, is_review) {
    if (!ids.length) {
        $message.info('请选择要审核的数据') // 如果没有选择留言，弹出提示
        return
    }

    // 调用 API 更新留言审核状态
    await api.updateMessageReview(ids, is_review)
    // 提示成功或失败
    $message?.success(is_review ? '审核成功' : '撤下成功')
    // 刷新表格
    $table.value?.handleSearch()
}

// 切换标签页，筛选不同的留言状态：全部、通过、审核中
function handleChangeTab(value) {
    switch (value) {
        case 'all':
            extraParams.value.is_review = null // 查看全部留言
            break
        case 'has_review': // 通过
            extraParams.value.is_review = 1
            break
        case 'not_review': // 审核中
            extraParams.value.is_review = 0
            break
    }
    // 切换标签后刷新表格
    $table.value?.handleSearch()
}
</script>

<style lang="scss" scoped></style>