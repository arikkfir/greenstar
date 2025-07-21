export function required<T>(value: T | null | undefined, message: string): T {
    if (value == null) {
        throw new Error(message)
    }
    return value
}
