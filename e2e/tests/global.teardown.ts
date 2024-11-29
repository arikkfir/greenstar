import {test as setup} from '@playwright/test';
import Descope from "@descope/node-sdk";

setup('teardown', async ({}) => {
    if (!process.env.DESCOPE_PROJECT_ID) {
        throw new Error("Missing DESCOPE_PROJECT_ID")
    } else if (!process.env.DESCOPE_MANAGEMENT_KEY_TOKEN) {
        throw new Error("Missing DESCOPE_MANAGEMENT_KEY_TOKEN")
    }

    const descope = Descope({
        projectId: process.env.DESCOPE_PROJECT_ID,
        managementKey: process.env.DESCOPE_MANAGEMENT_KEY_TOKEN,
    });

    await descope.management.user.delete(process.env.TEST_USER_LOGIN_ID);
});
