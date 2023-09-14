import {expect} from '@playwright/test';
import {test} from '../src/test'

test('available', async ({page, env, tenant, user}) => {
    await tenant.create()
    const u = await user.create({
        tenantRoles: [
            {tenantId: tenant.id, roleNames: ["Admin"]}
        ],
    });
    await user.login(page.context(), u)
    await page.goto(env.appURL(tenant.id));
    await expect(page).toHaveTitle("GreenSTAR");
    await expect(page.getByText("Homepage")).toBeVisible()
});
