import {test as setup} from '@playwright/test';
import * as crypto from "crypto";
import Descope from "@descope/node-sdk";
import * as fs from "node:fs";
import * as path from "node:path";
import {authFile} from "../util/globals";

setup('setup', async ({page, baseURL}) => {
    const testUserLoginID = crypto.randomBytes(5).toString("hex");
    const testUserEmail = `arikkfir+${testUserLoginID}@gmail.com`
    process.env.TEST_USER_LOGIN_ID = testUserLoginID;
    process.env.TEST_USER_EMAIL = testUserEmail;

    if (!process.env.DESCOPE_PROJECT_ID) {
        throw new Error("Missing DESCOPE_PROJECT_ID")
    } else if (!process.env.DESCOPE_MANAGEMENT_KEY_TOKEN) {
        throw new Error("Missing DESCOPE_MANAGEMENT_KEY_TOKEN")
    }

    const descope = Descope({
        projectId: process.env.DESCOPE_PROJECT_ID,
        managementKey: process.env.DESCOPE_MANAGEMENT_KEY_TOKEN,
    });

    await descope.management.user.createTestUser(testUserLoginID, `arikkfir+${testUserLoginID}@gmail.com`);

    const magicLink = await descope.management.user.generateMagicLinkForTestUser(
        "email",
        testUserLoginID,
        "https://acme.app.greenstar.test"
    );
    const token = magicLink.data.link.split("?t=")[1];
    const auth = await descope.magicLink.verify(token);

    await page.goto('/');
    await page.evaluate(
        ([ds, dsr]) => {
            window.localStorage.setItem("DS", ds);
            window.localStorage.setItem("DSR", dsr);
        },
        [auth.data.sessionJwt, auth.data?.refreshJwt]
    );

    if (!fs.existsSync(path.dirname(authFile))) {
        fs.mkdirSync(path.dirname(authFile))
    }
    fs.writeFileSync(authFile, JSON.stringify(await page.context().storageState()))
});
