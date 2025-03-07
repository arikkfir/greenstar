import { useCallback } from "react"

export type CurrencyFormatter = (v: number, currency: string, minFracDigits?: number, maxFracDigits?: number) => string

export function useCurrencyFormatter(): CurrencyFormatter {
    const formatters: Map<string, Intl.NumberFormat> = new Map()
    return useCallback(
        (v: number, currency: string, minFractionDigits?: number, maxFractionDigits?: number): string => {
            const key = `${currency}/${minFractionDigits || 2}/${maxFractionDigits || 2}`
            if (!formatters.has(key)) {
                formatters.set(key, new Intl.NumberFormat(navigator.language, {
                    style: "currency",
                    currency,
                    minimumFractionDigits: minFractionDigits || 2,
                    maximumFractionDigits: maxFractionDigits || 2,
                }))
            }
            return formatters.get(key)!.format(v)
        },
        [formatters],
    )
}
