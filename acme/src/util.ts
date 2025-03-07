import { v4 as uuid } from "uuid"
import RandomSeed from "random-seed"

class SeededUuidGenerator {
    private readonly rng: RandomSeed.RandomSeed

    constructor(seed?: string) {
        this.rng = RandomSeed.create(seed)
    }

    generateUUID(): string {
        const randomBytes = new Uint8Array(16)
        for (let i = 0; i < 16; i++) {
            randomBytes[i] = Math.floor(this.rng.random() * 256)
        }
        return uuid({ random: randomBytes })
    }

    generateNumber(): number {
        return this.rng.random()
    }
}

const generator = new SeededUuidGenerator("reproducible-seed")

export const randomUUID = () => generator.generateUUID()

export const randomNumber = () => generator.generateNumber()

export function splitCamelCase(input: string, acronyms: Set<string> = new Set([ "ID", "URL", "API" ])): string {
    const words = input
        .replace(/([a-z0-9])([A-Z])/g, "$1 $2") // Split between lowercase and uppercase
        .replace(/([A-Z]+)([A-Z][a-z])/g, "$1 $2") // Split acronyms from words (e.g., "APIName")
        .split(" ")
        .map((word) => {
            const upperWord = word.toUpperCase()
            return acronyms.has(upperWord) ? upperWord : word[0].toUpperCase() + word.slice(1).toLowerCase()
        })
    return words.join(" ")
}
