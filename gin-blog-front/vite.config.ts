import { defineConfig, loadEnv } from 'vite';
import vue from '@vitejs/plugin-vue';
import path from 'path';

// https://vitejs.dev/config/
// export default defineConfig({

//   plugins: [vue()],
//   resolve: {
//     alias: {
//       '@': path.resolve(__dirname, 'src'),  // 将 '@' 映射到 src 目录
//     },
//   },
// });
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
    plugins: [vue()],
    server: {
      host: '0.0.0.0',
      port: 3333,
      open: false,
      proxy: {
        '/api': {
          target: env.VITE_BACKEND_URL,
          changeOrigin: true,
        },
      },
    },
  }
})