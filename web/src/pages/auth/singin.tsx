import {createFileRoute} from '@tanstack/react-router';
import {SignInForm} from '#/feature/auth/signin-form';

export const Route = createFileRoute('/_auth/singin')({
	component: RouteComponent,
});

function RouteComponent() {
	return <SignInForm />;
}
