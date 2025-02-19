import { defineConfig, loadEnv } from 'vite';
import vue from '@vitejs/plugin-vue';
import unocss from 'unocss/vite'
import path from 'path';
import viteCompression from 'vite-plugin-compression'
import { visualizer } from 'rollup-plugin-visualizer'


export default defineConfig((configEnv) => {
  const env = loadEnv(configEnv.mode, process.cwd())
  console.log('Alias @:', path.resolve(__dirname, 'src'));
  return {
    base: env.VITE_PUBLIC_PATH || '/',
    resolve: {
      alias: {
        "@": path.resolve(__dirname, 'src'),
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
      port: 3333,
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
      chunkSizeWarningLimit: 1024, // chunk 大小警告的限制 (单位 kb)
    },
    esbuild: {
      drop: ['debugger'], // console
    },
  }
})