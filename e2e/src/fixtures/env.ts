export class EnvironmentHelper {
    public readonly name: string = process.env.ENV_NAME || '';

    constructor() {
        if (!this.name) {
            throw new Error('ENV_NAME is not defined');
        }
    }

    get domainSuffix(): string {
        return `${this.name}.greenstar.kfirs.com`;
    }

    get apiURL(): string {
        return `https://api.${this.domainSuffix}`;
    }

    get apiGraphURL(): string {
        return this.apiURL + '/query';
    }

    appURL(tenant: string): string {
        return `https://${tenant}.${this.domainSuffix}`;
    }
}
