import { expect, Locator, Page } from '@playwright/test';
import jwt from "jsonwebtoken";
import {DateTime} from "luxon";

export class LoginPage {
    readonly page: Page;
    readonly getStartedLink: Locator;
    readonly gettingStartedHeader: Locator;
    readonly pomLink: Locator;
    readonly tocList: Locator;

    constructor(page: Page) {
        this.page = page;
        this.getStartedLink = page.locator('a', { hasText: 'Get started' });
        this.gettingStartedHeader = page.locator('h1', { hasText: 'Installation' });
        this.pomLink = page.locator('li', { hasText: 'Guides' }).locator('a', { hasText: 'Page Object Model' });
        this.tocList = page.locator('article div.markdown ul > li > a');
    }

    async login() {
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
        await this.page.context().addCookies([{
            name: 'greenstar-session',
            value: token,
            path: '/',
            domain: "localhost",
            expires: DateTime.now().plus({minutes: 15}).toUnixInteger(),
            httpOnly: true,
            secure: false,
        }])

        // Navigate to the dashboard, and expect our session cookie to be accepted
        await this.page.goto('http://localhost');
        await expect(this.page.getByTestId("dashboard-title")).toHaveText('Dashboard here.');
    }
}
