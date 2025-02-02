import { baseRequest, request } from '@/utils/http'

export default {
    /** 首页文章列表 */
    getArticles: (params: any) => request.get('/article/list', { params }),
}