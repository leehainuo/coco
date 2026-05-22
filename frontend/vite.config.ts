import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import path from 'path'

export default defineConfig(() => {
  const backendPort = process.env.VITE_BACKEND_VERSION === 'v3' ? 8081 : 8080
  const backendUrl = `http://localhost:${backendPort}`

  return {
    plugins: [tailwindcss(), vue()],
    base: './',
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },
    server: {
      proxy: {
        '^/docs/(openapi|swagger|config)\\.json$': {
          target: backendUrl,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/docs/, '/dev'),
        },
        '/api': {
          target: backendUrl,
          changeOrigin: true,
        },
      },
    },
    build: {
      outDir: '../internal/assets/dist',
      emptyOutDir: true,
      rollupOptions: {
        output: {
          manualChunks: undefined,
          inlineDynamicImports: true,
        },
      },
      cssCodeSplit: false,
      minify: 'esbuild',
      sourcemap: false,
    },
  }
})
