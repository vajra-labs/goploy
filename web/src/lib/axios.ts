import axios from 'axios';

// Axios Instance
export const api = axios.create({
	baseURL: '/api',
	withCredentials: true,
	headers: {
		'Content-Type': 'application/json',
		Accept: 'application/json',
	},
});
