const locationIQURL = "https://us1.locationiq.com/v1/reverse"
const locationIQAccessToken = import.meta.env.VITE_LOCATION_IQ_ACCESS_TOKEN;

export interface LocationIQResponse {
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

export async function getLocation(position: GeolocationPosition): Promise<LocationIQResponse> {
    const urlParams = new URLSearchParams();
    urlParams.set("key", locationIQAccessToken)
    urlParams.set("lat", position.coords.latitude + "")
    urlParams.set("lon", position.coords.longitude + "")
    urlParams.set("format", "json")
    return await fetch(`${locationIQURL}?${urlParams}`, {method: 'GET'})
        .then(resp => {
            if (!resp.ok) {
                throw new Error("Unable to obtain your location: " + resp.statusText + " (" + resp.status + ")")
            } else {
                return resp.json()
            }
        })
}