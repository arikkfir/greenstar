import { createContext, useEffect, useState } from "react"
import { CurrenciesByCountryCode } from "../util/currencies.ts"

const url = "https://api.geoapify.com/v1/geocode/reverse"
const apiKey = import.meta.env.VITE_GEOAPIFY_API_KEY

export interface Locale {
    language: string
    currency: string
    country: string
    countryCode: string
}

const dummyLocale = {
    language: "en",
    currency: "USD",
    country: "United States",
    countryCode: "US",
}

export const LocaleContext = createContext<Locale>(dummyLocale)

export function LocaleProvider({ children }: { children: any }) {
    const [locale, setLocale] = useState<Locale>(dummyLocale)

    useEffect(() => {
        try {
            const storedLocale = localStorage.getItem("locale")
            if (storedLocale) {
                const parsedLocation = JSON.parse(storedLocale)

                // Ignore older versions of stored location
                if (!parsedLocation.location && !parsedLocation.error && !parsedLocation.resolved) {
                    setLocale(parsedLocation)
                }
            }
        } catch (error) {
            console.error("Failed to load location from local storage:", error)
        }

        if (!apiKey) {
            console.warn("No GeoAPIfy API key provided. Using dummy locale.")
            return
        }

        const getPositionPromise = new Promise<GeolocationPosition>((resolve, reject) =>
            navigator.geolocation.getCurrentPosition(
                (position) => resolve(position),
                (error) => reject(new Error(error.message)),
                {
                    timeout: 1000 * 60 * 5,
                    maximumAge: 1000 * 60 * 60 * 4,
                },
            ),
        )
        getPositionPromise
            .then((position) => {
                const urlParams = new URLSearchParams()
                urlParams.set("lat", position.coords.latitude + "")
                urlParams.set("lon", position.coords.longitude + "")
                urlParams.set("type", "country")
                urlParams.set("limit", "1")
                urlParams.set("format", "json")
                urlParams.set("apiKey", apiKey)
                return fetch(`${url}?${urlParams}`, { method: "GET" })
            })
            .then((resp) => {
                if (!resp.ok) {
                    throw new Error(`Unable to obtain your location: ${resp.statusText} (${resp.status})`)
                } else {
                    return resp.json()
                }
            })
            .then((json) => json.results)
            .then((results) => results[0])
            .then(
                (location): Locale => ({
                    language: navigator.language,
                    currency: CurrenciesByCountryCode[location.country_code.toUpperCase()] || "USD",
                    country: location.country,
                    countryCode: location.country_code,
                }),
            )
            .then((locale) => {
                localStorage.setItem("locale", JSON.stringify(locale))
                return locale
            })
            .then(setLocale)
            .catch((e) => console.error("Failed to obtain your locale: ", e))
    }, [setLocale])

    return <LocaleContext.Provider value={locale}>{children}</LocaleContext.Provider>
}
