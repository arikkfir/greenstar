import { ReactNode, useContext, useLayoutEffect, useState } from "react"
import { LoadedLanguagesStateCtx } from "../contexts/LoadingState.ts"
import { Language, LanguagesContext } from "../contexts/Languages.ts"
import { delayFactor } from "./util.ts"

export function LanguagesProvider({ children }: { children: ReactNode }) {
    const { setLoaded: markAsLoaded } = useContext(LoadedLanguagesStateCtx)
    const [languages, setLanguages]   = useState<Language[]>([])

    useLayoutEffect(() => {
        const handleLanguageChange = () => {
            const resolver = new Intl.DisplayNames(navigator.languages, { type: "language" })
            setLanguages(navigator.languages.map((code) => ({
                code,
                name: resolver.of(code.includes("-") ? code.split("-")[0] : code) || code,
            })))
            window.setTimeout(markAsLoaded, Math.random() * delayFactor)
        }
        window.addEventListener("languagechange", handleLanguageChange)

        handleLanguageChange()

        return () => window.removeEventListener("languagechange", handleLanguageChange)
    }, [markAsLoaded, setLanguages])

    return (
        <LanguagesContext.Provider value={languages}>{children}</LanguagesContext.Provider>
    )
}
