import {useAuthStore} from '#/stores/auth-store';
import {AuthLeftPanel} from '#/widgets/layout/auth-left-panel';
import {createFileRoute, Navigate, Outlet} from '@tanstack/react-router';

export const Route = createFileRoute('/_auth')({
	component: AuthLayoutComponent,
});

function AuthLayoutComponent() {
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
