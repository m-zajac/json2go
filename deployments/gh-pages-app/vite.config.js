import { defineConfig } from 'vite'
import preact from '@preact/preset-vite'
import { copyFileSync, mkdirSync } from 'fs'
import { resolve } from 'path'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    preact(),
    {
      name: 'copy-wasm-files',
      closeBundle() {
        // Copy WASM files from build/web to the output directory
        const wasmSource = resolve(__dirname, '../../build/web')
        const outDir = resolve(__dirname, '../../build/gh-pages')
        
        try {
          mkdirSync(outDir, { recursive: true })
          copyFileSync(`${wasmSource}/json2go.wasm`, `${outDir}/json2go.wasm`)
          copyFileSync(`${wasmSource}/wasm_exec.js`, `${outDir}/wasm_exec.js`)
          console.log('✓ Copied WASM files to build directory')
        } catch (err) {
          console.warn('Warning: Could not copy WASM files:', err.message)
        }
      }
    }
  ],
  base: './',
  build: {
    outDir: '../../build/gh-pages',
    emptyOutDir: true,
  },
})
