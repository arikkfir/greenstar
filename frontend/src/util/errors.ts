export const InternalError = new Error(`An internal error has occurred.`)
export const NotAuthenticatedError = new Error(`It seems you are not authenticated.`)

export class BadRequestError extends Error {
    constructor(
        public readonly code: number,
        message: string,
    ) {
        super(message)
    }
}
