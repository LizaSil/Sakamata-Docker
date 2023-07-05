import { defineConfig } from "vite"

export default defineConfig({
  base: "./",
  server: {
    port: 3000,
    strict: true,
  },
  preview: {
    port: 3000,
  },
})
