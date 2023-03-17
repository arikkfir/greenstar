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
            jwtid: "test-session-id",
            subject: "test|1",
            mutatePayload: true,
        });

        const session = {
            "claims": payload,
            "token": {
                "access_token": "bogus_access_token",
                "token_type": "Bearer",
                "refresh_token": "bogus_refresh_token",
                "expiry": DateTime.now().plus({minutes: 15}).toISO(),
            },
            "tenant": "test",
            "permissions": [
                "greenstar.auth.getUserInfo",
                "greenstar.auth.getUserInfo:mock",
                "greenstar.admin.createTenant",
            ],
            "mockUserInfo": {
                "email": "jack@ryan.com",
                "family_name": "Ryan",
                "given_name": "Jack",
                "hd": "test",
                "id": "test|1",
                "name": "Jack Ryan"
            }
        }

        // Save the session to Redis
        const redisClient = createClient();
        await redisClient.connect();
        await redisClient.set('session:test-session-id', JSON.stringify(session));
        await redisClient.disconnect();

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
        await expect(page.getByTestId("dashboard-title")).toHaveText('Dashboard here.');
    })
})
