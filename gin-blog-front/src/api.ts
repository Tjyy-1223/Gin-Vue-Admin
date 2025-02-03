import { baseRequest, request } from '@/utils/http'

export default {
  login: (data = {}) => baseRequest.post('/login', data),
  logout: () => baseRequest.get('/logout'),

  /** 首页文章列表 */
  getArticles: (params: any) => request.get('/article/list', { params }),
  /** 文章搜索 */
  searchArticles: (params = {}) => request.get('/article/search', { params }),


  // ! 需要 Token 的接口
  /** 根据 token 获取当前用户信息 */
  getUser: () => request.get('/user/info', { needToken: true }),
}