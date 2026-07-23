import ky, {isHTTPError} from 'ky';
import {useAuthStore} from '#/stores/auth-store';
import {tryCatch} from '#/lib/catch';

// Ky Instance
export const api = ky.create({
	prefix: '/api',
	credentials: 'include',
	headers: {
		'Content-Type': 'application/json',
		Accept: 'application/json',
	},
	hooks: {
		afterResponse: [
			async ({request, options, response, retryCount}) => {
				// If request fails with 401, try refreshing the token
				if (response.status === 401 && retryCount === 0) {
					// Prevent refreshing on login, register, or refresh endpoints
					if (
						request.url.includes('auth/login') ||
						request.url.includes('auth/register') ||
						request.url.includes('auth/refresh')
					) {
						return response;
					}
					const {err} = await tryCatch(
						ky.post('/api/auth/refresh', {
							credentials: 'include',
						}),
					);
					if (err != null) {
						console.error('Refresh token failed:', err);
						useAuthStore.getState().logout();
						return response;
					}
					// Retry the original request with new cookies
					return ky(request, options);
				}
				return response;
			},
		],
	},
});

// Helper to parse backend error code and message
export function parseApiError(err: unknown): {
	code?: string;
	message: string;
} {
	if (isHTTPError(err) && err.data) {
		const data = err.data as {code?: string; message?: string};
		return {
			code: data.code,
			message:
				data.message || 'An unexpected error occurred. Please try again.',
		};
	}
	if (err instanceof Error) {
		return {message: err.message};
	}
	return {message: 'An unexpected error occurred. Please try again.'};
}
