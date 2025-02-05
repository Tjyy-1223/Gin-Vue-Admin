import { createApp } from 'vue'
import App from './App.vue'

// unocss
import 'uno.css'
import '@unocss/reset/tailwind.css'

const app = createApp(App);
app.mount('#app')