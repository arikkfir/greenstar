import {expect} from '@playwright/test';
import {test} from '../src/test'
import {User} from "../src/fixtures/descope";

const defaultExpectTimeoutIntervals = {
    intervals: [1000, 2000, 5000, 10_000],
    timeout: 60_000,
}

type ExpectToPassOptions = { timeout?: number, intervals?: number[] }

async function expectAndReturn<T>(f: () => T | Promise<T>, options: ExpectToPassOptions = defaultExpectTimeoutIntervals): Promise<T> {
    let result: T | Promise<T> | undefined
    await expect(async () => result = await f()).toPass(options)
    if (result === undefined) {
        throw new Error("result is undefined")
    } else {
        return result
    }
}

test('available', async ({page, env, tenant, user}) => {
    await expect(async () => tenant.create()).toPass(defaultExpectTimeoutIntervals)
    const u: User = await expectAndReturn(async () => await user.create({
        tenantRoles: [{tenantId: tenant.id, roleNames: ["Admin"]}],
    }))

    await user.login(page.context(), u)
    await page.goto(env.appURL(tenant.id));
    await expect(page).toHaveTitle("GreenSTAR");
    await expect(page.getByText("Homepage")).toBeVisible()
});
