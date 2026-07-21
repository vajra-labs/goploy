import * as v from 'valibot';
import {useForm} from 'react-hook-form';
import {valibotResolver} from '@hookform/resolvers/valibot';
import {createFileRoute} from '@tanstack/react-router';

export const Route = createFileRoute('/_auth/singin')({
	component: RouteComponent,
});

const signInSchema = v.object({
	email: v.pipe(v.string(), v.email('Invalid email address.')),
	password: v.pipe(
		v.string(),
		v.minLength(8, 'Password must be at least 8 characters.'),
	),
});

type SignInSchema = v.InferInput<typeof signInSchema>;

function RouteComponent() {
	const {
		register,
		handleSubmit,
		formState: {errors},
	} = useForm<SignInSchema>({
		resolver: valibotResolver(signInSchema),
	});

	const onSubmit = (data: SignInSchema) => {
		console.log('Validated Form Data:', data);
	};

	return (
		<div>
			<form onSubmit={handleSubmit(onSubmit)}>
				<div>
					<label>Email</label>
					<input type="email" {...register('email')} />
					{errors.email && <p>{errors.email.message}</p>}
				</div>

				<div>
					<label>Password</label>
					<input type="password" {...register('password')} />
					{errors.password && <p>{errors.password.message}</p>}
				</div>
				<button type="submit">Submit</button>
			</form>
		</div>
	);
}
