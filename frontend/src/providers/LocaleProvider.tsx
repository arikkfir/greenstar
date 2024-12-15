import { createContext, useEffect, useState } from "react"
import { CurrenciesByCountryCode } from "../util/currencies.ts"
import { getLocation, LocationIQResponse } from "../services/location_iq.ts"

export interface Locale {
    language: string
    currency: string
    location?: LocationIQResponse
    resolved: boolean
    error?: Error
}

function getDefaultLocale(): Locale {
    try {
        const storedLocale = localStorage.getItem("locale")
        if (storedLocale) {
            return Object.assign({ resolved: true }, JSON.parse(storedLocale))
        } else {
            return {
                language: navigator.language,
                currency: "USD",
                resolved: false,
            }
        }
    } catch (error) {
        console.error("Failed to load location from local storage:", error)
        return {
            language: navigator.language,
            currency: "USD",
            resolved: false,
            error: error instanceof Error ? error : new Error("Failed generating default locale: " + error),
        }
    }
}

// TODO: verify we're using context correctly here - there are two default values (one on the context and one in the provider)

export const LocaleContext = createContext<Locale>({
    language: navigator.language,
    currency: "USD",
    resolved: false,
    error: new Error("Geolocation not initialized yet"),
})

export function LocaleProvider({ children }: { children: any }) {
    const [locale, setLocale] = useState<Locale>(getDefaultLocale())
    useEffect(() => {
        if (!navigator.geolocation) {
            setLocale({
                language: navigator.language,
                currency: "USD",
                resolved: false,
                error: new Error("Geolocation is not supported by your browser"),
            })
            return
        }

        if (!locale.resolved) {
            const getPositionPromise = new Promise<GeolocationPosition>((resolve, reject) =>
                navigator.geolocation.getCurrentPosition(
                    (position) => resolve(position),
                    (error) => reject(new Error(error.message)),
                    { timeout: 1000 * 60 * 5, maximumAge: 1000 * 60 * 60 * 4 },
                ),
            )
            getPositionPromise
                .then((position) => getLocation(position))
                .then(
                    (location): Locale => ({
                        language: navigator.language,
                        location: location,
                        currency: CurrenciesByCountryCode[location.address.country_code.toUpperCase()] || "USD",
                        resolved: true,
                        error: undefined,
                    }),
                )
                .then((locale) => {
                    localStorage.setItem("locale", JSON.stringify(locale))
                    return locale
                })
                .then(setLocale)
                .catch((e) => console.error("Failed to obtain your locale: ", e))
        }
    }, [locale, setLocale])

    return <LocaleContext.Provider value={locale}>{children}</LocaleContext.Provider>
}
