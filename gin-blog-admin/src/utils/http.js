import axios from 'axios'
import { useAuthStore } from '@/store'

// 创建 axios 实例
export const request = axios.create({
  baseURL: import.meta.env.VITE_BASE_API,  // 设置请求的基础 URL（从环境变量中读取）
  timeout: 12000, // 设置请求的超时时间为 12000 毫秒（12 秒）
})

// 请求拦截器
request.interceptors.request.use(
  // 请求成功拦截
  (config) => {
    // 判断该请求是否需要携带 Token，如果不需要，则直接返回 config
    if (config.noNeedToken) {
      return config
    }

    // 获取 token（通常在 store 中存储）
    const { token } = useAuthStore()

    // 如果 token 存在，则将 token 添加到请求头的 Authorization 字段
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }

    // 返回修改后的请求配置
    return config
  },
  // 请求失败拦截
  (error) => {
    // 请求发生错误时，直接返回拒绝的 Promise
    return Promise.reject(error)
  },
)

// 响应拦截器
request.interceptors.response.use(
  // 响应成功拦截
  (response) => {
    // 业务信息：从响应中提取数据
    const responseData = response.data
    const { code, message, data } = responseData

    // 判断响应中的业务状态码，如果不等于 0 说明业务失败
    if (code !== 0) {  // 这里的 `0` 是后端约定的成功状态码
      // 如果存在 data，且 message 和 data 不相等，则拼接错误信息
      if (data && message !== data) {
        window.$message.error(`${message} ${data}`)  // 使用 UI 库弹出错误消息
      } else {
        window.$message.error(message)  // 使用 UI 库弹出错误消息
      }

      // 在控制台打印错误信息，便于调试
      console.error(responseData)

      const authStore = useAuthStore() // 获取认证状态

      // 如果返回的 code 为 1201，则说明 Token 存在问题，跳转到登录页
      if (code === 1201) {
        authStore.toLogin() // 跳转到登录页面
        return
      }

      // 如果返回的 code 为 1202、1203 或 1207，说明 Token 过期或者被强制下线，执行强制下线操作
      if (code === 1202 || code === 1203 || code === 1207) {
        authStore.forceOffline() // 强制用户下线
        return
      }

      // 返回 Promise.reject，表示响应失败
      return Promise.reject(responseData)
    }

    // 如果业务成功，返回响应数据
    return Promise.resolve(responseData)
  },
  // 响应失败拦截
  (error) => {
    // 响应失败，错误通常是由网络问题或服务器问题引起的
    const responseData = error.response?.data
    const { message, data } = responseData

    // 如果 HTTP 状态码是 500（服务器错误），显示服务端异常
    if (error.response.status === 500) {
      if (message && data) {
        window.$message.error(`${message} ${data}`)
      } else {
        window.$message.error('服务端异常')  // 显示服务端异常错误信息
      }
    }

    // 返回 Promise.reject，表示响应失败
    return Promise.reject(error)
  },
)
