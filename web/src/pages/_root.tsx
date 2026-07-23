import {
	Outlet,
	HeadContent,
	createRootRoute,
} from '@tanstack/react-router';
import authApi from '#auth/auth.api';
import {RootPending} from '#/widgets/pending';
import {useAuthStore} from '#/stores/auth-store';
import {TooltipProvider} from '#/widgets/ui/tooltip';
import {QueryClient, QueryClientProvider} from '@tanstack/react-query';
import {Toaster} from '#/widgets/ui/sonner';
import '#/styles/index.css';

export const Route = createRootRoute({
	beforeLoad: async () => {
		const {res} = await authApi.whoami();
		if (res) {
			useAuthStore.getState().setAuth(res);
		} else {
			useAuthStore.getState().logout();
		}
	},
	component: RootComponent,
	pendingComponent: RootPending,
});

// Create a client
const queryClient = new QueryClient();

function RootComponent() {
	return (
		<>
			<HeadContent />
			<QueryClientProvider client={queryClient}>
				<TooltipProvider>
					<Outlet />
					<Toaster />
				</TooltipProvider>
			</QueryClientProvider>
		</>
	);
}
