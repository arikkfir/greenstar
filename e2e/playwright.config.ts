import {defineConfig, devices} from '@playwright/test';
import dotenv from 'dotenv';

dotenv.config();

export default defineConfig({
    testDir: './tests',
    fullyParallel: true,
    expect: {
        timeout: 1000 * 60 * 5,
    },
    timeout: 1000 * 60 * 5,
    globalTimeout: 1000 * 60 * 10,
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
        ignoreHTTPSErrors: true,
        locale: 'en-IL',
        screenshot: 'on',
        timezoneId: 'Asia/Jerusalem',
        trace: 'on',
        video: 'on',
    },
    projects: [
        {
            name: 'chromium',
            use: {...devices['Desktop Chrome']},
        },
    ],
});
