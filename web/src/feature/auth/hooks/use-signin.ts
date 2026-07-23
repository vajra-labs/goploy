import {toast} from 'sonner';
import authApi from '../auth.api';
import {parseApiError} from '#/lib/api';
import {useMutation} from '@tanstack/react-query';
import {useNavigate} from '@tanstack/react-router';
import {useAuthStore} from '#/stores/auth-store';

export function useSignin(onErrorCallback?: (msg: string) => void) {
	const navigate = useNavigate();
	const setAuth = useAuthStore(state => state.setAuth);

	return useMutation({
		mutationFn: authApi.login,
		onSuccess: ({res, err}) => {
			if (err != null) {
				const errorInfo = parseApiError(err);
				let message = errorInfo.message;
				if (errorInfo.code === 'INVALID_CREDENTIALS') {
					message = 'Invalid email or password. Please try again.';
				}
				onErrorCallback?.(message);
			} else if (res) {
				setAuth(res);
				toast.success('Successfully logged in!');
				navigate({to: '/'});
			}
		},
	});
}
