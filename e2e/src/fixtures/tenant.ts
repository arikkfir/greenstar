import {DescopeHelper} from "./descope";
import {BackendHelper} from "./backend";

export class TenantHelper {
    private readonly tenantAdminAccessKey: string = process.env.DESCOPE_TENANT_ADMIN_ACCESS_KEY || '';
    private _id?: string
    private _displayName?: string
    private deleteOnCleanup: boolean = false

    constructor(private readonly backend: BackendHelper, private readonly descope: DescopeHelper) {
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
        const response = await this.backend.createTenant(jwt, `test-${random}`, `Test (${random})`);
        this._id = response.id
        this._displayName = response.displayName
        this.deleteOnCleanup = true
    }

    async load(id: string) {
        await this.cleanup()
        const jwt = await this.descope.exchangeAccessKey(this.tenantAdminAccessKey)
        const response = await this.backend.loadTenant(jwt, id);
        this._id = response.id
        this._displayName = response.displayName
        this.deleteOnCleanup = false
    }

    async cleanup() {
        if (this._id && this.deleteOnCleanup) {
            const jwt = await this.descope.exchangeAccessKey(this.tenantAdminAccessKey)
            await this.backend.deleteTenant(jwt, this._id);
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
