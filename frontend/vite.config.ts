import { defineConfig } from 'vite';
import solidPlugin from 'vite-plugin-solid';

const backendHost = process.env.VITE_BACKEND_URL || 'localhost';

export default defineConfig({
  plugins: [solidPlugin()],
  server: {
    port: 3000,
    proxy: {
      '/media': {
        target: `http://${backendHost}:8080`,
        changeOrigin: true, 
        secure: false,
      }
    }
  },
  build: {
    target: 'esnext',
    outDir: 'dist',
  },
});