import { defineConfig } from "vite";

import tailwindcss from "@tailwindcss/vite";
import franken from "franken-ui/plugin-vite";

export default defineConfig({
  root: ".",
  build: {
    outDir: "dist",
    minify: "esbuild",
    sourcemap: true,
  },
  server: {
    port: 3000,
    host: true,
  },
  optimizeDeps: {
    include: ["alpinejs"],
  },
  plugins: [tailwindcss()],
});
