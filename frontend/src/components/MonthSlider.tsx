import { Slider } from "@mui/material"
import { DateTime } from "luxon"
import { useCallback, useMemo } from "react"

// Define the props for the MonthSlider component
export interface MonthSliderProps {
    minDate: DateTime;
    maxDate: DateTime;
    value: DateTime;
    onChange: (newValue: DateTime) => void;
    disabled?: boolean;
}

/**
 * MonthSlider component that displays a slider with ticks for months between min and max dates.
 *
 * @param minDate - The minimum date (Luxon DateTime representing a year and month)
 * @param maxDate - The maximum date (Luxon DateTime representing a year and month)
 * @param value - The current selected date (Luxon DateTime)
 * @param onChange - Callback function called when the value changes
 * @param disabled - Whether the slider is disabled
 */
export function MonthSlider({ minDate, maxDate, value, onChange, disabled = false }: MonthSliderProps) {
    // Calculate the number of months between min and max dates
    const totalMonths = useMemo(() => {
        // Ensure minDate is before maxDate
        if (minDate > maxDate) {
            console.warn("minDate is after maxDate, swapping values")
            return maxDate.diff(minDate, "months").months
        }
        return maxDate.diff(minDate, "months").months
    }, [ minDate, maxDate ])

    // Convert DateTime value to slider value (0 to totalMonths)
    const dateToSliderValue = useCallback(
        (date: DateTime): number => {
            // Ensure the date is within bounds
            if (date < minDate) {
                return 0
            }
            if (date > maxDate) {
                return totalMonths
            }
            return date.diff(minDate, "months").months
        },
        [ minDate, maxDate, totalMonths ],
    )

    // Convert slider value (0 to totalMonths) to DateTime
    const sliderValueToDate = useCallback(
        (sliderValue: number): DateTime => minDate.plus({ months: Math.round(sliderValue) }),
        [ minDate ],
    )

    // Generate marks for each month
    const marks = useMemo(
        () => {
            const result = []
            for (let i = 0; i <= totalMonths; i++) {
                const date = minDate.plus({ months: i })
                result.push({
                    value: i,
                    label: i % 3 === 0 || i == totalMonths ? date.toFormat("MM/yy") : "",
                })
            }
            return result
        },
        [ minDate, totalMonths ],
    )

    // Handle slider change
    const handleChange = (_event: Event, newValue: number | number[]) => {
        const numericValue = Array.isArray(newValue) ? newValue[0] : newValue
        onChange(sliderValueToDate(numericValue))
    }

    // Current slider value
    const sliderValue = dateToSliderValue(value)

    return (
        <Slider value={sliderValue}
                min={0}
                max={totalMonths}
                step={1}
                marks={marks}
                onChange={handleChange}
                disabled={disabled}
                valueLabelDisplay="auto"
                aria-label="Month slider"
        />
    )
}