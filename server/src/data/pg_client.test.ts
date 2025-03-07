import { convertObjectKeysToCamelCase } from "./pg_client.js"

describe("convertObjectKeysToCamelCase", () => {
    test("should convert object keys from snake_case to camelCase", () => {
        const input  = { first_name: "John", last_name: "Doe", is_active: true }
        const result = convertObjectKeysToCamelCase(input)
        expect(result).toEqual({ firstName: "John", lastName: "Doe", isActive: true })
    })

    test("should handle empty objects without errors", () => {
        const input  = {}
        const result = convertObjectKeysToCamelCase(input)
        expect(result).toEqual({})
    })

    test("should keep the same values while transforming keys", () => {
        const input  = { user_id: 123, is_verified: false }
        const result = convertObjectKeysToCamelCase(input)
        expect(result).toEqual({ userID: 123, isVerified: false })
    })

    test("should not modify keys that are already in camelCase", () => {
        const input  = { userId: 123, firstName: "Alice" }
        const result = convertObjectKeysToCamelCase(input)
        expect(result).toEqual({ userID: 123, firstName: "Alice" })
    })

    test("should skip transformation for non-string keys (numeric)", () => {
        const input  = { 1: "one", 2: "two" }
        const result = convertObjectKeysToCamelCase(input)
        expect(result).toEqual({ 1: "one", 2: "two" })
    })
})
