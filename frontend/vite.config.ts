import { defineConfig } from 'vite';
import solidPlugin from 'vite-plugin-solid';

export default defineConfig({
  plugins: [solidPlugin()],
  server: {
    port: 3000,
    proxy: {
      '/media': {
        target: 'http://localhost:8080',
        changeOrigin: true, 
        secure: false,
      }
    }
  },
  build: {
    target: 'esnext',
  },
});