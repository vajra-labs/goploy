import {api} from '#/lib/axios';
import {tryCatch} from '#/lib/catch';
import {createFileRoute} from '@tanstack/react-router';
import {Button} from '#/widgets/ui/button';

export const Route = createFileRoute('/_app/')({
	component: Home,
});

function Home() {
	const onHealth = async () => {
		const {res, err} = await tryCatch(api.get('/health'));
		if (err != null) {
			console.log(err.message);
			return;
		}
		console.log(res?.data);
	};

	return (
		<div className="flex min-h-svh p-6">
			<div className="flex max-w-md min-w-0 flex-col gap-4 text-sm leading-loose">
				<div>
					<h1 className="font-medium">Project ready!</h1>
					<p>You may now add components and start building.</p>
					<p>We&apos;ve already added the button component for you.</p>
					<Button className="mt-2" onClick={onHealth}>
						Button
					</Button>
				</div>
			</div>
		</div>
	);
}
