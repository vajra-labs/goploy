import {
	Outlet,
	redirect,
	Navigate,
	createFileRoute,
} from '@tanstack/react-router';
import authApi from '#auth/auth.api';
import {useAuthStore} from '#/stores/auth-store';
import {AuthLeftPanel} from '#/widgets/layout/auth-left-panel';

export const Route = createFileRoute('/_auth')({
	beforeLoad: async ({location}) => {
		// Verify setup status
		const {res} = await authApi.checkSetup();
		const isOwnerPresent = res?.isOwnerPresent ?? true;
		// Owner not present & user is not on signup
		if (!isOwnerPresent && location.pathname !== '/singup') {
			throw redirect({
				to: '/singup',
				replace: true,
			});
		}
		// Owner is present & user attempts to signup
		if (isOwnerPresent && location.pathname === '/singup') {
			throw redirect({
				to: '/singin',
				replace: true,
			});
		}
	},
	component: AuthLayout,
});

function AuthLayout() {
	const isAuth = useAuthStore(state => state.isAuth);
	if (isAuth) return <Navigate to="/" replace />;
	return (
		<div className="flex min-h-svh w-screen bg-background">
			{/* Left Side */}
			<AuthLeftPanel />
			{/* Right Side */}
			<div className="flex w-full items-center justify-center bg-background p-6 md:p-10 lg:w-1/2">
				<Outlet />
			</div>
		</div>
	);
}
