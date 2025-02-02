import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
// 清除默认样式
import "./style/reset.scss"
import "./style/style.css"

import { router } from './router'
import store from './store'

const app = createApp(App);
app.use(router); // 注册路由
app.use(store); // 注册pinia
app.mount('#app')


console.log(import.meta.env.VITE_APP_TITLE)