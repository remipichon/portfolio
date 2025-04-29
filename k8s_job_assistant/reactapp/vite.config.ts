import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  base: "./",
  build: {
    outDir: "build",
  },
  server: {
    port: 3000,
    strictPort: true,
    proxy: {
      '/list': 'http://localhost:8080',
      '/run': 'http://localhost:8080',
      '/kill': 'http://localhost:8080',
    },
  }
});
