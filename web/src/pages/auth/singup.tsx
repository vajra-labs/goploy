import {createFileRoute} from '@tanstack/react-router';
import {SignUpForm} from '#/feature/auth/signup-form';

export const Route = createFileRoute('/_auth/singup')({
	component: RouteComponent,
});

function RouteComponent() {
	return <SignUpForm />;
}
