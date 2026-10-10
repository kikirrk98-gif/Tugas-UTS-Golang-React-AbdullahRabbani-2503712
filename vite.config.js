import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// Saat development, permintaan /api diteruskan ke server Go di port 8080.
export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: { "/api": "http://localhost:8080" },
  },
});
