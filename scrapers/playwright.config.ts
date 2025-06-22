import { defineConfig, devices } from "@playwright/test"

import * as dotenv from "dotenv"

dotenv.config()

export default defineConfig({
    testDir: "./src",
    testMatch: /scrapers\/.+\/scraper\.ts/,
    fullyParallel: true,
    forbidOnly: false,
    workers: 1,
    reporter: [
        ["list", { printSteps: true }],
        ["html", { open: "never" }],
    ],
    use: {
        baseURL: "https://www.bank-yahav.co.il",
        // contextOptions: {
        //     geolocation: {
        //         longitude: 34.89242,
        //         latitude: 32.1869342,
        //     },
        //     locale: "he-IL",
        //     permissions: ["geolocation"],
        // },
        contextOptions: {
            reducedMotion: "reduce",
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
