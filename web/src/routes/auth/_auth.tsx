import {useAuthStore} from '#/stores/auth-store';
import {createFileRoute, Navigate, Outlet} from '@tanstack/react-router';

export const Route = createFileRoute('/_auth')({
	component: AuthLayoutComponent,
});

function AuthLayoutComponent() {
	const isAuth = useAuthStore(state => state.isAuth);
	if (isAuth) return <Navigate to="/" replace />;
	return <Outlet />;
}
