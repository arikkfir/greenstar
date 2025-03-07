import type { Config } from "jest"

const config: Config = {
    preset: "ts-jest",
    testEnvironment: "node",
    testMatch: [ "**/__tests__/**/*.ts", "**/?(*.)+(spec|test).ts" ],
    extensionsToTreatAsEsm: ['.ts'],
    moduleNameMapper: {
        '^(\\.{1,2}/.*)\\.js$': '$1',
    },
    transform: {
        '^.+\\.ts$': ['ts-jest', {
            useESM: true,
        }]
    }
}

export default config
