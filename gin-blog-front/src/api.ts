import { baseRequest, request } from '@/utils/http'

export default {
  login: (data = {}) => baseRequest.post('/login', data),
  logout: () => baseRequest.get('/logout'),

  /** 获取页面 */
  getPageList: () => request.get('/page'),
  /** 首页文章列表 */
  getArticles: (params: any) => request.get('/article/list', { params }),
  /** 文章详情 */
  getArticleDetail: id => request.get(`/article/${id}`),
  /** 文章搜索 */
  searchArticles: (params = {}) => request.get('/article/search', { params }),
  /** 文章归档 */
  getArchives: (params = {}) => request.get('/article/archive', { params }),


  /** 菜单列表 */
  getCategorys: () => request.get('/category/list'),
  /** 标签列表 */
  getTags: () => request.get('/tag/list'),
  /** 评论列表 */
  getComments: (params = {}) => request.get('/comment/list', { params }),
  /** 评论回复列表 */
  getCommentReplies: (id, params = {}) => request.get(`/comment/replies/${id}`, { params }),

  // ! 需要 Token 的接口
  /** 根据 token 获取当前用户信息 */
  getUser: () => request.get('/user/info', { needToken: true }),
  /** 点赞评论 */
  saveLikeComment: id => request.get(`/comment/like/${id}`, { needToken: true }),
  /** 评论 */
  saveComment: data => request.post('/comment', data, { needToken: true }),
  /** 点赞文章 */
  saveLikeArticle: id => request.get(`/article/like/${id}`, { needToken: true }),
}