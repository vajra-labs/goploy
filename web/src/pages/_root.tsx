import {
	Outlet,
	HeadContent,
	createRootRoute,
} from '@tanstack/react-router';
import {RootPending} from '#/widgets/pending';
import {useAuthStore} from '#/stores/auth-store';
import {TooltipProvider} from '#/widgets/ui/tooltip';
import {tryCatch} from '#/lib/catch';
import {api} from '#/lib/axios';
import '#/styles/index.css';

const fetchUser = async () => {
	const {res} = await tryCatch(api.get('/user/me'));
	if (res?.data) {
		useAuthStore.getState().setAuth(res.data);
	} else {
		useAuthStore.getState().logout();
	}
};

export const Route = createRootRoute({
	beforeLoad: fetchUser,
	component: RootComponent,
	pendingComponent: RootPending,
});

function RootComponent() {
	return (
		<>
			<HeadContent />
			<TooltipProvider>
				<Outlet />
			</TooltipProvider>
		</>
	);
}
