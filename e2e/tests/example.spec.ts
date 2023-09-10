import {expect, test} from '@playwright/test';

// TODO: https://docs.descope.com/knowledgebase/testing/cypresstesting/

const ENV_NAME = process.env.ENV_NAME;
if (!ENV_NAME) {
    throw new Error('ENV_NAME is not defined');
}

test('has title', async ({page}) => {
    const url = `https://acme.${ENV_NAME}.greenstar.kfirs.com/`;
    console.info(`Navigating to: ${url}`)

    await page.goto(url, {
        timeout: 1000 * 60,
    });

    await expect(page).toHaveTitle("GreenSTAR", {
        timeout: 1000 * 60,
    });
});

