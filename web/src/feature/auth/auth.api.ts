import {api} from '#/lib/api';
import {tryCatch} from '#/lib/catch';
import type {User} from '#/types/user';
import type {SignInSchema} from './schema/signin.schema';
import type {SignUpSchema} from './schema/signup.schema';

export class AuthApi {
	login(data: SignInSchema) {
		return tryCatch(api.post('auth/login', {json: data}).json<any>());
	}
	register(data: SignUpSchema) {
		return tryCatch(api.post('auth/register', {json: data}).json<any>());
	}
	checkSetup() {
		return tryCatch(
			api.get('auth/setup').json<{isOwnerPresent: boolean}>(),
		);
	}
	whoami() {
		return tryCatch(api.get('user/me').json<User>());
	}
}

export default new AuthApi();
