import {test as base} from '@playwright/test';
import {EnvironmentHelper} from "./fixtures/env";
import {DescopeHelper} from "./fixtures/descope";
import {BackendHelper} from "./fixtures/backend";
import {TenantHelper} from "./fixtures/tenant";
import {UserHelper} from "./fixtures/user";

interface AdditionalFixtures {
    backend: BackendHelper,
    descope: DescopeHelper
    env: EnvironmentHelper
    tenant: TenantHelper
    user: UserHelper
}

export const test = base.extend<AdditionalFixtures>({

    backend: async ({env}, use) => {
        await use(new BackendHelper(env))
    },

    descope: async ({}, use) => {
        const descope = new DescopeHelper();
        await use(descope)
        await descope.cleanup()
    },

    env: new EnvironmentHelper(),

    tenant: async ({env, descope}, use) => {
        const tenant = new TenantHelper(env, descope);
        await use(tenant)
        await tenant.cleanup()
    },

    user: async ({descope}, use) => {
        const user = new UserHelper(descope)
        await use(user)
        await user.cleanup()
    }
});
