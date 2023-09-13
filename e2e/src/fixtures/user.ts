import {AssociatedTenant, DescopeHelper, User} from "./descope";
import {BrowserContext} from "@playwright/test";

/** Default name for the session cookie name / local storage key */
const SESSION_TOKEN_KEY = 'DS';
/** Default name for the refresh local storage key */
const REFRESH_TOKEN_KEY = 'DSR';

export class UserHelper {
    private usersToDelete: User[] = []
    private loggedInUsers: { context: BrowserContext, user: User }[] = []

    constructor(private readonly descope: DescopeHelper) {
    }

    async create(u: { picture?: string, globalRoles?: string[], tenantRoles?: AssociatedTenant[] }): Promise<User> {
        const random = Math.random().toString(36).substring(7);
        const loginId = "playwright-" + random
        const email = `support+${loginId}@kfirs.com`
        const name = `GreenSTAR Playwright (${random})`
        const user = await this.descope.createUser(
            {
                loginId,
                name,
                email,
                picture: u.picture,
                globalRoles: u.globalRoles || [],
                tenantRoles: u.tenantRoles || [],
            }
        );
        this.usersToDelete.push(user)
        return user
    }

    async cleanup() {
        await Promise.all(this.loggedInUsers.map(async i => await i.context.clearCookies()))
        this.loggedInUsers = []

        await Promise.all(this.usersToDelete.map(u => this.descope.deleteUser(u.loginIds[0])))
        this.usersToDelete = []
    }

    async login(context: BrowserContext, user: User) {
        const loginId = user.loginIds[0]
        const otpCode = await this.descope.generateOTP(loginId)
        const response = await this.descope.verifyOTP(loginId, otpCode);

        const domain = ".descope.com"
        const path = "/"
        const expires = Date.now() / 1000 + 60 * 10
        const secure = true
        const sameSite = "None"
        await context.addCookies([
            {name: SESSION_TOKEN_KEY, value: response.jwt.session, domain, expires, path, secure, sameSite},
        ])

        if (response.jwt.refresh) {
            await context.addCookies([
                {name: REFRESH_TOKEN_KEY, value: response.jwt.refresh, domain, expires, path, secure, sameSite},
            ])
        }
        this.loggedInUsers.push({context, user})
    }

    async logout(context: BrowserContext) {
        this.loggedInUsers = this.loggedInUsers.filter(i => i.context != context)
        await context.clearCookies()
    }
}
