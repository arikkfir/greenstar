import {Locator, Page} from "@playwright/test";
import {randomInt, randomUUID} from "node:crypto";

export class LoginPage {
    readonly testUserEmail: string = `arikkfir+${randomInt(5)}@gmail.com`
    readonly testUserPassword: string = randomUUID()
    readonly page: Page;
    readonly emailField: Locator;
    readonly passwordField: Locator;
    readonly loginButton: Locator;
    readonly dashboardHeader: Locator;

    constructor(page: Page) {
        this.page = page;
        this.emailField = page.getByLabel('Email');
        this.passwordField = page.getByLabel('Password').nth(0);
        this.loginButton = page.getByRole('button').filter({hasText: 'LOGIN'});
        this.dashboardHeader = page.locator('h3:has-text("Dashboard")');
    }

    async loginWithTestUser() {
        await this.page.context().clearCookies({domain: 'acme.app.greenstar.test'})
        await this.page.context().clearCookies({domain: 'descope.com'})
        await this.page.goto('https://acme.app.greenstar.test/');
        await this.page.getByLabel('Email').fill(this.testUserEmail);
        await this.page.getByLabel('Password').nth(0).fill(this.testUserPassword);
        await this.page.getByRole('button').filter({hasText: 'LOGIN'}).click()
        await this.dashboardHeader.waitFor({state: "visible"})
    }
}
