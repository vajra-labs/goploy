import {Header, Wrapper} from '#auth/widgets';
import {SignInForm} from '#auth/forms/signin.form';
import {createFileRoute} from '@tanstack/react-router';

export const Route = createFileRoute('/_auth/singin')({
	component: RouteComponent,
});

function RouteComponent() {
	return (
		<Wrapper>
			<Header
				title="Welcome Back"
				subtitle="Sign in to your account to manage your deployments"
			/>
			<SignInForm />
		</Wrapper>
	);
}
