<template>
    <!-- 如果插槽 queryBar 存在，则渲染查询栏 -->
    <div v-if="$slots.queryBar"
        class="mb-7 min-h-[60px] flex items-start justify-between border border-gray-200 border-gray-400 rounded-2 border-solid bg-gray-50 p-3.5 dark:bg-black dark:bg-opacity-5">

        <!-- 使用 Naive UI 的 NSpace 来布局查询栏，设置间距 -->
        <NSpace wrap :size="[35, 15]">
            <!-- 渲染插槽 queryBar，允许父组件自定义查询条件部分 -->
            <slot name="queryBar" />
        </NSpace>

        <!-- 按钮区域，包含重置和搜索按钮 -->
        <div class="flex-shrink-0 space-x-4">
            <!-- 重置按钮，点击时触发 handleReset 方法 -->
            <NButton ghost type="primary" @click="handleReset">
                <template #icon>
                    <!-- 重置按钮的图标 -->
                    <i class="i-lucide:rotate-ccw" />
                </template>
                重置
            </NButton>

            <!-- 搜索按钮，点击时触发 handleSearch 方法 -->
            <NButton type="primary" @click="handleSearch">
                <template #icon>
                    <!-- 搜索按钮的图标 -->
                    <i class="i-fe:search" />
                </template>
                搜索
            </NButton>

            <!-- TODO: 未来可以添加额外的插槽，允许用户自定义其他按钮 -->
        </div>
    </div>

    <NDataTable :remote="remote" :loading="loading" :scroll-x="scrollX" :columns="columns" :data="tableData"
        :row-key="(row) => row[rowKey]" :single-line="singleLine" :pagination="isPagination ? pagination : false"
        :checked-row-keys="selections" @update:checked-row-keys="onChecked" @update:page="onPageChange"
        @update:sorter="onSorterChange" />

</template>


<script setup>
// 导入 Vue 和 Naive UI 的相关模块
import { nextTick, reactive, ref } from 'vue'
import { NButton, NDataTable, NSpace } from 'naive-ui'
import { utils, writeFile } from 'xlsx'

// 定义接收的 props 参数
const props = defineProps({
    /** 是否不设定列的分割线 */
    singleLine: { type: Boolean, default: false },
    /** true: 后端分页 false: 前端分页 */
    remote: { type: Boolean, default: true },
    /** 是否分页 */
    isPagination: { type: Boolean, default: true },
    /** 表格内容的横向宽度 */
    scrollX: { type: Number, default: 1200 },
    /** 主键 name */
    rowKey: { type: String, default: 'id' },
    /** 需要展示的列 */
    columns: { type: Array, required: true },
    /** queryBar 中的参数 */
    queryItems: {
        type: Object,
        default() { return {} },
    },
    /** 补充参数（可选） */
    extraParams: {
        type: Object,
        default() { return {} },
    },
    /** 获取数据的请求 API */
    getData: {
        type: Function,
        required: true,
    },
})

// 定义组件需要触发的事件
const emit = defineEmits(['update:queryItems', 'checked', 'dataChange', 'sorterChange'])

// 定义表格的加载状态、选中行、表格数据等响应式变量
const loading = ref(false) // 表示是否正在加载数据
const selections = ref([]) // 存储当前选中的行的 rowKey（主键）
const tableData = ref([]) // 存储表格的数据
const initQuery = { ...props.queryItems } // 初始化查询条件

// 分页配置，控制分页行为
const pagination = reactive({
    page: 1, // 当前页
    pageSize: 10, // 每页条数
    showSizePicker: true, // 是否显示选择每页多少条的选项
    pageSizes: [5, 10, 20], // 每页条数的选择范围
    // 分页变化时触发
    onChange: (page) => {
        pagination.page = page
        handleQuery()
    },
    // 页面大小变化时触发
    onUpdatePageSize: (pageSize) => {
        pagination.page = 1 // 重置为第一页
        pagination.pageSize = pageSize
        handleQuery()
    },
    // 显示分页信息的前缀
    prefix({ itemCount }) {
        return `共 ${itemCount} 条`
    },
})

// 请求数据的核心方法
async function handleQuery() {
    selections.value = [] // 重置选中的行

    try {
        loading.value = true // 开始加载数据
        let paginationParams = {}
        // 如果启用了分页并且是远程分页，则加入分页参数
        if (props.isPagination && props.remote) {
            paginationParams = {
                page_num: pagination.page,
                page_size: pagination.pageSize,
            }
        }
        // 调用父组件传递的 getData 函数请求数据
        const { data } = await props.getData({
            ...props.queryItems, // 当前的查询条件
            ...props.extraParams, // 补充的额外参数
            ...paginationParams, // 分页参数
        })
        // 更新表格数据
        tableData.value = data?.page_data || data
        pagination.itemCount = data?.total ?? data.length // 更新总数据条数
    }
    catch (error) {
        tableData.value = [] // 请求失败时清空表格数据
        pagination.itemCount = 0 // 重置数据总数
    }
    finally {
        emit('dataChange', tableData.value) // 通知父组件表格数据发生变化
        loading.value = false // 结束加载状态
    }
}

// 搜索按钮点击时的处理函数
function handleSearch() {
    pagination.page = 1 // 搜索时回到第一页
    handleQuery() // 重新请求数据
}

// 重置按钮点击时的处理函数
async function handleReset() {
    const queryItems = { ...props.queryItems } // 拷贝查询条件
    // 重置查询条件中的所有字段为 null
    for (const key in queryItems) {
        queryItems[key] = null // 注意类型问题，可能需要根据实际类型来重置
    }
    // 更新查询条件
    emit('update:queryItems', { ...queryItems, ...initQuery })
    await nextTick() // 等待 DOM 更新
    pagination.page = 1 // 回到第一页
    handleQuery() // 重新请求数据
}

// 分页变化时的处理函数
function onPageChange(currentPage) {
    pagination.page = currentPage // 更新当前页
    props.remote && handleQuery() // 如果是远程分页，则重新请求数据
}

// 表格行选择变化时的处理函数
function onChecked(rowKeys) {
    selections.value = rowKeys // 更新选中的行
    // 如果表格有选择列，则触发父组件的 'checked' 事件
    if (props.columns.some(item => item.type === 'selection')) {
        emit('checked', rowKeys)
    }
}

// 排序变化时的处理函数
function onSorterChange(sorter) {
    emit('sorterChange', sorter) // 通知父组件排序变化
}

// 导出功能，导出当前表格数据为 Excel 文件
function handleExport(columns = props.columns, data = tableData.value) {
    if (!data?.length) {
        return window.$message.warning('没有数据') // 如果没有数据，则提示
    }
    // 过滤掉没有标题或者设置了隐藏的列
    const columnsData = columns.filter(item => !!item.title && !item.hideInExcel)
    const thKeys = columnsData.map(item => item.key) // 获取列的 key
    const thData = columnsData.map(item => item.title) // 获取列的标题
    const trData = data.map(item => thKeys.map(key => item[key])) // 获取每一行的数据

    // 使用 xlsx 库创建工作表
    const sheet = utils.aoa_to_sheet([thData, ...trData])
    const workBook = utils.book_new() // 创建新的工作簿
    utils.book_append_sheet(workBook, sheet, '数据报表') // 将工作表添加到工作簿
    writeFile(workBook, '数据报表.xlsx') // 导出为 Excel 文件
}

// 暴露给父组件的 API 方法
defineExpose({
    handleQuery,
    handleSearch,
    handleReset,
    handleExport,
    selections,
    tableData,
})
</script>


<style lang="scss" scoped></style>