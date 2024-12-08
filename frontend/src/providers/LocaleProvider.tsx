import {createContext, useEffect, useState} from "react";
import {CurrenciesByCountryCode} from "../util/currencies.ts";
import {getLocation, LocationIQResponse} from "../services/location_iq.ts";

export interface Locale {
    language: string
    currency: string
    location?: LocationIQResponse,
    error?: Error
}

function getDefaultLocale(): Locale {
    try {
        const storedLocale = localStorage.getItem("locale");
        if (storedLocale) {
            return JSON.parse(storedLocale)
        } else {
            return {
                language: navigator.language,
                currency: 'USD',
            }
        }
    } catch (error) {
        console.error('Failed to load location from local storage:', error);
        return {
            language: navigator.language,
            currency: 'USD',
            error: error instanceof Error ? error : new Error("Failed generating default locale: " + error),
        }
    }
}
const defaultLocale = getDefaultLocale()

export const LocaleContext = createContext<Locale>(defaultLocale);

export function LocaleProvider({children}: { children: any }) {
    const [locale, setLocale] = useState<Locale>(defaultLocale);

    useEffect(() => {
        if (!navigator.geolocation) {
            setLocale(defaultLocale)
            console.warn("Geolocation is not supported by your browser.")
            return
        }

        const getPositionPromise = new Promise<GeolocationPosition>(
            (resolve, reject) =>
                navigator.geolocation.getCurrentPosition(
                    (position) => resolve(position),
                    (error) => reject(new Error(error.message)),
                    {timeout: 5000, maximumAge: 1000 * 60 * 60 * 4},
                ),
        )

        getPositionPromise
            .then(position => {
                console.info("Obtained geolocation position: ", position)
                return getLocation(position)
            })
            .then((location): Locale => {
                console.info("Obtained location from coordinates: ", location)
                return {
                    language: navigator.language,
                    location: location,
                    currency: CurrenciesByCountryCode[location.address.country_code.toUpperCase()] || 'USD',
                    error: undefined,
                }
            })
            .then(locale => {
                console.info("Stored locale in local storage: ", locale)
                localStorage.setItem("locale", JSON.stringify(locale))
                return locale
            })
            .then(setLocale)
            .catch(e => console.error("Failed to obtain your locale: ", e))
    }, []);

    return (
        <LocaleContext.Provider value={locale}>
            {children}
        </LocaleContext.Provider>
    );
}
