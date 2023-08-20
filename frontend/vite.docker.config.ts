import {defineConfig} from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
    plugins: [react()],
    server: {
        cors: false,
        port: 80,
        strictPort: true,
        hmr: {
            clientPort: 443
        },
        host: "0.0.0.0",
    }
})
