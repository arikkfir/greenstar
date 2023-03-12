import {defineConfig, devices} from '@playwright/test';

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
    reporter: 'list',
    use: {
        actionTimeout: 0,
        baseURL: 'http://localhost:3000',
        trace: 'on',
    },
    projects: [
        {
            name: 'chromium',
            use: {...devices['Desktop Chrome']},
        },
        // {
        //     name: 'firefox',
        //     use: {...devices['Desktop Firefox']},
        // },
        // {
        //     name: 'webkit',
        //     use: {...devices['Desktop Safari']},
        // },
    ],
});
