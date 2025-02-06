import { createApp } from 'vue'
import App from './App.vue'
import { setupRouter } from './router'
import { setupStore } from './store'

// unocss
import 'uno.css'
import '@unocss/reset/tailwind.css'

const app = createApp(App);
setupStore(app); // 优先级最高
await setupRouter(app);
app.mount('#app')