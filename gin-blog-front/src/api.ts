import { baseRequest, request } from '@/utils/http'

export default {
  login: (data = {}) => baseRequest.post('/login', data),
  logout: () => baseRequest.get('/logout'),

  /** 获取页面 */
  getPageList: () => request.get('/page'),
  /** 首页文章列表 */
  getArticles: (params: any) => request.get('/article/list', { params }),
  /** 文章搜索 */
  searchArticles: (params = {}) => request.get('/article/search', { params }),
  /** 文章归档 */
  getArchives: (params = {}) => request.get('/article/archive', { params }),


  /** 菜单列表 */
  getCategorys: () => request.get('/category/list'),
  /** 标签列表 */
  getTags: () => request.get('/tag/list'),

  // ! 需要 Token 的接口
  /** 根据 token 获取当前用户信息 */
  getUser: () => request.get('/user/info', { needToken: true }),
}