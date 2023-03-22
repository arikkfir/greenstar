import {expect, test} from '@playwright/test';
import {config} from 'dotenv'
import jwt from 'jsonwebtoken';
import {DateTime} from 'luxon'
import * as path from "path";
import {createClient} from 'redis';

config({
    path: path.resolve(__dirname, '..', 'google-client-app.env')
});

test.describe('Login', () => {

    // test.beforeEach(async ({page}, testInfo) => {
    //     await page.context().tracing.start({ screenshots: true, snapshots: true });
    // })
    // test.afterEach(async ({page}, testInfo) => {
    //     await page.context().tracing.stop({ path: `trace/${testInfo.title}-trace.zip` });
    // })

    test('login', async ({page}, testInfo) => {
        expect(process.env.GOOGLE_CLIENT_SECRET).toBeDefined()

        const payload = {}
        let token = jwt.sign(payload, process.env.GOOGLE_CLIENT_SECRET, {
            algorithm: 'HS256',
            expiresIn: '15m',
            notBefore: 0,
            audience: [
                "greenstar.admin",
                "greenstar.auth",
                "greenstar.operations",
                "greenstar.public"
            ],
            issuer: "greenstar.auth",
            jwtid: "test:" + DateTime.now().toISO(),
            subject: "test|1",
            mutatePayload: true,
        });
        console.info("Token: ", token)

        // Set the cookie
        await page.context().addCookies([{
            name: 'greenstar-session',
            value: token,
            path: '/',
            domain: "localhost",
            expires: DateTime.now().plus({minutes: 15}).toUnixInteger(),
            httpOnly: true,
            secure: false,
        }])

        // Navigate to the dashboard, and expect our session cookie to be accepted
        await page.goto('http://localhost');
        console.info("Went to page")
        await expect(page.getByTestId("dashboard-title")).toHaveText('Dashboard here.');
    })
})
