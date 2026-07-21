import {createFileRoute} from '@tanstack/react-router';

export const Route = createFileRoute('/_auth/singup')({
	component: RouteComponent,
});

function RouteComponent() {
	return <div>Hello "/singup"!</div>;
}
