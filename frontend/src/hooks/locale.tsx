import {useCallback, useContext, useMemo} from "react";
import {LocaleContext} from "../providers/LocaleProvider.tsx";

export function useCurrencyFormatter() {
    const locale = useContext(LocaleContext)
    const currencyFormat = useMemo(() => {
        const language = locale.language || navigator.language
        const currency = locale.currency || 'USD'
        return new Intl.NumberFormat(language, {style: 'currency', currency: currency})
    }, [locale, navigator.language])
    return useCallback((v: number) => (currencyFormat.format(v)), [currencyFormat])
}

export function useDateFormatter(style?: "full" | "long" | "medium" | "short") {
    const locale = useContext(LocaleContext)
    const dateTimeFormat = useMemo(() => {
        const language = locale.language || navigator.language
        const dateStyle = style || 'short'
        return new Intl.DateTimeFormat(language, {dateStyle})
    }, [locale, style, navigator.language])
    return useCallback((v: Date | number) => (dateTimeFormat.format(v)), [dateTimeFormat])
}
