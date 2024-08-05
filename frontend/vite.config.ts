import {defineConfig} from 'vite';
import react from '@vitejs/plugin-react';

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [react()],
    server: {
        proxy: {
            // 代理所有以 /api 开头的请求
            '/api': {
                target: 'http://127.0.0.1:8000', // 目标服务器地址
                changeOrigin: true, // 是否更改请求头中的 origin 字段
                rewrite: (path) => path.replace(/^\/api/, ''), // 去掉请求路径中的 /api 前缀
            },
        },
    },
});
