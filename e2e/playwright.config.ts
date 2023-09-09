import {defineConfig, devices} from '@playwright/test';

/**
 * Read environment variables from file.
 * https://github.com/motdotla/dotenv
 */
// require('dotenv').config();

/**
 * See https://playwright.dev/docs/test-configuration.
 */
export default defineConfig({
    testDir: './tests',
    timeout: 30 * 1000,
    expect: {
        timeout: 5000
    },
    fullyParallel: true,
    forbidOnly: !!process.env.CI,
    retries: 0,
    workers: 1,
    reporter:
        process.env.CI
            ? [
                ['dot'],
                ['github'],
                ['@skilbourn/playwright-report-summary', {outputFile: 'custom-summary.txt'}],
                ['html', {open: 'never'}],
            ]
            : [
                ['list', {printSteps: true}],
                ['html', {open: 'never'}],
            ],
    use: {
        actionTimeout: 0,
        /* Base URL to use in actions like `await page.goto('/')`. */
        // baseURL: 'http://127.0.0.1:3000',
        ignoreHTTPSErrors: true,
        trace: 'on',
    },
    projects: [
        {
            name: 'chromium',
            use: {...devices['Desktop Chrome']},
        },
        // {
        //   name: 'firefox',
        //   use: { ...devices['Desktop Firefox'] },
        // },
        // {
        //   name: 'webkit',
        //   use: { ...devices['Desktop Safari'] },
        // },
        // {
        //   name: 'Google Chrome',
        //   use: { ...devices['Desktop Chrome'], channel: 'chrome' },
        // },
        // {
        //   name: 'Mobile Chrome',
        //   use: { ...devices['Pixel 7'] },
        // },
        // {
        //   name: 'Mobile Safari',
        //   use: { ...devices['iPhone 13'] },
        // },
        // {
        //   name: 'Microsoft Edge',
        //   use: { ...devices['Desktop Edge'], channel: 'msedge' },
        // },
    ],
});
