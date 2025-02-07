import { ref } from 'vue'

/**
 * 可复用的表单对象
 * 该函数用于封装表单的状态管理、验证逻辑以及表单规则的配置，可以在多个组件中复用。
 * @param {any} initForm 表单初始值，用来初始化表单数据
 * @returns {object} 返回一个包含表单引用、表单模型、验证函数和表单规则的对象
 */
export function useForm(initForm = {}) {
  const formRef = ref(null) // 表单的引用，用来获取表单实例进行操作（如验证）
  const formModel = ref({ ...initForm }) // 表单模型，保存表单数据，初始值为传入的 `initForm`

  /**
   * 表单验证函数
   * 该函数会触发表单验证，验证通过返回 `true`，否则返回 `false`
   * @returns {boolean} 验证是否通过
   */
  const validation = async () => {
    try {
      // 调用表单实例的 `validate` 方法进行表单验证
      await formRef.value?.validate()
      return true // 验证成功，返回 true
    }
    catch (error) {
      return false // 验证失败，返回 false
    }
  }

  // 表单字段的验证规则
  const rules = {
    required: {
      required: true,  // 必填项
      message: '此为必填项',  // 错误提示信息
      trigger: ['blur', 'change'],  // 在 `blur`（失去焦点）或 `change`（值改变）时触发验证
    },
  }

  // 返回包含表单引用、表单模型、验证方法和表单规则的对象
  return { formRef, formModel, validation, rules }
}
