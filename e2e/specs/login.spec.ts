import {test} from '@playwright/test';
import {config} from 'dotenv'
import * as path from "path";
import {LoginPage} from "../pom/login";

config({
    path: path.resolve(__dirname, '..', 'google-client-app.env')
});

test.describe('Login', () => {
    test('login', async ({page}, _) => {
        const loginPage = new LoginPage(page);
        await loginPage.login()
    })
})
