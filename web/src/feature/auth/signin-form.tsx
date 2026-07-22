import * as React from 'react';
import * as v from 'valibot';
import {useForm} from 'react-hook-form';
import {valibotResolver} from '@hookform/resolvers/valibot';
import {useNavigate} from '@tanstack/react-router';
import {Mail, Lock, Eye, EyeOff, ArrowRight, Loader2} from 'lucide-react';
import {Button} from '#/widgets/ui/button';
import {Input} from '#/widgets/ui/input';
import {
	Field,
	FieldLabel,
	FieldError,
	FieldGroup,
} from '#/widgets/ui/field';
import {Logo} from '#/widgets/shared/logo';
import {api} from '#/lib/axios';
import {tryCatch} from '#/lib/catch';
import {useAuthStore} from '#/stores/auth-store';

export const signInSchema = v.object({
	email: v.pipe(v.string(), v.email('Invalid email address.')),
	password: v.pipe(
		v.string(),
		v.minLength(8, 'Password must be at least 8 characters.'),
	),
});

export type SignInSchema = v.InferInput<typeof signInSchema>;

export function SignInForm() {
	const [showPassword, setShowPassword] = React.useState(false);
	const [isLoading, setIsLoading] = React.useState(false);
	const [apiError, setApiError] = React.useState<string | null>(null);
	const navigate = useNavigate();
	const setAuth = useAuthStore(state => state.setAuth);

	const {
		register,
		handleSubmit,
		formState: {errors},
	} = useForm<SignInSchema>({
		resolver: valibotResolver(signInSchema),
	});

	const onSubmit = async (data: SignInSchema) => {
		setIsLoading(true);
		setApiError(null);
		const {res, err} = await tryCatch(api.post('/auth/login', data));
		if (err != null) {
			console.warn('Login with /auth/login failed...', err);
		} else if (res?.data) {
			setAuth(res.data);
			navigate({to: '/'});
		}
		setIsLoading(false);
	};

	return (
		<div className="relative z-10 w-full max-w-md rounded-2xl border border-border bg-card/70 p-8 shadow-2xl backdrop-blur-md transition-all hover:border-border">
			<div className="mb-8 flex flex-col items-center gap-2">
				<Logo className="size-16" />
				<h1 className="text-2xl font-bold tracking-tight text-foreground">
					Welcome back
				</h1>
				<p className="text-center text-sm text-muted-foreground">
					Sign in to your account to manage your deployments
				</p>
			</div>

			<form
				onSubmit={handleSubmit(onSubmit)}
				className="flex flex-col gap-6">
				{apiError && (
					<div className="animate-in rounded-lg border border-destructive/20 bg-destructive/10 p-3 text-sm text-destructive duration-200 fade-in slide-in-from-top-1">
						{apiError}
					</div>
				)}
				<FieldGroup>
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
						<div className="flex items-center justify-between">
							<FieldLabel
								htmlFor="password"
								className="text-xs font-semibold tracking-wider text-muted-foreground uppercase">
								Password
							</FieldLabel>
							<a
								href="#"
								onClick={e => e.preventDefault()}
								className="text-xs font-medium text-primary/80 transition-colors hover:text-primary hover:underline">
								Forgot password?
							</a>
						</div>
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
				</FieldGroup>
				<Button
					type="submit"
					className="flex h-10 w-full items-center justify-center rounded-lg bg-primary font-semibold text-primary-foreground transition-all hover:bg-primary/90 hover:shadow-lg hover:shadow-primary/10"
					disabled={isLoading}>
					{isLoading ? (
						<>
							<Loader2 className="animate-spin" data-icon="inline-start" />
							Signing in...
						</>
					) : (
						<>
							Sign In
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
