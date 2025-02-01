import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
// 清除默认样式
import "./style/reset.scss"
import "./style/style.css"

let app = createApp(App);
app.mount('#app')
