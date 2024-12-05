import {createContext, useEffect, useMemo, useState} from "react";
import {CurrenciesByCountryCode} from "../util/currencies.ts";
import {getLocation, LocationIQResponse} from "../services/location_iq.ts";

export interface Locale {
    language: string
    position?: GeolocationPosition
    location?: LocationIQResponse,
    currency?: string
    error?: Error
}

export const LocaleContext = createContext<Locale>({language: navigator.language});

export function LocaleProvider({children}: { children: any }) {
    const [position, setPosition] = useState<GeolocationPosition | undefined>();
    const [positionError, setPositionError] = useState<Error | undefined>();
    const [location, setLocation] = useState<LocationIQResponse | undefined>();
    const [locationError, setLocationError] = useState<Error | undefined>();
    const [currency, setCurrency] = useState<string | undefined>();
    const [currencyError, setCurrencyError] = useState<Error | undefined>();

    useEffect(() => {
        if (navigator.geolocation) {
            navigator.geolocation.getCurrentPosition(
                (position) => setPosition(position),
                (error) => setPositionError(new Error(error.message)));
        } else {
            setPosition(undefined)
            setPositionError(new Error('Unable to obtain your geolocation'));
        }
    }, [setPosition, setPositionError]);

    useEffect(() => {
        if (positionError || !position) {
            setLocation(undefined)
            setLocationError(undefined)
        } else {
            getLocation(position)
                .then(data => {
                    setLocation(data)
                    setLocationError(undefined)
                })
                .catch(e => {
                    setLocation(undefined)
                    setLocationError(e)
                })
        }
    }, [position, positionError, setLocation, setLocationError]);

    useEffect(() => {
        if (locationError || !location) {
            setCurrency(undefined)
            setCurrencyError(undefined)
        } else {
            const countryCode = location.address.country_code.toUpperCase()
            const currencyCode = CurrenciesByCountryCode[countryCode]
            if (currencyCode) {
                setCurrency(currencyCode)
                setCurrencyError(undefined)
            } else {
                setCurrency(undefined)
                setCurrencyError(new Error("Could not discover your preferred currency (defaulting to USD)"))
            }
        }
    }, [location, locationError, setCurrency, setCurrencyError]);

    const locale = useMemo((): Locale => {
        return {
            language: navigator.language,
            position: position,
            location: location,
            currency: currency,
            error: positionError || locationError || currencyError
        }
    }, [navigator.language, position, location, currency, positionError, locationError, currencyError])

    return (
        <LocaleContext.Provider value={locale}>
            {children}
        </LocaleContext.Provider>
    );
}
