## 第七章 gin-blog-admin 首页搭建

### 7.1 初始化项目

**进入终端并使用命令行：** 

1. pnpm create vite -> gin-blog-admin -> vue -> ts -> 初始化完毕

2. 安装 vue 核心依赖：cd gin-blog-admin -> pnpm i
3. 运行项目 : pnpm run dev
4. 项目起步完成

------

**之后，将下面的内容放在 package.json 之后， 执行 pnpm i，完成相关配置：**

```
{
  "name": "gin-blog-front",
  "private": true,
  "version": "0.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite --open",
    "build": "vue-tsc -b && vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "@iconify/vue": "^4.3.0",
    "@unocss/reset": "^65.4.3",
    "@vueuse/core": "^12.5.0",
    "axios": "^1.7.9",
    "dayjs": "^1.11.13",
    "easy-typer-js": "^2.1.0",
    "highlight.js": "^11.11.1",
    "marked": "^15.0.6",
    "nprogress": "^0.2.0",
    "pinia": "^2.3.1",
    "pinia-plugin-persistedstate": "^4.2.0",
    "rollup-plugin-visualizer": "^5.14.0",
    "unocss": "^65.4.3",
    "v3-infinite-loading": "^1.3.2",
    "vite-plugin-compression": "^0.5.1",
    "vue": "^3.5.13",
    "vue-router": "^4.5.0",
    "vue3-danmaku": "^1.6.1"
  },
  "devDependencies": {
    "@egoist/tailwindcss-icons": "^1.8.2",
    "@iconify-json/mdi-light": "^1.2.2",
    "@iconify/json": "^2.2.302",
    "@iconify/tailwind": "^1.2.0",
    "@mdi/font": "^7.4.47",
    "@types/node": "^22.12.0",
    "@vitejs/plugin-vue": "^5.2.1",
    "@vue/tsconfig": "^0.7.0",
    "autoprefixer": "^10.4.20",
    "postcss": "^8.5.1",
    "sass-embedded": "^1.83.4",
    "tailwindcss": "^3.4.17",
    "typescript": "~5.6.2",
    "vite": "^6.0.5",
    "vue-tsc": "^2.2.0"
  }
}
```

**postcss.config.js**

```
export default {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
}
```

**tsconfig.app.json**

```json
{
  "extends": "@vue/tsconfig/tsconfig.dom.json",
  "compilerOptions": {
    "tsBuildInfoFile": "./node_modules/.tmp/tsconfig.app.tsbuildinfo",

    /* Linting */
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    "noUncheckedSideEffectImports": true,
    "lib": ["es6", "dom", "es2017"],
    "target": "es5",
    "module": "ESNext",
    "moduleResolution": "node",
    "jsx": "preserve",
    "allowJs": true,
    "baseUrl": "./",
    "paths": {
      "@/*": ["src/*"]  // 将 '@' 映射到 src 目录
    },
    "types": ["node","pinia"],
    "esModuleInterop": true,
    "skipLibCheck": true,
  },
  "include": ["src/**/*.ts", "src/**/*.tsx", "src/**/*.vue"]
}
```

**tailwind.config.js 中使用**

```javascript
/** @type {import('tailwindcss').Config} */

const { iconsPlugin, getIconCollections } = require('@egoist/tailwindcss-icons');
export default {
  content: ["index.html", "./src/**/*.{html,js,ts,jsx,tsx,vue}"],
  theme: {
    extend: {
      transitionDuration: {
        '500': '500ms',
      },
    },
  },
  plugins: [
    iconsPlugin({
      collections: getIconCollections(['mdi', 'lucide']), // 或使用 "all" 来使用全部图标
    }),
  ],
}
```

**修改 vite.config.ts 如下：**

```typescript
import path from 'path'
import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import unocss from 'unocss/vite'
import viteCompression from 'vite-plugin-compression'
import { visualizer } from 'rollup-plugin-visualizer'

export default defineConfig((configEnv) => {
  const env = loadEnv(configEnv.mode, process.cwd())

  return {
    base: env.VITE_PUBLIC_PATH || '/',
    resolve: {
      alias: {
        '@': path.resolve(__dirname, 'src'),
        '~': path.resolve(process.cwd()),
      },
    },
    plugins: [
      vue(),
      unocss(),
      viteCompression({ algorithm: 'gzip' }),
      visualizer({ open: false, gzipSize: true, brotliSize: true }),
    ],
    server: {
      host: '0.0.0.0',
      port: 3000,
      open: false,
      proxy: {
        '/api': {
          target: env.VITE_SERVER_URL,
          changeOrigin: true,
        },
      },
    },
    // https://cn.vitejs.dev/guide/api-javascript.html#build
    build: {
      chunkSizeWarningLimit: 1024, // chunk 大小警告的限制（单位kb）
    },
    esbuild: {
      drop: ['debugger'], // console
    },
  }
})
```

根目录下创建 uno.config.js

```javascript
import {
  defineConfig,
  presetIcons,
  presetTypography,
  presetUno,
  transformerDirectives,
  transformerVariantGroup,
} from 'unocss'

export default defineConfig({
  shortcuts: [
    ['f-c-c', 'flex justify-center items-center'],
  ],
  presets: [
    presetUno(),
    presetIcons({ warn: true }),
    presetTypography(),
  ],
  transformers: [
    transformerDirectives(),
    transformerVariantGroup(),
  ],
})

```

**然后，可以将 main.ts 编写如下：**

```typescript
import { createApp } from 'vue'
import App from './App.vue'

// unocss
import 'uno.css'
import '@unocss/reset/tailwind.css'

const app = createApp(App);
app.mount('#app')
```

之后，就可以开始编写对应的 vue 代码了

使用 pnpm run dev 查看目前的效果





### 7.2 主页搭建以及相关 store \ api

首先复制三个配置文件到根目录下，启动时会设置好相关的后端请求路径配置:

+ .env
+ .env.development
+ .env.production

