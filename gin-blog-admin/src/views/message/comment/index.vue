<template>
    <CommonPage title="评论管理">
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
            :get-data="api.getComments">
            <template #queryBar>
                <QueryItem label="用户" :label-width="40" :content-width="180">
                    <NInput v-model:value="queryItems.nickname" clearable type="text" placeholder="请输入用户昵称"
                        @keydown.enter="$table?.handleSearch()" />
                </QueryItem>
                <QueryItem label="来源" :label-width="40" :content-width="160">
                    <NSelect v-model:value="queryItems.type" clearable filterablec placeholder="请选择评论来源"
                        :options="commentTypeOptions" @update:value="$table?.handleSearch()" />
                </QueryItem>
            </template>
        </CrudTable>
    </CommonPage>
</template>


<script setup>
import { h, onMounted, ref } from 'vue'
import { NButton, NImage, NInput, NPopconfirm, NSelect, NTabPane, NTabs, NTag } from 'naive-ui'

import CommonPage from '@/components/common/CommonPage.vue'
import QueryItem from '@/components/crud/QueryItem.vue'
import CrudTable from '@/components/crud/CrudTable.vue'

import { commentTypeMap, commentTypeOptions } from '@/assets/config'
import { convertImgUrl, formatDate } from '@/utils'
import { useCRUD } from '@/composables'
import api from '@/api'

// 设置组件的名称
defineOptions({ name: '评论管理' })

// 组件挂载完成后，默认显示所有评论
onMounted(() => {
    handleChangeTab('all') // 默认查看全部
})

// 定义响应式引用
const $table = ref(null) // 用来引用 CrudTable 组件实例
const queryItems = ref({
    nickname: '', // 搜索条件：评论人昵称
    type: '', // 搜索条件：评论类型
})
const extraParams = ref({
    is_review: null, // 评论状态：审核中 | 通过
})

// 使用自定义的 useCRUD hook 处理删除操作
const { handleDelete } = useCRUD({
    name: '评论', // 资源名称
    doDelete: api.deleteComments, // 删除 API
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
                'src': convertImgUrl(row.user?.info?.avatar), // 获取用户头像 URL
                'fallback-src': 'http://dummyimage.com/400x400', // 如果头像加载失败，显示占位图
                'show-toolbar-tooltip': true,
            })
        },
    },
    {
        title: '评论人', // 列标题：评论人
        key: 'nickname', // 列字段：评论人昵称
        width: 50, // 列宽
        align: 'center',
        ellipsis: { tooltip: true }, // 超过宽度时显示 tooltip
        render(row) {
            return h('span', row.user?.info?.nickname || '无') // 显示评论人的昵称，若为空则显示 "无"
        },
    },
    // TODO: 合理的显示评论的文章信息
    {
        title: '评论类型', // 列标题：评论类型
        key: '', // 没有直接映射的字段
        width: 50,
        align: 'center',
        render(row) {
            if (row.type === 1) { // 如果评论类型是 1，则显示 "文章"
                return [
                    h(NTag, { type: 'info' }, { default: () => '文章' }),
                ]
            }
            if (row.type === 2) { // 如果评论类型是 2，则显示 "友链"
                return h(NTag, { type: 'success' }, { default: () => '友链' })
            }
        },
    },
    {
        title: '回复对象', // 列标题：回复对象
        key: 'reply_nick_name', // 列字段：回复对象昵称
        width: 50,
        align: 'center',
        render(row) {
            return h('span', row.reply_user?.info?.nickname || '-') // 显示回复对象的昵称，若为空则显示 "-"
        },
    },
    {
        title: '评论内容', // 列标题：评论内容
        key: 'content', // 列字段：评论内容
        width: 140, // 列宽
        align: 'center',
        ellipsis: { tooltip: true }, // 超过宽度时显示 tooltip
    },
    {
        title: '评论时间', // 列标题：评论时间
        key: 'created_at', // 列字段：评论创建时间
        align: 'center',
        width: 60,
        render(row) {
            return h(
                NButton,
                { size: 'small', type: 'text', ghost: true }, // 按钮类型为文本，外观透明
                {
                    default: () => formatDate(row.created_at), // 格式化评论时间
                    icon: () => h('i', { class: 'i-mdi:update' }), // 设置图标
                },
            )
        },
    },
    {
        title: '状态', // 列标题：评论状态
        key: 'is_review', // 列字段：评论审核状态
        width: 50,
        align: 'center',
        render(row) {
            return h(
                NTag,
                { type: row.is_review ? 'success' : 'error' }, // 根据审核状态显示不同颜色
                { default: () => (row.is_review ? '通过' : '审核中') }, // 显示 "通过" 或 "审核中"
            )
        },
    },
    {
        title: '来源', // 列标题：来源
        key: 'type', // 列字段：评论类型
        width: 50,
        align: 'center',
        render(row) {
            return h(
                NTag,
                { type: commentTypeMap[row.type].tag }, // 根据评论类型显示不同的标签
                { default: () => commentTypeMap[row.type].name }, // 显示评论类型的名称
            )
        },
    },
    {
        title: '操作', // 列标题：操作
        key: 'actions', // 列字段：操作
        width: 100,
        align: 'center',
        fixed: 'right', // 操作列固定在右侧
        render(row) {
            return [
                // 根据评论审核状态显示不同的按钮：通过 / 撤下
                row.is_review
                    ? h(
                        NButton,
                        {
                            size: 'small',
                            type: 'warning',
                            style: 'margin-left: 15px;', // 按钮间隔
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
                        default: () => h('div', {}, '确定删除该条评论吗?'), // 删除确认提示文本
                    },
                ),
            ]
        },
    },
]

// 修改评论审核状态
async function handleUpdateReview(ids, is_review) {
    if (!ids.length) {
        window.$message.info('请选择要审核的数据') // 如果没有选择评论，弹出提示
        return
    }
    // 调用 API 更新评论审核状态
    await api.updateCommentReview(ids, is_review)
    // 提示成功或失败
    window.$message?.success(is_review ? '审核成功' : '撤下成功')
    // 刷新表格
    $table.value?.handleSearch()
}

// 切换标签页，筛选不同的评论状态：全部、通过、审核中
function handleChangeTab(value) {
    switch (value) {
        case 'all':
            extraParams.value.is_review = null // 查看全部评论
            break
        case 'has_review': // 通过
            extraParams.value.is_review = true
            break
        case 'not_review': // 审核中
            extraParams.value.is_review = false
            break
    }
    // 切换标签后刷新表格
    $table.value?.handleSearch()
}

</script>

<style lang="scss" scoped></style>