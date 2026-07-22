import {useAuthStore} from '#/stores/auth-store';
import {createFileRoute, Navigate, Outlet} from '@tanstack/react-router';

export const Route = createFileRoute('/_app')({
	component: AppLayoutComponent,
});

function AppLayoutComponent() {
	const isAuth = useAuthStore(state => state.isAuth);

	if (!isAuth) {
		return <Navigate to="/singin" replace />;
	}

	return (
		<main className="flex-1 p-6">
			<Outlet />
		</main>
	);
}
