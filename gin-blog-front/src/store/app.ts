import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', {
    state: () => ({
        searchFlag: false,
        loginFlag: false,
        registerFlag: false,
        collapsed: false, // 侧边栏折叠（移动端）

        page_list: [], // 页面数据
        // TODO: 优化
        blogInfo: {
            article_count: 0,
            category_count: 0,
            tag_count: 0,
            view_count: 0,
            user_count: 0,
        },
        blog_config: {
            website_name: 'CodingHome',
            website_author: 'Tjyy',
            website_intro: 'coding is coding',
            website_avatar: 'https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png',
            website_notice: 'hello, world',
            qq: '',
            github: '',
            gitee: '',
            website_createtime: '2025-2-1',
            website_record: 'CodingHome'
        },
    }),
    getters: {
        isMobile: () => !!navigator.userAgent.match(/(phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone)/i),
        articleCount: state => state.blogInfo.article_count ?? 0,
        categoryCount: state => state.blogInfo.category_count ?? 0,
        tagCount: state => state.blogInfo.tag_count ?? 0,
        viewCount: state => state.blogInfo.view_count ?? 0,
        pageList: state => state.page_list ?? [],
        blogConfig: state => state.blog_config,
    },
    actions: {
        setCollapsed(flag: any) { this.collapsed = flag },
        setLoginFlag(flag: any) { this.loginFlag = flag },
        setRegisterFlag(flag: any) { this.registerFlag = flag },
        setSearchFlag(flag: any) { this.searchFlag = flag },
        // async getBlogInfo() {
        //     try {
        //         const resp = await api.getHomeData()
        //         if (resp.code === 0) {
        //             this.blogInfo = resp.data
        //             this.blog_config = resp.data.blog_config
        //             this.blog_config.website_avatar = convertImgUrl(this.blog_config.website_avatar)
        //         }
        //         else {
        //             return Promise.reject(resp)
        //         }
        //     }
        //     catch (err) {
        //         return Promise.reject(err)
        //     }
        // },

        // async getPageList() {
        //     const resp = await api.getPageList()
        //     if (resp.code === 0) {
        //         this.page_list = resp.data
        //         this.page_list?.forEach(e => (e.cover = convertImgUrl(e.cover)))
        //     }
        // },
    },
})
