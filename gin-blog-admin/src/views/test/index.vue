<template>
    <CommonPage title="测试页面">
        <CrudTable ref="$table" v-model:query-items="queryItems" :columns="columns" :get-data="api.getCategorys">
            <template #queryBar>
                <QueryItem label="分类名" :label-width="50">
                    <NInput v-model:value="queryItems.keyword" clearable type="text" placeholder="请输入分类名"
                        @keydown.enter="$table?.handleSearch()" />
                </QueryItem>
            </template>
        </CrudTable>
    </CommonPage>
</template>


<script setup>
import { defineComponent, h, nextTick, onMounted, ref } from 'vue'  // 导入 Vue 相关的 API
import { NInput } from 'naive-ui'  // 导入 Naive UI 中的输入框组件

// 导入其他自定义组件
import CommonPage from '@/components/common/CommonPage.vue'
import QueryItem from '@/components/crud/QueryItem.vue'
import CrudTable from '@/components/crud/CrudTable.vue'

import api from '@/api'  // 导入 API 模块，用于与后台进行交互

// 设置当前组件的名称为 "分类管理"
defineOptions({ name: '分类管理' })

// 定义响应式数据
const $table = ref(null)  // 用于引用表格组件，进行操作
const queryItems = ref({  // 存储查询条件，当前只有一个关键字查询项
    keyword: '',
})

// 当前编辑行的索引
const editIndex = ref(-1)

// 在组件挂载完成后，调用表格的搜索方法
onMounted(() => {
    $table.value?.handleSearch()  // 调用表格实例的 handleSearch 方法进行数据查询
})

// 定义一个组件，用于渲染可编辑的表格单元格
const ShowOrEdit = defineComponent({
    props: {  // 组件的输入属性
        value: [String, Number],  // 单元格的初始值，可以是字符串或数字
        onUpdateValue: [Function, Array],  // 更新单元格值时调用的回调函数
        rowIndex: [Number],  // 当前单元格所在行的索引
    },
    setup(props) {  // 组件的逻辑部分
        const inputRef = ref(null)  // 用于引用输入框实例
        const inputValue = ref(props.value)  // 输入框的值，默认等于传入的 value 值

        // 处理点击单元格事件，进入编辑模式
        function handleOnClick() {
            editIndex.value = props.rowIndex  // 设置当前编辑行的索引
            nextTick(() => {
                inputRef.value.focus()  // 等待 DOM 更新后聚焦到输入框
            })
        }

        // 处理输入框值的变化
        function handleChange() {
            props.onUpdateValue(inputValue.value)  // 更新父组件传入的值
        }

        // 返回渲染的虚拟 DOM，条件渲染输入框或值文本
        return () =>
            h(  // 返回一个 div 元素
                'div',
                {
                    style: 'min-height: 22px',  // 设置最小高度
                    onClick: handleOnClick,  // 点击时进入编辑状态
                },
                editIndex.value === props.rowIndex  // 如果当前行是编辑状态
                    ? h(NInput, {  // 渲染一个输入框
                        ref: inputRef,  // 引用输入框
                        value: inputValue.value,  // 绑定输入框的值
                        onUpdateValue: (v) => {  // 处理输入框值变化
                            inputValue.value = v
                        },
                        onChange: handleChange,  // 输入框值变化时更新父组件
                        onBlur: handleChange,  // 输入框失去焦点时更新父组件
                    })
                    : props.value,  // 否则直接渲染单元格的原始值
            )
    },
})

// 表格的列定义
const columns = [
    { type: 'selection', width: 15, fixed: 'left' },  // 选择框列，宽度为 15，固定在左侧
    {
        title: '创建日期',  // 表头显示 "创建日期"
        key: 'created_at',  // 字段对应的 key
        width: 80,  // 列宽度
        align: 'center',  // 列内容居中对齐
        render(row, index) {  // 自定义渲染函数
            return h(
                ShowOrEdit,  // 使用 ShowOrEdit 组件进行渲染
                {
                    value: row.created_at,  // 将行数据中的 "创建日期" 传递给组件
                    rowIndex: index,  // 当前行的索引
                    onUpdateValue(v) {  // 组件更新时的回调函数
                        data.value[index].created_at = v  // 更新当前行的 "创建日期" 字段
                    },
                },
            )
        },
    },
    {
        title: '更新日期',  // 表头显示 "更新日期"
        key: 'updated_at',  // 字段对应的 key
        width: 80,  // 列宽度
        align: 'center',  // 列内容居中对齐
        render(row, index) {  // 自定义渲染函数
            return h(
                ShowOrEdit,  // 使用 ShowOrEdit 组件进行渲染
                {
                    value: row.updated_at,  // 将行数据中的 "更新日期" 传递给组件
                    rowIndex: index,  // 当前行的索引
                    onUpdateValue(v) {  // 组件更新时的回调函数
                        data.value[index].updated_at = v  // 更新当前行的 "更新日期" 字段
                    },
                },
            )
        },
    },
    {
        title: '操作',  // 表头显示 "操作"
        key: 'action',  // 字段对应的 key
        width: 80,  // 列宽度
        align: 'center',  // 列内容居中对齐
        render(row, index) {  // 自定义渲染函数
            if (editIndex.value === index) {  // 如果当前行在编辑状态
                return h(
                    'button',  // 渲染一个按钮
                    {
                        onClick() {  // 点击按钮时关闭编辑
                            editIndex.value = -1  // 关闭编辑状态
                        },
                    },
                    '关闭编辑',  // 按钮文本
                )
            }
            else {  // 如果当前行不在编辑状态
                return h(
                    'button',  // 渲染一个按钮
                    {
                        onClick() {  // 点击按钮时进入编辑状态
                            editIndex.value = index  // 设置当前编辑行的索引
                        },
                    },
                    '进入编辑',  // 按钮文本
                )
            }
        },
    },
]
</script>


<style lang="scss" scoped></style>