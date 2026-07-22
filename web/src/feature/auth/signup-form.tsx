import * as React from 'react';
import * as v from 'valibot';
import {useForm} from 'react-hook-form';
import {valibotResolver} from '@hookform/resolvers/valibot';
import {useNavigate} from '@tanstack/react-router';
import {
	Mail,
	Lock,
	User,
	Eye,
	EyeOff,
	ArrowRight,
	Loader2,
} from 'lucide-react';
import {Button} from '#/widgets/ui/button';
import {Input} from '#/widgets/ui/input';
import {
	Field,
	FieldLabel,
	FieldError,
	FieldGroup,
} from '#/widgets/ui/field';
import {api} from '#/lib/axios';
import {tryCatch} from '#/lib/catch';
import {Logo} from '#/widgets/shared/logo';
import {useAuthStore} from '#/stores/auth-store';

export const signUpSchema = v.pipe(
	v.object({
		name: v.pipe(v.string(), v.nonEmpty('First name is required')),
		lastName: v.pipe(v.string(), v.nonEmpty('Last name is required')),
		email: v.pipe(
			v.string(),
			v.nonEmpty('Email is required'),
			v.email('Email must be a valid email'),
		),
		password: v.pipe(
			v.string(),
			v.nonEmpty('Password is required'),
			v.minLength(8, 'Password must be at least 8 characters'),
		),
		confirmPassword: v.pipe(
			v.string(),
			v.nonEmpty('Confirm password is required'),
			v.minLength(8, 'Password must be at least 8 characters'),
		),
	}),
	v.forward(
		v.check(
			input => input.password === input.confirmPassword,
			'Passwords do not match',
		),
		['confirmPassword'],
	),
);

export type SignUpSchema = v.InferInput<typeof signUpSchema>;

export function SignUpForm() {
	const [showPassword, setShowPassword] = React.useState(false);
	const [showConfirmPassword, setShowConfirmPassword] =
		React.useState(false);
	const [isLoading, setIsLoading] = React.useState(false);
	const [apiError, setApiError] = React.useState<string | null>(null);
	const navigate = useNavigate();
	const setAuth = useAuthStore(state => state.setAuth);

	const {
		register,
		handleSubmit,
		formState: {errors},
	} = useForm<SignUpSchema>({
		resolver: valibotResolver(signUpSchema),
	});

	const onSubmit = async (data: SignUpSchema) => {
		setIsLoading(true);
		setApiError(null);
		const {res, err} = await tryCatch(api.post('/auth/register', data));
		if (err != null) {
			console.warn('Registration with /user/register failed...', err);
		} else if (res?.data) {
			setAuth(res.data);
			navigate({to: '/'});
		}
		setIsLoading(false);
	};

	return (
		<div className="relative z-10 w-full max-w-md rounded-2xl border border-border bg-card/70 p-8 shadow-2xl backdrop-blur-md transition-all hover:border-border">
			<div className="mb-6 flex flex-col items-center gap-2">
				<Logo className="size-16" />
				<h1 className="text-2xl font-bold tracking-tight text-foreground">
					Create an account
				</h1>
				<p className="text-center text-sm text-muted-foreground">
					Start managing your servers and deploying with GoPloy
				</p>
			</div>

			<form
				onSubmit={handleSubmit(onSubmit)}
				className="flex flex-col gap-5">
				{apiError && (
					<div className="animate-in rounded-lg border border-destructive/20 bg-destructive/10 p-3 text-sm text-destructive duration-200 fade-in slide-in-from-top-1">
						{apiError}
					</div>
				)}
				<FieldGroup>
					<div className="grid grid-cols-2 gap-4">
						<Field data-invalid={!!errors.name}>
							<FieldLabel
								htmlFor="name"
								className="text-xs font-semibold tracking-wider text-muted-foreground uppercase">
								First Name
							</FieldLabel>
							<div className="relative">
								<span className="absolute inset-y-0 left-0 flex items-center pl-3.5 text-muted-foreground">
									<User className="size-4" />
								</span>
								<Input
									id="name"
									type="text"
									placeholder="John"
									className="h-10 pl-10"
									disabled={isLoading}
									aria-invalid={!!errors.name}
									{...register('name')}
								/>
							</div>
							<FieldError errors={[errors.name]} />
						</Field>
						<Field data-invalid={!!errors.lastName}>
							<FieldLabel
								htmlFor="lastName"
								className="text-xs font-semibold tracking-wider text-muted-foreground uppercase">
								Last Name
							</FieldLabel>
							<div className="relative">
								<span className="absolute inset-y-0 left-0 flex items-center pl-3.5 text-muted-foreground">
									<User className="size-4" />
								</span>
								<Input
									id="lastName"
									type="text"
									placeholder="Doe"
									className="h-10 pl-10"
									disabled={isLoading}
									aria-invalid={!!errors.lastName}
									{...register('lastName')}
								/>
							</div>
							<FieldError errors={[errors.lastName]} />
						</Field>
					</div>
					<Field data-invalid={!!errors.email}>
						<FieldLabel
							htmlFor="email"
							className="text-xs font-semibold tracking-wider text-muted-foreground uppercase">
							Email Address
						</FieldLabel>
						<div className="relative">
							<span className="absolute inset-y-0 left-0 flex items-center pl-3.5 text-muted-foreground">
								<Mail className="size-4" />
							</span>
							<Input
								id="email"
								type="email"
								placeholder="name@example.com"
								className="h-10 pl-10"
								disabled={isLoading}
								aria-invalid={!!errors.email}
								{...register('email')}
							/>
						</div>
						<FieldError errors={[errors.email]} />
					</Field>
					<Field data-invalid={!!errors.password}>
						<FieldLabel
							htmlFor="password"
							className="text-xs font-semibold tracking-wider text-muted-foreground uppercase">
							Password
						</FieldLabel>
						<div className="relative">
							<span className="absolute inset-y-0 left-0 flex items-center pl-3.5 text-muted-foreground">
								<Lock className="size-4" />
							</span>
							<Input
								id="password"
								type={showPassword ? 'text' : 'password'}
								placeholder="••••••••"
								className="h-10 pr-10 pl-10"
								disabled={isLoading}
								aria-invalid={!!errors.password}
								{...register('password')}
							/>
							<button
								type="button"
								onClick={() => setShowPassword(!showPassword)}
								className="absolute inset-y-0 right-0 flex items-center pr-3 text-muted-foreground transition-colors hover:text-foreground"
								tabIndex={-1}>
								{showPassword ? (
									<EyeOff className="size-4" />
								) : (
									<Eye className="size-4" />
								)}
							</button>
						</div>
						<FieldError errors={[errors.password]} />
					</Field>
					<Field data-invalid={!!errors.confirmPassword}>
						<FieldLabel
							htmlFor="confirmPassword"
							className="text-xs font-semibold tracking-wider text-muted-foreground uppercase">
							Confirm Password
						</FieldLabel>
						<div className="relative">
							<span className="absolute inset-y-0 left-0 flex items-center pl-3.5 text-muted-foreground">
								<Lock className="size-4" />
							</span>
							<Input
								id="confirmPassword"
								type={showConfirmPassword ? 'text' : 'password'}
								placeholder="••••••••"
								className="h-10 pr-10 pl-10"
								disabled={isLoading}
								aria-invalid={!!errors.confirmPassword}
								{...register('confirmPassword')}
							/>
							<button
								type="button"
								onClick={() =>
									setShowConfirmPassword(!showConfirmPassword)
								}
								className="absolute inset-y-0 right-0 flex items-center pr-3 text-muted-foreground transition-colors hover:text-foreground"
								tabIndex={-1}>
								{showConfirmPassword ? (
									<EyeOff className="size-4" />
								) : (
									<Eye className="size-4" />
								)}
							</button>
						</div>
						<FieldError errors={[errors.confirmPassword]} />
					</Field>
				</FieldGroup>
				<Button
					type="submit"
					className="flex h-10 w-full items-center justify-center rounded-lg bg-primary font-semibold text-primary-foreground transition-all hover:bg-primary/90 hover:shadow-lg hover:shadow-primary/10"
					disabled={isLoading}>
					{isLoading ? (
						<>
							<Loader2 className="animate-spin" data-icon="inline-start" />
							Creating account...
						</>
					) : (
						<>
							Create Account
							<ArrowRight
								className="transition-transform group-hover/button:translate-x-0.5"
								data-icon="inline-end"
							/>
						</>
					)}
				</Button>
			</form>
		</div>
	);
}
