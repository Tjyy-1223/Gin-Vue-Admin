import { createApp } from 'vue'
import App from './App.vue'

// custom style
// import './style/index.css'
import './style/common.css'
// import './style/animate.css'

// unocss
import 'uno.css'
import '@unocss/reset/tailwind.css'

import { router } from './router'
import store from './store'

const app = createApp(App);
app.use(router); // 注册路由
app.use(store); // 注册pinia
app.mount('#app')


console.log(import.meta.env.VITE_APP_TITLE)