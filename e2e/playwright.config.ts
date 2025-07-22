import { defineConfig, devices } from "@playwright/test"

import * as dotenv from "dotenv"
import * as path from "path"

dotenv.config({ path: path.resolve("../.env") })

export default defineConfig({
    testDir: "./tests",
    fullyParallel: true,
    forbidOnly: !!process.env.CI,
    workers: process.env.CI ? 1 : undefined,
    reporter: process.env.CI
        ? [
            [ "list", { printSteps: true } ],
            [ "github" ],
            [
                "monocart-reporter", {
                name: "End-to-end Report",
                outputFile: "./monocart-report/index.html",
            },
            ],
        ]
        : [
            [ "list", { printSteps: true } ],
            [ "html", { open: "never" } ],
        ],
    use: {
        baseURL: "https://acme.app.greenstar.test",
        screenshot: "on",
        trace: "on",
        video: "on",
        ignoreHTTPSErrors: true,
    },
    timeout: 120_000,
    globalTimeout: 120_000,
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
