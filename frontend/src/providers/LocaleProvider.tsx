import {createContext, useEffect, useMemo, useState} from "react";
import {CurrenciesByCountryCode} from "../util/currencies.ts";

const locationIQAccessToken = "pk.8bf13229466b44319e00651d2d6e15db";

interface LocationIQResponse {
    display_name: string
    address: {
        name: string
        house_number: string
        road: string
        city: string
        state: string
        postcode: string
        country: string
        country_code: string
    }
}

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
            const urlParams = new URLSearchParams();
            urlParams.set("key", locationIQAccessToken)
            urlParams.set("lat", position.coords.latitude + "")
            urlParams.set("lon", position.coords.longitude + "")
            urlParams.set("format", "json")
            fetch("https://us1.locationiq.com/v1/reverse?" + urlParams, {method: 'GET'})
                .then(resp => {
                    if (!resp.ok) {
                        setLocation(undefined)
                        setLocationError(new Error("Unable to obtain your location: " + resp.statusText + " (" + resp.status + ")"))
                    } else {
                        return resp.json()
                    }
                })
                .then((data: LocationIQResponse) => {
                    setLocation(data)
                    setLocationError(undefined)
                })
                .catch(e => {
                    setLocation(undefined)
                    setLocationError(e)
                })
        }
    }, [position, positionError, setLocation, setLocationError, locationIQAccessToken]);

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
