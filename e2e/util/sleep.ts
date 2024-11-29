/**
 * Sleep for a specified number of milliseconds.
 * @param seconds - The number of seconds to sleep.
 */
export function sleep(seconds: number) {
    return new Promise(resolve => setTimeout(resolve, seconds * 1000));
}
