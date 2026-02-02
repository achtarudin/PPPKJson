import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig(({ command, mode }) => {
  // Determine API base URL based on build mode
  const API_BASE_URL = mode === 'production' 
    ? 'https://pppk-json.cutbray.tech/api/v1'
    : 'http://localhost:8080/api/v1';

  return {
    plugins: [react()],
    define: {
      // Make API_BASE_URL available globally
      __API_BASE_URL__: JSON.stringify(API_BASE_URL),
    },
    server: {
      port: 5173,
      open: true,
    },
    build: {
      outDir: 'dist',
    }
  };
});
