import {test} from '@playwright/test';

test('dashboard', async ({page}) => {
    await page.goto('/');
    await page.waitForLoadState('networkidle')
    await page.locator('main div:has-text("Dashboard")').waitFor({state: "visible"})
});
