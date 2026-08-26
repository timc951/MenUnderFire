import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'

// Set VITE_DEBUG_BUILD=true to disable minification for easier debugging
const isDebugBuild = process.env.VITE_DEBUG_BUILD === 'true'

export default defineConfig({
  plugins: [react()],
  server: {
    host: '0.0.0.0',
    port: 8001,
  },
  build: {
    // Disable minification for debug builds
    minify: isDebugBuild ? false : 'esbuild',
    // Generate source maps for debug builds
    sourcemap: isDebugBuild ? true : false,
  },
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: './src/__tests__/setup.ts',
    coverage: {
      provider: 'v8',
      reporter: ['text', 'html', 'lcov'],
      exclude: ['node_modules', 'src/__tests__', 'src/__mocks__'],
      thresholds: {
        branches: 80,
        functions: 80,
        lines: 80,
        statements: 80,
      },
    },
  },
})
