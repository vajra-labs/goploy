// Define result types
type Success<T> = {res: T; err: null};
type Failure<E> = {res: null; err: E};
type Result<T, E = Error> = Success<T> | Failure<E>;

/**
 * Wraps an async operation and returns a typed result object.
 *
 * @example
 * const { data, error } = await tryCatch(fetchUser());
 */
export async function tryCatch<T, E = Error>(
	promise: Promise<T>,
): Promise<Result<T, E>> {
	try {
		const res = await promise;
		return {res, err: null};
	} catch (err) {
		return {res: null, err: err as E};
	}
}
