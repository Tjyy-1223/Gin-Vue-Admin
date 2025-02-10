<template>
    <CommonPage title="用户列表">
        <CrudTable ref="$table" v-model:query-items="queryItems" :columns="columns" :get-data="api.getUsers">
            <template #queryBar>
                <QueryItem label="昵称" :label-width="40" :content-width="160">
                    <NInput v-model:value="queryItems.nickname" clearable type="text" placeholder="请输入昵称"
                        @keydown.enter="$table?.handleSearch()" />
                </QueryItem>
                <QueryItem label="用户名" :label-width="60" :content-width="160">
                    <NInput v-model:value="queryItems.username" clearable type="text" placeholder="请输入用户名"
                        @keydown.enter="$table?.handleSearch()" />
                </QueryItem>
                <QueryItem label="登录方式" :label-width="70" :content-width="160">
                    <NSelect v-model:value="queryItems.login_type" clearable filterable placeholder="请选择登录方式"
                        :options="loginTypeOptions" @update:value="$table?.handleSearch()" />
                </QueryItem>
            </template>
        </CrudTable>

        <CrudModal v-model:visible="modalVisible" title="修改用户" :loading="modalLoading" @save="handleSave">
            <NForm ref="modalFormRef" label-placement="left" label-align="left" :label-width="80" :model="modalForm">
                <NFormItem label="用户昵称" path="name">
                    <NInput v-model:value="modalForm.nickname" clearable placeholder="请输入用户昵称" />
                </NFormItem>
                <NFormItem label="角色" path="role_ids">
                    <NCheckboxGroup v-model:value="modalForm.role_ids">
                        <NSpace item-style="display: flex;">
                            <NCheckbox v-for="item in roleOptions" :key="item.value" :value="item.value"
                                :label="item.label" />
                        </NSpace>
                    </NCheckboxGroup>
                </NFormItem>
            </NForm>
        </CrudModal>
    </CommonPage>
</template>


<script setup>
// 导入 Vue 和 Naive UI 相关组件
import { h, onMounted, ref } from 'vue'
import { NButton, NCheckbox, NCheckboxGroup, NForm, NFormItem, NImage, NInput, NSelect, NSpace, NSwitch, NTag } from 'naive-ui'

// 导入自定义组件和工具函数
import CommonPage from '@/components/common/CommonPage.vue'
import QueryItem from '@/components/crud/QueryItem.vue'
import CrudModal from '@/components/crud/CrudModal.vue'
import CrudTable from '@/components/crud/CrudTable.vue'

import { loginTypeMap, loginTypeOptions } from '@/assets/config'  // 登录方式相关配置
import { convertImgUrl, formatDate } from '@/utils'  // 图片路径转换和日期格式化工具函数
import { useCRUD } from '@/composables'  // 自定义的CRUD逻辑钩子
import api from '@/api'  // API 请求模块

// 设置当前组件名称为 "用户列表"
defineOptions({ name: '用户列表' })

// 表格引用和查询条件
const $table = ref(null)  // 存储 CrudTable 组件的引用
const queryItems = ref({
    username: '',  // 用户名查询
    nickname: '',  // 昵称查询
    login_type: null,  // 登录方式查询
})

// 使用 CRUD 操作的自定义钩子，包含添加、编辑、保存等逻辑
const {
    modalVisible,  // 模态框的显示状态
    modalLoading,  // 模态框加载状态
    handleSave,  // 保存数据的方法
    handleEdit,  // 编辑数据的方法
    modalForm,  // 模态框表单数据
    modalFormRef,  // 模态框表单引用
} = useCRUD({
    name: '用户',  // 当前操作的数据模型名称
    doUpdate: api.updateUser,  // 更新用户的 API 请求
    refresh: () => $table.value?.handleSearch(),  // 更新成功后刷新表格
})

// 存储角色选项
const roleOptions = ref([])

// 在组件挂载时，加载角色选项并触发表格搜索
onMounted(() => {
    api.getRoleOption().then(resp => roleOptions.value = resp.data)  // 获取角色选项
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
            return h(NImage, {  // 渲染头像图片
                'height': 30,
                'imgProps': { style: { 'border-radius': '3px' } },  // 设置圆角
                'src': convertImgUrl(row.info?.avatar),  // 获取图片 URL
                'fallback-src': 'http://dummyimage.com/400x400',  // 加载失败的占位图
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
            return h('span', row.info?.nickname)  // 渲染昵称
        },
    },
    {
        title: '登录方式',
        key: 'login_type',
        width: 40,
        align: 'center',
        render(row) {
            return h(
                NTag,
                { type: loginTypeMap[row.login_type]?.tag },  // 显示不同类型的登录方式
                { default: () => loginTypeMap[row.login_type]?.name || '未知' },  // 显示登录方式名称
            )
        },
    },
    {
        title: '用户角色',
        key: 'role',
        width: 80,
        align: 'center',
        render(row) {
            // 如果是超级管理员，则显示特殊标签
            if (row.is_super) {
                return h(NTag, { type: 'error' }, { default: () => '超级管理员' })
            }
            // 渲染用户的多个角色
            const roles = row.roles ?? []
            const groups = []
            for (let i = 0; i < roles.length; i++) {
                groups.push(h(NTag, { type: 'info', style: { margin: '2px 3px' } }, { default: () => roles[i].name }))
            }
            return h('span', groups.length ? groups : '无')  // 如果没有角色，显示 "无"
        },
    },
    {
        title: '登录 IP',
        key: 'ip_address',
        width: 70,
        align: 'center',
        ellipsis: { tooltip: true },
        render(row) {
            return h('span', row.ip_address || '未知')  // 显示登录 IP 地址
        },
    },
    {
        title: '登录地址',
        key: 'ip_source',
        width: 70,
        align: 'center',
        ellipsis: { tooltip: true },
        render(row) {
            return h('span', row.ip_source || '未知')  // 显示登录地址
        },
    },
    {
        title: '创建时间',
        key: 'created_at',
        align: 'center',
        width: 70,
        render(row) {
            return h(
                NButton,
                { size: 'small', type: 'text', ghost: true },
                {
                    default: () => formatDate(row.created_at),  // 格式化创建时间
                    icon: () => h('i', { class: 'i-mdi:update' }),  // 显示图标
                },
            )
        },
    },
    {
        title: '上次登录时间',
        key: 'last_login_time',
        align: 'center',
        width: 70,
        render(row) {
            return h(
                NButton,
                { size: 'small', type: 'text', ghost: true },
                {
                    default: () => formatDate(row.last_login_time),  // 格式化上次登录时间
                    icon: () => h('i', { class: 'i-mdi:update' }),  // 显示图标
                },
            )
        },
    },
    {
        title: '禁用',
        key: 'is_disable',
        width: 30,
        align: 'center',
        fixed: 'left',
        render(row) {
            return h(NSwitch, {
                size: 'small',
                rubberBand: false,
                value: row.is_disable,  // 控制禁用状态
                loading: !!row.publishing,  // 显示加载状态
                onUpdateValue: () => handleUpdateDisable(row),  // 更新禁用状态的回调函数
            })
        },
    },
    {
        title: '操作',
        key: 'actions',
        width: 60,
        align: 'center',
        fixed: 'right',
        render(row) {
            return [
                h(
                    NButton,
                    {
                        size: 'small',
                        type: 'primary',
                        onClick: () => {
                            row.nickname = row.info?.nickname  // 复制昵称
                            row.role_ids = row.roles.map(e => e.id)  // 获取角色ID
                            handleEdit(row)  // 触发编辑操作
                        },
                    },
                    {
                        default: () => '编辑',
                        icon: () => h('i', { class: 'i-material-symbols:delete-outline' }),  // 编辑按钮图标
                    },
                ),
            ]
        },
    },
]

// 修改用户禁用状态
async function handleUpdateDisable(row) {
    if (!row.id) {
        return  // 如果没有用户ID，则直接返回
    }
    row.publishing = true  // 设置为正在发布状态
    row.is_disable = !row.is_disable  // 切换禁用状态
    try {
        await api.updateUserDisable(row.id, row.is_disable)  // 调用 API 更新禁用状态
        $message?.success(row.is_disable ? '已禁用该用户' : '已取消禁用该用户')  // 显示成功消息
        $table.value?.handleSearch()  // 刷新表格数据
    }
    catch (err) {
        row.is_disable = !row.is_disable  // 如果失败，恢复原来的禁用状态
        console.error(err)  // 打印错误
    }
    finally {
        row.publishing = false  // 无论成功与否，取消发布状态
    }
}
</script>



<style lang="scss" scoped></style>