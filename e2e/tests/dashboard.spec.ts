import {test} from '@playwright/test';

test('dashboard', async ({page}) => {
    await page.goto('/');
    await page.waitForLoadState('networkidle')
    await page.locator('h3:has-text("Dashboard")').waitFor({state: "visible"})
});
