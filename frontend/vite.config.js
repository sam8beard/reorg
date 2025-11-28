import { defineConfig } from 'vite';

export default defineConfig({
	build: { 
		manifest: "manifest.json",
		outDir: '../cmd/server/dist',
	},
})

