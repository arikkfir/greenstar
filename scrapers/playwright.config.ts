import { defineConfig, devices } from "@playwright/test"

import * as dotenv from "dotenv"

dotenv.config()

export default defineConfig({
    testDir: "./src/scrapers",
    testMatch: /.+\/scraper\.ts$/,
    fullyParallel: true,
    forbidOnly: true,
    retries: 0,
    workers: 1,
    reporter: [
        ["list", { printSteps: true }],
        ["html", { open: "never" }],
    ],
    use: {
        contextOptions: {
            reducedMotion: "reduce",
        },
        headless: !!process.env.KUBERNETES_SERVICE_HOST,
        launchOptions: {
            args: ['--no-sandbox'],
        },
        screenshot: "on",
        trace: "on",
        video: "on",
    },
    expect: {
        timeout: 1000 * 15,
    },
    reportSlowTests: null,
    timeout: 1000 * 60 * 30,
    projects: [
        {
            name: "Google Chrome",
            use: {
                ...devices["Desktop Chrome"],
                channel: "chrome",
            },
        },
    ],
})
