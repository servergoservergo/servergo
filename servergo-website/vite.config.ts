import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  base: '/',
  server: {
    // 在开发模式下正确处理单页应用路由
    port: 3000,
    cors: true,
  },
  build: {
    // 输出目录
    outDir: 'dist',
    // 静态资源文件生成时的路径前缀
    assetsDir: 'assets',
    // 是否生成源码映射文件
    sourcemap: false,
  }
})
