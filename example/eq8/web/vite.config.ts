import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  base: "./",
  plugins: [react()],
  build: {
    outDir: "editor",
    emptyOutDir: false,
    rollupOptions: {
      output: {
        entryFileNames: "editor.js",
        chunkFileNames: "editor-[name].js",
        assetFileNames: (assetInfo) => {
          if (assetInfo.name?.endsWith(".css")) {
            return "editor.css";
          }
          return "[name][extname]";
        },
      },
    },
  },
});
