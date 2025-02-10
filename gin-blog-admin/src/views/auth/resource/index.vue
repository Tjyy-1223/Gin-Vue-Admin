<template>
    <CommonPage title="接口管理">
        <template #action>
            <NButton type="primary" @click="handleAddModule">
                <template #icon>
                    <span class="i-material-symbols:add" />
                </template>
                新增模块
            </NButton>
        </template>

        <CrudTable ref="$table" v-model:query-items="queryItems" :is-pagination="false" :columns="columns"
            :get-data="api.getResources" :single-line="true">
            <template #queryBar>
                <QueryItem label="资源名" :label-width="50">
                    <NInput v-model:value="queryItems.keyword" clearable type="text" placeholder="请输入资源名"
                        @keydown.enter="$table?.handleSearch()" />
                </QueryItem>
            </template>
        </CrudTable>

        <CrudModal v-model:visible="modalVisible" :title="modalTitle" :loading="modalLoading" @save="handleSave">
            <NForm ref="modalFormRef" label-placement="left" label-align="left" :label-width="80" :model="modalForm">
                <NFormItem label="资源名" path="name">
                    <NInput v-model:value="modalForm.name" placeholder="请输入资源名" />
                </NFormItem>
                <NFormItem label="资源路径" path="url">
                    <NInput v-model:value="modalForm.url" placeholder="请输入资源路径" />
                </NFormItem>
                <NFormItem label="请求方式" path="request_method">
                    <NRadioGroup v-model:value="modalForm.request_method" name="radiogroup">
                        <NSpace>
                            <NRadio v-for="method of requestMethods" :key="method" :value="method">
                                <NGradientText :type="tagType(method)">
                                    {{ method }}
                                </NGradientText>
                            </NRadio>
                        </NSpace>
                    </NRadioGroup>
                </NFormItem>
            </NForm>
        </CrudModal>

        <CrudModal v-model:visible="moduleModalVisible" :title="`${modalAction === 'add' ? '新增' : '编辑'}模块`"
            :loading="modalVisible" @save="handleModuleSave">
            <NForm ref="modalFormRef" label-placement="left" label-align="left" :label-width="80" :model="modalForm">
                <NFormItem label="模块名" path="name">
                    <NInput v-model:value="modalForm.name" placeholder="请输入模块名" />
                </NFormItem>
            </NForm>
        </CrudModal>
    </CommonPage>
</template>


<script setup>
import { h, onMounted, ref } from 'vue'
import { NButton, NForm, NFormItem, NGradientText, NInput, NPopconfirm, NRadio, NRadioGroup, NSpace, NSwitch, NTag } from 'naive-ui'

import CommonPage from '@/components/common/CommonPage.vue'
import QueryItem from '@/components/crud/QueryItem.vue'
import CrudModal from '@/components/crud/CrudModal.vue'
import CrudTable from '@/components/crud/CrudTable.vue'

import { formatDate } from '@/utils'
import { useCRUD } from '@/composables'
import api from '@/api'

defineOptions({ name: '接口管理' }) // 设置页面名称

// 定义引用，用于操作表格
const $table = ref(null)
// 定义查询条件
const queryItems = ref({
    keyword: '',
})

// 使用 CRUD 封装逻辑
const {
    modalVisible,  // 控制模态框显示与隐藏
    modalAction,   // 当前操作，新增、编辑
    modalTitle,    // 模态框标题
    modalLoading,  // 模态框加载状态
    handleAdd,     // 处理新增操作
    handleDelete,  // 处理删除操作
    handleEdit,    // 处理编辑操作
    handleSave,    // 处理保存操作
    modalForm,     // 模态框表单数据
    modalFormRef,  // 模态框表单引用
} = useCRUD({
    name: '接口',  // 资源名称
    doCreate: api.saveOrUpdateResource,  // 创建接口资源
    doDelete: api.deleteResource,       // 删除接口资源
    doUpdate: api.saveOrUpdateResource, // 更新接口资源
    refresh: () => $table.value?.handleSearch(), // 刷新表格数据
})

// 页面挂载时进行查询
onMounted(() => {
    $table.value?.handleSearch() // 初始加载接口数据
})

// 请求方法类型，用于后续的标签渲染
const requestMethods = ['GET', 'POST', 'DELETE', 'PUT']

// 根据请求方法返回不同的标签类型
function tagType(type) {
    switch (type) {
        case 'GET':
            return 'info'  // 信息标签
        case 'POST':
            return 'success' // 成功标签
        case 'PUT':
            return 'warning' // 警告标签
        case 'DELETE':
            return 'error' // 错误标签
        default:
            return 'info'
    }
}

// 表格列定义
const columns = [
    {
        title: '资源名称',  // 列标题
        key: 'name',       // 数据字段
        width: 80,         // 列宽度
        ellipsis: { tooltip: true },  // 超出文本显示为省略号并提供工具提示
    },
    {
        title: '资源路径',
        key: 'url',
        width: 80,
        ellipsis: { tooltip: true },
        render(row) {
            return row.children ? '-' : h('span', { class: 'color-[#1890ff]' }, row.url)  // 渲染资源路径
        },
    },
    {
        title: '请求方式',
        key: 'request_method',
        width: 50,
        align: 'center',
        render(row) {
            return row.children
                ? '-'
                : h(
                    NTag,
                    { type: tagType(row.request_method) }, // 使用计算属性渲染标签
                    { default: () => row.request_method },
                )
        },
    },
    {
        title: '匿名访问',
        key: 'is_hidden',
        width: 50,
        align: 'center',
        fixed: 'left',
        render(row) {
            return row.children
                ? '-'
                : h(NSwitch, {
                    size: 'small',  // 设置开关大小
                    rubberBand: false,  // 禁用橡皮筋效果
                    value: row.is_anonymous,  // 当前匿名访问状态
                    loading: !!row.publishing, // 显示加载状态
                    onUpdateValue: () => handleUpdateAnonymous(row),  // 处理切换匿名访问
                })
        },
    },
    {
        title: '创建日期',
        key: 'created_at',
        width: 60,
        render(row) {
            return h('span', formatDate(row.created_at))  // 格式化日期并显示
        },
    },
    {
        title: '操作',
        key: 'actions',
        width: 115,
        align: 'center',
        fixed: 'right',
        render(row) {
            return [
                // 新增按钮
                h(
                    NButton,
                    {
                        size: 'tiny',
                        quaternary: true,
                        type: 'primary',
                        style: `display: ${row.children ? '' : 'none'};`, // 根据是否有子资源决定显示
                        onClick: () => {
                            handleAdd()  // 打开新增模态框
                            modalForm.value.parent_id = row.id  // 设置父资源id
                        },
                    },
                    { default: () => '新增', icon: () => h('i', { class: 'i-material-symbols:add' }) },
                ),
                // 编辑按钮
                h(
                    NButton,
                    {
                        size: 'tiny',
                        quaternary: true,
                        type: 'info',
                        onClick: () => (row.children ? handleEditModule(row) : handleEdit(row)),  // 判断是编辑模块还是资源
                    },
                    { default: () => '编辑', icon: () => h('i', { class: 'i-material-symbols:edit-outline' }) },
                ),
                // 删除按钮
                h(
                    NPopconfirm,
                    {
                        onPositiveClick: () => {
                            handleDelete(row.id, false)  // 删除操作
                        },
                    },
                    {
                        trigger: () =>
                            h(
                                NButton,
                                { size: 'tiny', quaternary: true, type: 'error' },
                                { default: () => '删除', icon: () => h('i', { class: 'i-material-symbols:delete-outline' }) },
                            ),
                        default: () => h('div', {}, '确定删除该接口吗?'),  // 弹出确认框提示
                    },
                ),
            ]
        },
    },
]

// 修改匿名访问状态
async function handleUpdateAnonymous(row) {
    if (!row.id) {
        return  // 确保有id
    }
    row.publishing = true  // 启动加载动画
    row.is_anonymous = !row.is_anonymous  // 切换匿名访问状态
    try {
        await api.updateResourceAnonymous(row)  // 调用 API 更新
        $message?.success(row.is_anonymous ? '已允许匿名访问' : '已禁止匿名访问')  // 成功提示
    }
    catch (err) {
        row.is_anonymous = !row.is_anonymous  // 回滚操作
        console.error(err)
    }
    finally {
        row.publishing = false  // 结束加载动画
    }
}

// 模块相关
const moduleModalVisible = ref(false)  // 控制模块模态框的显示
// 打开新增模块的模态框
function handleAddModule() {
    modalAction.value = 'add'  // 设置操作类型为新增
    modalForm.value = {}  // 清空表单数据
    moduleModalVisible.value = true  // 显示模态框
}
// 编辑模块
function handleEditModule(row) {
    modalAction.value = 'edit'  // 设置操作类型为编辑
    modalForm.value = { ...row }  // 复制数据到表单
    moduleModalVisible.value = true  // 显示模态框
}
// 保存模块数据
async function handleModuleSave() {
    handleSave()  // 保存操作
    moduleModalVisible.value = false  // 关闭模态框
}
</script>


<style lang="scss" scoped></style>