import {DescopeHelper} from "./descope";
import {BackendHelper} from "./backend";
import {EnvironmentHelper} from "./env";

export class TenantHelper extends BackendHelper {
    private readonly tenantAdminAccessKey: string = process.env.DESCOPE_TENANT_ADMIN_ACCESS_KEY || '';
    private _id?: string
    private _displayName?: string
    private deleteOnCleanup: boolean = false

    constructor(env: EnvironmentHelper, private readonly descope: DescopeHelper) {
        super(env)
        if (!this.tenantAdminAccessKey) {
            throw new Error('DESCOPE_TENANT_ADMIN_ACCESS_KEY is not defined');
        }
    }

    async create() {
        if (this._id) {
            throw new Error("tenant already created or loaded (perhaps you forgot to call 'cleanup()'?)")
        }
        const random = Math.random().toString(36).substring(7);
        const jwt = await this.descope.exchangeAccessKey(this.tenantAdminAccessKey)
        const response = await this.sdk(jwt).createTenant({
            tenantID: `test-${random}`,
            displayName: `Test (${random})`,
        })
        this._id = response.createTenant.id
        this._displayName = response.createTenant.displayName
        this.deleteOnCleanup = true
    }

    async load(id: string) {
        await this.cleanup()
        const jwt = await this.descope.exchangeAccessKey(this.tenantAdminAccessKey)
        const response = await this.sdk(jwt).loadTenant({tenantID: id})
        if (!response.tenant) {
            throw new Error(`Empty tenant received`)
        }
        this._id = response.tenant.id
        this._displayName = response.tenant.displayName
    }

    async cleanup() {
        if (this._id && this.deleteOnCleanup) {
            const jwt = await this.descope.exchangeAccessKey(this.tenantAdminAccessKey)
            await this.sdk(jwt).deleteTenant({tenantID: this._id})
            this.deleteOnCleanup = false
        }
        this._id = undefined
        this._displayName = undefined
    }

    get id(): string {
        if (this._id) {
            return this._id
        } else {
            throw new Error("tenant not created or loaded")
        }
    }

    get displayName(): string {
        if (this._displayName) {
            return this._displayName
        } else {
            throw new Error("tenant not created or loaded")
        }
    }
}
