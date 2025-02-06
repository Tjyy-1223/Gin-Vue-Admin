import { defineStore } from 'pinia'

// 用户全局变量
export const useUserStore = defineStore('user', {
    state: () => ({
        userInfo: {
            id: null,
            nickname: '',
            avatar: '',
            intro: '',
            website: '',
            // roles: [], // TODO: 后端返回 roles
        },
    }),
    getters: {
        userId: state => state.userInfo.id,
        nickname: state => state.userInfo.nickname,
        intro: state => state.userInfo.intro,
        website: state => state.userInfo.website,
        // avatar: state => convertImgUrl(state.userInfo.avatar),
        // roles: state => state.userInfo.roles,
    },
    actions: {
    },
})
