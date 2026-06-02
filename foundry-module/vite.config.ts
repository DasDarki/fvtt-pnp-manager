import { defineConfig } from 'vite'

export default defineConfig({
  publicDir: 'public',
  build: {
    outDir: 'dist',
    emptyOutDir: true,
    minify: false,
    sourcemap: true,
    lib: {
      entry: 'src/main.ts',
      formats: ['es'],
      fileName: () => 'aetherwright.js',
    },
    rollupOptions: {
      output: { entryFileNames: 'aetherwright.js' },
    },
  },
})
