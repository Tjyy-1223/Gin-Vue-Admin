import { unref } from "vue";
import { defineStore } from "pinia";

export const useAuthStore = defineStore('auth', {
    persist: {
        key: 'gvb_admin_auth',
        paths: ['token'],
    },
    state: () => ({
        token: null,
    }),
    actions: {
        setToken(token) {
            this.token = token
        },
    },
    getters: {

    }
})
