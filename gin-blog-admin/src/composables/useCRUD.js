import { computed, ref } from 'vue'
import { useForm } from './useForm'

const ACTIONS = {
  view: '查看',   // 查看操作
  edit: '编辑',   // 编辑操作
  add: '新增',    // 新增操作
}

/**
 * @typedef {object} FormObject
 * @property {string} name - 名称，表单模块的名称
 * @property {object} initForm - 初始表单数据，作为表单默认值
 * @property {Function} doCreate - 执行创建操作的函数
 * @property {Function} doDelete - 执行删除操作的函数
 * @property {Function} doUpdate - 执行更新操作的函数
 * @property {Function} refresh - 刷新操作，用来更新界面
 */

/**
 * 可复用的 CRUD 操作
 * @param {FormObject} options - 包含了CRUD操作所需的配置项，如表单名称、初始化表单数据、创建、删除、更新函数等
 */
export function useCRUD({ name, initForm = {}, doCreate, doDelete, doUpdate, refresh }) {
  const modalVisible = ref(false) // 弹框是否显示的状态
  /** @type {'add' | 'edit' | 'view'} 弹窗操作类型 */
  const modalAction = ref('') // 当前操作类型，用于控制是新增、编辑还是查看
  /** 弹窗加载状态 */
  const modalLoading = ref(false) // 弹窗是否在加载中
  /** 弹窗标题 */
  const modalTitle = computed(() => ACTIONS[modalAction.value] + name) // 根据当前操作类型动态生成弹窗标题

  // 表单模型和表单引用
  const { formModel: modalForm, formRef: modalFormRef, validation } = useForm(initForm)

  /** 新增操作 */
  function handleAdd() {
    modalAction.value = 'add' // 设置操作类型为“新增”
    modalVisible.value = true // 显示弹框
    modalForm.value = { ...initForm } // 重置表单为初始化的默认值
  }

  /** 修改操作 */
  function handleEdit(row) {
    modalAction.value = 'edit' // 设置操作类型为“编辑”
    modalVisible.value = true // 显示弹框
    modalForm.value = { ...row } // 设置表单数据为选中的行数据
  }

  /** 查看操作 */
  function handleView(row) {
    modalAction.value = 'view' // 设置操作类型为“查看”
    modalVisible.value = true // 显示弹框
    modalForm.value = { ...row } // 设置表单数据为选中的行数据
  }

  /** 保存操作，处理新增或编辑 */
  async function handleSave() {
    // 只有在“新增”或“编辑”时才进行保存操作
    if (!['edit', 'add'].includes(modalAction.value)) {
      modalVisible.value = false // 关闭弹框
      return
    }

    // 校验表单是否合法
    if (!(await validation())) {
      return false
    }

    // 根据操作类型选择对应的 API 函数和回调
    const actions = {
      add: {
        api: () => doCreate(modalForm.value), // 调用新增接口
        cb: () => window.$message.success('新增成功'), // 新增成功后的回调
      },
      edit: {
        api: () => doUpdate(modalForm.value), // 调用更新接口
        cb: () => window.$message.success('编辑成功'), // 编辑成功后的回调
      },
    }
    const action = actions[modalAction.value]

    try {
      modalLoading.value = true // 开始加载
      const data = await action.api() // 调用对应的 API 函数
      action.cb() // 执行操作成功后的回调
      modalLoading.value = modalVisible.value = false // 关闭加载状态和弹框
      data && refresh(data) // 刷新数据
    }
    catch (error) {
      console.error(error) // 错误处理
      modalLoading.value = false // 关闭加载状态
    }
  }

  /**
   * 删除操作，支持单条删除和批量删除
   * @param {Array} ids - 要删除的主键数组，单条删除传入单个 id，批量删除传入 id 数组
   * @param {boolean} needConfirm - 是否需要确认窗口
   */
  async function handleDelete(ids, needConfirm = true) {
    // 如果没有选中任何数据，则提示用户选择数据
    if (!ids || (Array.isArray(ids) && !ids.length)) {
      window.$message.info('请选择要删除的数据')
      return
    }

    // 调用删除接口
    const callDeleteAPI = async () => {
      try {
        modalLoading.value = true // 显示加载状态

        // 判断是否是批量删除或单条删除
        let data
        if (typeof ids === 'number' || typeof ids === 'string') {
          data = await doDelete(ids) // 单条删除
        }
        else {
          data = await doDelete(JSON.stringify(ids)) // 批量删除
        }

        // 针对软删除的情况做判断
        if (data?.code === 0) {
          window.$message.success('删除成功') // 删除成功后提示
        }
        modalLoading.value = false // 关闭加载状态
        refresh(data) // 刷新数据
      }
      catch (error) {
        console.error(error) // 错误处理
        modalLoading.value = false // 关闭加载状态
      }
    }

    // 如果需要确认窗口，则弹出确认框
    if (needConfirm) {
      window.$dialog.confirm({
        content: '确定删除？', // 确认删除提示文本
        confirm: () => callDeleteAPI(), // 点击确认后的操作
      })
    }
    else {
      callDeleteAPI() // 直接调用删除操作
    }
  }

  // 返回可供外部使用的函数和状态
  return {
    modalVisible,
    modalAction,
    modalTitle,
    modalLoading,
    handleAdd,
    handleDelete,
    handleEdit,
    handleView,
    handleSave,
    modalForm,
    modalFormRef,
  }
}