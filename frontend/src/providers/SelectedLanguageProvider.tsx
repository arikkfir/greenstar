import { ReactNode, useCallback, useContext, useLayoutEffect, useMemo, useState } from "react"
import { LoadedSelectedLanguageStateCtx } from "../contexts/LoadingState.ts"
import { Language, LanguagesContext, SelectedLanguageContext } from "../contexts/Languages.ts"
import { delayFactor } from "./util.ts"

export function SelectedLanguageProvider({ children }: { children: ReactNode }) {
    const { setLoaded: markAsLoaded } = useContext(LoadedSelectedLanguageStateCtx)
    const storageKey                  = "language"
    const languages: Language[]       = useContext(LanguagesContext)
    const [language, setLanguage]     = useState<Language | undefined>()

    useLayoutEffect(() => {
        const storedCode = localStorage.getItem(storageKey)
        if (storedCode) {
            const language = languages.find((l) => l.code === storedCode)
            if (language) {
                setLanguage(language)
            }
        }
        window.setTimeout(markAsLoaded, Math.random() * delayFactor)
    }, [setLanguage, markAsLoaded])

    const saveLanguage = useCallback(
        (l: Language) => {
            setLanguage(l)
            localStorage.setItem(storageKey, l.code)
        },
        [setLanguage],
    )

    const value = useMemo(
        () => ({ language, setLanguage: saveLanguage }),
        [language, saveLanguage],
    )

    return (
        <SelectedLanguageContext.Provider value={value}>{children}</SelectedLanguageContext.Provider>
    )
}
