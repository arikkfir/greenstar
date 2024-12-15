import { useDomain } from "../hooks/domain.ts"

export const BaseAPIURL = `https://api.${useDomain()}`

export const QueryNilValue = "<nil>"
