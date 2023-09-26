export const InternalError = new Error(`An internal error has occurred.`);

export class BadRequestError extends Error {
    constructor(public readonly code: number, message: string) {
        super(message);
    }
}
