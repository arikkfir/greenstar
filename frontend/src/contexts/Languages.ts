import { createContext } from "react"

export type Language = {
    code: string,
    name: string,
}

export const LanguagesContext = createContext<Language[]>([])

export type SelectedLanguage = {
    language?: Language
    setLanguage(l: Language): void
}

export const SelectedLanguageContext = createContext<SelectedLanguage>({
    setLanguage: () => {
        throw new Error("setLanguage called on uninitialized SelectedLanguageContext.")
    },
})
