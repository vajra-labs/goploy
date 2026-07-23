import {toast} from 'sonner';
import authApi from '../auth.api';
import {useMutation} from '@tanstack/react-query';
import {useNavigate} from '@tanstack/react-router';
import {useAuthStore} from '#/stores/auth-store';
import {parseApiError} from '#/lib/api';

export function useSignup(onErrorCallback?: (msg: string) => void) {
	const navigate = useNavigate();
	const setAuth = useAuthStore(state => state.setAuth);

	return useMutation({
		mutationFn: authApi.register,
		onSuccess: ({res, err}) => {
			if (err != null) {
				const errorInfo = parseApiError(err);
				onErrorCallback?.(errorInfo.message);
			} else if (res) {
				setAuth(res);
				toast.success('Owner registered and logged in successfully!');
				navigate({to: '/'});
			}
		},
	});
}
