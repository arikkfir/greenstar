import { BaseConfig } from "./base-config.ts"

class Config extends BaseConfig {

    public readonly jobName: string

    public readonly podName: string

    public readonly podNamespace: string

    /**
     * Creates a new configuration instance
     *
     * Loads all required environment variables for Bank Yahav scraper.
     * Throws an error if any required environment variable is missing.
     */
    constructor() {
        super()
        this.podName      = this.requireStringEnvVar("POD_NAME")
        this.podNamespace = this.requireStringEnvVar("POD_NAMESPACE")
        this.jobName      = this.requireStringEnvVar("JOB_NAME")
    }
}

export const k8sConfig = new Config()
