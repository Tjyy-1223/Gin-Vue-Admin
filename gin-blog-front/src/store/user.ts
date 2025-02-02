import { defineStore } from 'pinia';

// 自定义 state 的类型
interface UserState {
  userInfo: {
    id: string;
    nickname: string;
    avatar: string;
    website: string;
    intro: string;
    email: string;
    articleLikeSet: any[];
    commentLikeSet: any[];
  };
  token: string | null;
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    userInfo: {
      id: '',
      nickname: '',
      avatar: 'https://example.com/default-avatar.png',
      website: '',
      intro: '',
      email: '',
      articleLikeSet: [],
      commentLikeSet: [],
    },
    token: null,
  }),
  getters: {
    userId: state => state.userInfo.id ?? '',
    nickname: state => state.userInfo.nickname ?? '',
    avatar: state => state.userInfo.avatar ?? 'https://example.com/default-avatar.png',
    website: state => state.userInfo.website ?? '',
    intro: state => state.userInfo.intro ?? '',
    email: state => state.userInfo.email ?? '',
    articleLikeSet: state => state.userInfo.articleLikeSet || [],
    commentLikeSet: state => state.userInfo.commentLikeSet || [],
  },
  actions: {
    setToken(token:string) {
      this.token = token
    },
    resetLoginState() {
      this.$reset()
    },
    // async logout() {
    //   await api.logout()
    //   this.$reset()
    // },
    // async getUserInfo() {
    //   if (!this.token) {
    //     return
    //   }
    //   try {
    //     const resp = await api.getUser()
    //     if (resp.code === 0) {
    //       const data = resp.data
    //       this.userInfo = {
    //         id: data.id,
    //         nickname: data.nickname,
    //         avatar: data.avatar ? convertImgUrl(data.avatar) : 'https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png',
    //         website: data.website,
    //         intro: data.intro,
    //         email: data.email,
    //         articleLikeSet: data.article_like_set.map(e => +e),
    //         commentLikeSet: data.comment_like_set.map(e => +e),
    //       }
    //       return Promise.resolve(resp.data)
    //     }
    //     else {
    //       return Promise.reject(resp)
    //     }
    //   }
    //   catch (error) {
    //     return Promise.reject(error)
    //   }
    // },
    // commentLike(commentId) {
    //   this.commentLikeSet.includes(commentId)
    //     ? this.commentLikeSet.splice(this.commentLikeSet.indexOf(commentId), 1)
    //     : this.commentLikeSet.push(commentId)
    // },
    // articleLike(articleId) {
    //   this.articleLikeSet.includes(articleId)
    //     ? this.articleLikeSet.splice(this.articleLikeSet.indexOf(articleId), 1)
    //     : this.articleLikeSet.push(articleId)
    // },
  },
  persist: true
  // persist: {
  //   key: 'gvb_blog_user',
  //   paths: ['token'], // 确保 token 是直接在 state 的根路径下
  // }
});