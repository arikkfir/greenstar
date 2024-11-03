import {defineConfig, ProxyOptions, ViteDevServer} from 'vite'
import react from '@vitejs/plugin-react'
import * as core from "express-serve-static-core";
import express from 'express'
import pkg from './package.json'

const app = express()
app.use((req, res, next) => {
    res.setHeader('Server', 'Greenstar');
    res.setHeader('X-Greenstar-Version', pkg.version);
    next();
})
app.get('/healthz', (req, res) => {
    res.status(200).send('OK').end()
})

const proxy: Record<string, string | ProxyOptions> = {
    '/healthz': {} // proxy our /api route to nowhere
}

function healthPlugin(app: core.Express) {
    return {
        name: 'health',
        config() {
            return {
                server: {proxy},
                preview: {proxy}
            }
        },
        configureServer(server: ViteDevServer) {
            server.middlewares.use(app)
        }
    }
}

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [react(), healthPlugin(app)],
})
