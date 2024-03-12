import {defineConfig, loadEnv} from 'vite'
import vue from '@vitejs/plugin-vue'


export default ({mode})=>defineConfig({
  plugins: [vue()],
  server: {
    // 配置代理
    proxy: {
      '/api': {
        target: loadEnv(mode, process.cwd()).VITE_PROXY_API,
        ws: true,
        changeOrigin: true,
      },
    },
  },
})
