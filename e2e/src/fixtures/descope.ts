import DescopeClientFactory, {AuthenticationInfo, SdkResponse} from '@descope/node-sdk';
import * as DescopeSDK from "@descope/core-js-sdk";
import * as _descope_core_js_sdk from "@descope/core-js-sdk";

export type AttributesTypes = string | boolean | number;

export type AssociatedTenant = {
    tenantId: string;
    roleNames: string[];
};

type UserTenant = {
    tenantId: string;
    roleNames?: string[];
};
type GenerateOTPForTestResponse = {
    loginId: string;
    code: string;
};

export interface User {
    id: string
    loginIds: string[]
    name?: string
    email?: string
    phone?: string
    picture?: string
    globalRoles?: string[]
    tenantRoles?: UserTenant[]
}

interface DescopeClient {
    exchangeAccessKey: (accessKey: string) => Promise<AuthenticationInfo>
    management: {
        user: {
            createTestUser: (
                loginId: string,
                email?: string,
                phone?: string,
                name?: string,
                roles?: string[],
                userTenants?: AssociatedTenant[],
                customAttributes?: Record<string, AttributesTypes>,
                picture?: string
            ) => Promise<SdkResponse<DescopeSDK.UserResponse>>,
            delete: (loginId: string) => Promise<SdkResponse<never>>
            deleteAllTestUsers: () => Promise<SdkResponse<never>>
            generateOTPForTestUser: (
                deliveryMethod: "email" | "sms" | "whatsapp",
                loginId: string
            ) => Promise<SdkResponse<GenerateOTPForTestResponse>>
        }
    }
    otp: {
        verify: {
            email: (loginId: string, code: string) => Promise<SdkResponse<_descope_core_js_sdk.JWTResponse & {
                refreshJwt?: string;
                cookies?: string[];
            }>>
        }
    }
}

export class DescopeHelper {
    private readonly descopeProjectID: string = process.env.DESCOPE_PROJECT_ID || '';
    private readonly descopeManagementKey: string = process.env.DESCOPE_MANAGEMENT_KEY || '';
    private readonly client: DescopeClient;

    constructor() {
        if (!this.descopeProjectID) {
            throw new Error('DESCOPE_PROJECT_ID is not defined');
        } else if (!this.descopeManagementKey) {
            throw new Error('DESCOPE_MANAGEMENT_KEY is not defined');
        }
        this.client = DescopeClientFactory({
            projectId: this.descopeProjectID,
            managementKey: this.descopeManagementKey,
        })
    }

    async cleanup() {
        // no-op
    }

    async exchangeAccessKey(accessKey: string): Promise<string> {
        const authInfo = await this.client.exchangeAccessKey(accessKey);
        return authInfo.jwt
    }

    async createUser(p: {
        loginId: string,
        name?: string,
        email: string,
        phone?: string,
        picture?: string,
        globalRoles: string[],
        tenantRoles: AssociatedTenant[],
    }): Promise<User> {
        const response = await this.client.management.user.createTestUser(
            p.loginId, p.email, p.phone, p.name, p.globalRoles, p.tenantRoles, undefined, p.picture
        )
        if (!response.ok) {
            throw new Error(`Failed to create user: ${JSON.stringify(response.error)}`)
        } else if (!response.data) {
            throw new Error(`No data in user creation response: ${JSON.stringify(response)}`)
        } else {
            return {
                id: response.data.userId,
                loginIds: response.data.loginIds,
                name: response.data.name,
                email: response.data.email,
                phone: response.data.phone,
                picture: response.data.picture,
                globalRoles: response.data.roleNames,
                tenantRoles: response.data.userTenants
            }
        }
    }

    async deleteAllTestUsers() {
        await this.client.management.user.deleteAllTestUsers()
    }

    async deleteUser(loginId: string): Promise<void> {
        const response = await this.client.management.user.delete(loginId)
        if (!response.ok) {
            throw new Error(`Failed to create user: ${JSON.stringify(response.error)}`)
        } else if (!response.data) {
            throw new Error(`No data in user creation response: ${JSON.stringify(response)}`)
        }
    }

    async generateOTP(loginId: string): Promise<string> {
        const response = await this.client.management.user.generateOTPForTestUser("email", loginId)
        if (!response.ok) {
            throw new Error(`Failed to generate OTP: ${JSON.stringify(response.error)}`)
        } else if (!response.data) {
            throw new Error(`No data in OTP generation response: ${JSON.stringify(response)}`)
        } else {
            return response.data.code
        }
    }

    async verifyOTP(loginId: string, code: string): Promise<{ jwt: { session: string, refresh?: string } }> {
        const response = await this.client.otp.verify.email(loginId, code)
        if (!response.ok) {
            throw new Error(`Failed to generate OTP: ${JSON.stringify(response.error)}`)
        } else if (!response.data) {
            throw new Error(`No data in OTP generation response: ${JSON.stringify(response)}`)
        } else {
            return {
                jwt: {
                    session: response.data.sessionJwt,
                    refresh: response.data.refreshJwt,
                },
            }
        }
    }
}
