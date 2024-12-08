import {defineConfig, devices} from '@playwright/test';

import * as dotenv from 'dotenv';
import * as path from 'path';
import {authFile} from "./util/globals";

dotenv.config({path: path.resolve("../.env")});

export default defineConfig({
    testDir: './tests',
    fullyParallel: true,
    forbidOnly: !!process.env.CI,
    workers: process.env.CI ? 1 : undefined,
    reporter: process.env.CI
        ? [
            ['list', {printSteps: true}],
            ['github'],
            ['monocart-reporter', {
                name: "End-to-end Report",
                outputFile: "./monocart-report/index.html",
            }],
        ]
        : [
            ['list', {printSteps: true}],
            ['html', {open: 'never'}],
        ],
    use: {
        baseURL: 'https://acme.app.greenstar.test',
        contextOptions: {
            geolocation: { longitude: 34.89242, latitude: 32.1869342 },
            locale: 'he-IL',
            permissions: ['geolocation'],
        },
        screenshot: 'on',
        trace: 'on',
        video: 'on',
    },
    expect: {
        timeout: 1000 * 5,
    },
    timeout: 1000 * 60 * 5,
    projects: [
        {
            name: 'Setup',
            testMatch: /global\.setup\.ts/,
            teardown: 'Teardown'
        },
        {
            name: 'Teardown',
            testMatch: /global\.teardown\.ts/,
        },
        {
            name: 'Google Chrome',
            use: {
                ...devices['Desktop Chrome'],
                channel: 'chrome',
                storageState: authFile,
            },
            dependencies: ['Setup']
        },
    ],
});
