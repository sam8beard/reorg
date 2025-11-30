import { defineConfig, loadEnv } from 'vite';
const distPath = process.env.DIST_DIR; 

export default defineConfig({
	server: {
		cors: { 
			origin: ['http://localhost:5137'],
			methods: ['GET', 'POST', 'DELETE', 'PUT'],
			allowedHeaders: ['Content-Type'],
		},
		proxy: {
			// Proxy for /api requests
			'/api': {
				target: 'http://localhost:8080',
				// Change header of each request to point 
				// to target
				changeOrigin: true,
				// Remove prefix of request for correct 
				// resolution on the backend server
				rewrite: (path) => path.replace(/^\/api/, ''),
			},
		},
	},

	build: { 
		manifest: "manifest.json",
		outDir: '../cmd/server/dist',
	},
})

