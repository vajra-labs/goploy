import {
	Outlet,
	HeadContent,
	createRootRoute,
} from '@tanstack/react-router';
import {api} from '#/lib/axios';
import {tryCatch} from '#/lib/catch';
import {RootPending} from '#/widgets/pending';
import {useAuthStore} from '#/stores/auth-store';
import '../index.css';

export const Route = createRootRoute({
	beforeLoad: async () => {
		const {res} = await tryCatch(api.get('/user/me'));
		if (res?.data) {
			useAuthStore.getState().setAuth(res.data);
		} else {
			useAuthStore.getState().logout();
		}
	},
	component: RootComponent,
	pendingComponent: RootPending,
});

function RootComponent() {
	return (
		<>
			<HeadContent />
			<Outlet />
		</>
	);
}
