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
    // 如果有 actions，可以在这里添加
  },
  persist: {
    key: 'gvb_blog_user',
    paths: ['token'], // 确保 token 是直接在 state 的根路径下
  }
});