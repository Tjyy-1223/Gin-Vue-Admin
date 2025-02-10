<template>
    <CommonPage title="在线用户">
        <CrudTable ref="$table" v-model:query-items="queryItems" :columns="columns" :get-data="api.getOnlineUsers"
            :is-pagination="false">
            <template #queryBar>
                <QueryItem label="用户名 | 昵称" :label-width="100" :content-width="200">
                    <NInput v-model:value="queryItems.keyword" clearable type="text" placeholder="搜索关键字"
                        @keydown.enter="$table?.handleSearch()" />
                </QueryItem>
            </template>
        </CrudTable>
    </CommonPage>
</template>


<script setup>
// 导入 Vue 和 Naive UI 相关组件
import { h, onMounted, ref } from 'vue'
import { NButton, NImage, NInput, NPopconfirm } from 'naive-ui'

// 导入自定义组件和工具函数
import CommonPage from '@/components/common/CommonPage.vue'
import QueryItem from '@/components/crud/QueryItem.vue'
import CrudTable from '@/components/crud/CrudTable.vue'

import { convertImgUrl, formatDate } from '@/utils'  // 图片路径转换和日期格式化工具函数
import api from '@/api'  // API 请求模块

// 设置当前组件名称为 "在线用户"
defineOptions({ name: '在线用户' })

// 表格引用和查询条件
const $table = ref(null)  // 存储 CrudTable 组件的引用
const queryItems = ref({
    keyword: '',  // 用户名或昵称查询
})

// 在组件挂载时触发表格搜索
onMounted(() => {
    $table.value?.handleSearch()  // 触发表格的初始搜索
})

// 定义表格列
const columns = [
    {
        title: '头像',  // 表格列标题
        key: 'avatar',  // 列的 key
        width: 30,  // 列宽度
        align: 'center',  // 内容居中对齐
        render(row) {  // 自定义渲染函数
            return h(NImage, {
                'height': 30,  // 图片高度
                'src': convertImgUrl(row.info.avatar),  // 获取图片 URL
                'fallback-src': 'http://dummyimage.com/400x400',  // 加载失败时显示的占位图
                'show-toolbar-tooltip': true,  // 显示工具提示
            })
        },
    },
    {
        title: '昵称',
        key: 'nickname',
        width: 60,
        align: 'center',
        ellipsis: { tooltip: true },  // 启用溢出文本显示工具提示
        render(row) {
            return h('span', row.info.nickname || '未知')  // 渲染昵称
        },
    },
    {
        title: '登录 IP',
        key: 'ip_address',
        width: 70,
        align: 'center',
        ellipsis: { tooltip: true },
        render(row) {
            return h('span', row.ip_address || '未知')  // 渲染登录 IP 地址
        },
    },
    {
        title: '登录地址',
        key: 'ip_source',
        width: 70,
        align: 'center',
        ellipsis: { tooltip: true },
        render(row) {
            return h('span', row.ip_source || '未知')  // 渲染登录地址
        },
    },
    {
        title: '登录浏览器',
        key: 'browser',
        width: 70,
        align: 'center',
        ellipsis: { tooltip: true },
        render(row) {
            return h('span', row.browser || '未知')  // 渲染浏览器信息
        },
    },
    {
        title: '操作系统',
        key: 'os',
        width: 70,
        align: 'center',
        ellipsis: { tooltip: true },
        render(row) {
            return h('span', row.os || '未知')  // 渲染操作系统信息
        },
    },
    {
        title: '登录时间',
        key: 'last_login_time',
        align: 'center',
        width: 70,
        render(row) {
            return h('span', formatDate(row.last_login_time, 'YYYY-MM-DD HH:mm:ss'))  // 格式化登录时间
        },
    },
    {
        title: '操作',
        key: 'actions',
        width: 60,
        align: 'center',
        fixed: 'right',  // 固定在右侧
        render(row) {
            return h(
                NPopconfirm,
                { onPositiveClick: () => handleForceOffline(row) },  // 用户点击确认后触发下线操作
                {
                    trigger: () =>
                        h(
                            NButton,
                            { size: 'small', type: 'warning' },  // 下线按钮
                            {
                                default: () => '下线',  // 按钮文字
                                icon: () => h('i', { class: 'i-material-symbols:delete-outline' }),  // 图标
                            },
                        ),
                    default: () => h('div', {}, '确定强制该用户下线吗?'),  // 弹出确认框内容
                },
            )
        },
    },
]

// 强制用户下线的处理函数
async function handleForceOffline(row) {
    try {
        await api.forceOfflineUser(row.id)  // 调用 API 强制下线用户
        window.$message.success('该用户已被强制下线!')  // 弹出成功消息
        $table.value?.handleSearch()  // 刷新表格数据
    }
    catch (err) {
        console.error(err)  // 捕获错误并打印日志
    }
}
</script>


<style lang="scss" scoped></style>